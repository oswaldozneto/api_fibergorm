package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Config armazena as configurações da aplicação
// Todas as variáveis de ambiente são opcionais e possuem valores padrão
type Config struct {
	// Servidor
	ServerPort         string // SERVER_PORT (padrão: 3000)
	ServerReadTimeout  int    // SERVER_READ_TIMEOUT em segundos (padrão: 10)
	ServerWriteTimeout int    // SERVER_WRITE_TIMEOUT em segundos (padrão: 10)

	// Banco de Dados PostgreSQL
	DBHost            string // DB_HOST (padrão: localhost)
	DBPort            string // DB_PORT (padrão: 5432)
	DBUser            string // DB_USER (padrão: postgres)
	DBPassword        string // DB_PASSWORD (padrão: postgres)
	DBName            string // DB_NAME (padrão: produtos_db)
	DBSSLMode         string // DB_SSLMODE (padrão: disable) - valores: disable, require, verify-ca, verify-full
	DBMaxOpenConns    int    // DB_MAX_OPEN_CONNS (padrão: 10)
	DBMaxIdleConns    int    // DB_MAX_IDLE_CONNS (padrão: 5)
	DBConnMaxLifetime int    // DB_CONN_MAX_LIFETIME em minutos (padrão: 30)

	// Logging
	LogLevel  string // LOG_LEVEL (padrão: debug) - valores: debug, info, warn, error
	LogFormat string // LOG_FORMAT (padrão: json) - valores: json, text
}

// Load carrega as configurações a partir de variáveis de ambiente
// Todas as variáveis são opcionais e possuem valores padrão sensatos
func Load() *Config {
	cfg := &Config{
		// Servidor
		ServerPort:         getEnv("SERVER_PORT", "3000"),
		ServerReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
		ServerWriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),

		// Banco de Dados
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "admin"),
		DBName:            getEnv("DB_NAME", "produtos_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 30),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}

	return cfg
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da variável de ambiente como int ou o valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool retorna o valor da variável de ambiente como bool ou o valor padrão
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// SetupLogger configura o logger da aplicação
func SetupLogger(level string) *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	log.SetLevel(logLevel)

	return log
}

// PrintConfig exibe as configurações carregadas (útil para debug)
func (c *Config) PrintConfig(log *logrus.Logger) {
	log.WithFields(logrus.Fields{
		"server_port": c.ServerPort,
		"db_host":     c.DBHost,
		"db_port":     c.DBPort,
		"db_name":     c.DBName,
		"db_user":     c.DBUser,
		"db_sslmode":  c.DBSSLMode,
		"log_level":   c.LogLevel,
	}).Info("Configurações carregadas")
}
