package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config armazena as configurações da aplicação
type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	LogLevel   string
}

// Load carrega as configurações a partir de variáveis de ambiente
func Load() *Config {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "3000"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "produtos_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		LogLevel:   getEnv("LOG_LEVEL", "debug"),
	}

	return cfg
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
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
