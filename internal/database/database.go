package database

import (
	"fmt"

	"api_fibergorm/internal/config"
	"api_fibergorm/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect estabelece conexão com o banco de dados PostgreSQL
func Connect(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	log.WithFields(logrus.Fields{
		"host": cfg.DBHost,
		"port": cfg.DBPort,
		"db":   cfg.DBName,
	}).Info("Conectando ao banco de dados PostgreSQL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.WithError(err).Error("Falha ao conectar ao banco de dados")
		return nil, err
	}

	log.Info("Conexão com o banco de dados estabelecida com sucesso")
	return db, nil
}

// Migrate executa as migrações automáticas do GORM
func Migrate(db *gorm.DB, log *logrus.Logger) error {
	log.Info("Executando migrações do banco de dados")

	err := db.AutoMigrate(&models.Produto{})
	if err != nil {
		log.WithError(err).Error("Falha ao executar migrações")
		return err
	}

	log.Info("Migrações executadas com sucesso")
	return nil
}
