package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// LokiConfig configurações do Loki
type LokiConfig struct {
	URL         string            // URL do Loki push endpoint
	BatchSize   int               // Número de logs para enviar em batch (padrão: 10)
	BatchWait   time.Duration     // Tempo máximo para aguardar antes de enviar batch (padrão: 5s)
	Labels      map[string]string // Labels estáticos para todos os logs
	ServiceName string            // Nome do serviço/job
	Enabled     bool              // Se a integração está habilitada
	Timeout     time.Duration     // Timeout para requisições HTTP (padrão: 10s)
}

// LokiHook hook do Logrus para enviar logs ao Loki
type LokiHook struct {
	config   LokiConfig
	client   *http.Client
	entries  []lokiEntry
	mutex    sync.Mutex
	quit     chan struct{}
	hostname string
}

// lokiEntry representa uma entrada de log
type lokiEntry struct {
	timestamp time.Time
	line      string
	level     string
}

// lokiPushRequest estrutura de requisição do Loki
type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

// lokiStream representa um stream de logs
type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// NewLokiHook cria um novo hook do Loki
func NewLokiHook(config LokiConfig) (*LokiHook, error) {
	if !config.Enabled {
		return nil, nil
	}

	if config.URL == "" {
		return nil, fmt.Errorf("loki URL é obrigatória")
	}

	// Valores padrão
	if config.BatchSize <= 0 {
		config.BatchSize = 10
	}
	if config.BatchWait <= 0 {
		config.BatchWait = 5 * time.Second
	}
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Second
	}
	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	hostname, _ := os.Hostname()

	hook := &LokiHook{
		config:   config,
		client:   &http.Client{Timeout: config.Timeout},
		entries:  make([]lokiEntry, 0, config.BatchSize),
		quit:     make(chan struct{}),
		hostname: hostname,
	}

	// Inicia goroutine para flush periódico
	go hook.runFlusher()

	return hook, nil
}

// Levels retorna os níveis de log que o hook suporta
func (h *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire é chamado quando um log é emitido
func (h *LokiHook) Fire(entry *logrus.Entry) error {
	// Formata a linha de log como JSON
	line, err := h.formatEntry(entry)
	if err != nil {
		return err
	}

	h.mutex.Lock()
	h.entries = append(h.entries, lokiEntry{
		timestamp: entry.Time,
		line:      line,
		level:     entry.Level.String(),
	})

	// Envia se atingiu o batch size
	shouldFlush := len(h.entries) >= h.config.BatchSize
	h.mutex.Unlock()

	if shouldFlush {
		go h.flush()
	}

	return nil
}

// formatEntry formata uma entrada de log como JSON
func (h *LokiHook) formatEntry(entry *logrus.Entry) (string, error) {
	data := make(map[string]interface{})

	// Adiciona campos do log
	for k, v := range entry.Data {
		data[k] = v
	}

	// Adiciona campos padrão
	data["level"] = entry.Level.String()
	data["msg"] = entry.Message
	data["time"] = entry.Time.Format(time.RFC3339Nano)
	data["hostname"] = h.hostname
	data["service"] = h.config.ServiceName

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// runFlusher executa flush periódico
func (h *LokiHook) runFlusher() {
	ticker := time.NewTicker(h.config.BatchWait)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.flush()
		case <-h.quit:
			h.flush() // Flush final antes de sair
			return
		}
	}
}

// flush envia os logs acumulados para o Loki
func (h *LokiHook) flush() {
	h.mutex.Lock()
	if len(h.entries) == 0 {
		h.mutex.Unlock()
		return
	}

	entries := h.entries
	h.entries = make([]lokiEntry, 0, h.config.BatchSize)
	h.mutex.Unlock()

	// Agrupa por nível de log
	streamsByLevel := make(map[string][][]string)
	for _, entry := range entries {
		ts := strconv.FormatInt(entry.timestamp.UnixNano(), 10)
		streamsByLevel[entry.level] = append(streamsByLevel[entry.level], []string{ts, entry.line})
	}

	// Cria streams para cada nível
	var streams []lokiStream
	for level, values := range streamsByLevel {
		labels := make(map[string]string)
		for k, v := range h.config.Labels {
			labels[k] = v
		}
		labels["job"] = h.config.ServiceName
		labels["level"] = level
		labels["hostname"] = h.hostname

		streams = append(streams, lokiStream{
			Stream: labels,
			Values: values,
		})
	}

	// Envia para o Loki
	h.send(lokiPushRequest{Streams: streams})
}

// send envia a requisição para o Loki
func (h *LokiHook) send(req lokiPushRequest) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao serializar logs para Loki: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", h.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao criar requisição para Loki: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao enviar logs para Loki: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "Loki retornou status %d\n", resp.StatusCode)
	}
}

// Close fecha o hook e faz flush final
func (h *LokiHook) Close() {
	close(h.quit)
}

// DefaultLokiConfig retorna configuração padrão do Loki
func DefaultLokiConfig() LokiConfig {
	return LokiConfig{
		URL:         getEnv("LOKI_URL", "http://10.110.0.239:3100/loki/api/v1/push"),
		BatchSize:   getEnvAsInt("LOKI_BATCH_SIZE", 10),
		BatchWait:   time.Duration(getEnvAsInt("LOKI_BATCH_WAIT_SECONDS", 5)) * time.Second,
		ServiceName: getEnv("LOKI_SERVICE_NAME", "ARQUITETURA_FIBER_GORM"),
		Enabled:     getEnvAsBool("LOKI_ENABLED", true),
		Timeout:     time.Duration(getEnvAsInt("LOKI_TIMEOUT_SECONDS", 10)) * time.Second,
		Labels: map[string]string{
			"app":         "api_fibergorm",
			"environment": getEnv("ENVIRONMENT", "development"),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
