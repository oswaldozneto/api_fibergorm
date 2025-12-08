package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestsTotal contador de requisições HTTP
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP recebidas",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration histograma de duração das requisições
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestsInFlight gauge de requisições em andamento
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Número de requisições HTTP em processamento",
		},
	)

	// HTTPResponseSize histograma do tamanho das respostas
	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Tamanho das respostas HTTP em bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"method", "path", "status"},
	)

	// DatabaseQueriesTotal contador de queries no banco de dados
	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total de queries executadas no banco de dados",
		},
		[]string{"operation", "table"},
	)

	// DatabaseQueryDuration histograma de duração das queries
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duração das queries no banco de dados em segundos",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"operation", "table"},
	)
)

// PrometheusMiddleware middleware para coletar métricas das requisições HTTP
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ignora o endpoint de métricas para não criar recursão
		if c.Path() == "/metrics" {
			return c.Next()
		}

		// Incrementa requisições em andamento
		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		// Marca o início da requisição
		start := time.Now()

		// Processa a requisição
		err := c.Next()

		// Calcula a duração
		duration := time.Since(start).Seconds()

		// Obtém informações da resposta
		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()
		path := c.Route().Path // Usa o padrão da rota para evitar cardinalidade alta

		// Se não encontrou a rota, usa o path original (para 404s)
		if path == "" {
			path = c.Path()
		}

		// Registra as métricas
		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
		HTTPResponseSize.WithLabelValues(method, path, status).Observe(float64(len(c.Response().Body())))

		return err
	}
}

// MetricsHandler retorna o handler do Prometheus para expor as métricas
func MetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}

// RecordDatabaseQuery registra métricas de uma query no banco de dados
func RecordDatabaseQuery(operation, table string, duration time.Duration) {
	DatabaseQueriesTotal.WithLabelValues(operation, table).Inc()
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

