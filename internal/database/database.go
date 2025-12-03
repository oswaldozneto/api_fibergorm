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
	// Primeiro, tenta criar o banco de dados se não existir
	if err := createDatabaseIfNotExists(cfg, log); err != nil {
		log.WithError(err).Warn("Não foi possível verificar/criar o banco de dados")
	}

	// Conecta ao banco de dados da aplicação
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

// createDatabaseIfNotExists conecta ao postgres e cria o banco se não existir
func createDatabaseIfNotExists(cfg *config.Config, log *logrus.Logger) error {
	// Conecta ao banco postgres padrão para criar o banco da aplicação
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco postgres: %w", err)
	}

	// Fecha a conexão ao final
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("falha ao obter conexão SQL: %w", err)
	}
	defer sqlDB.Close()

	// Verifica se o banco de dados já existe
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	if err := sqlDB.QueryRow(query, cfg.DBName).Scan(&exists); err != nil {
		return fmt.Errorf("falha ao verificar existência do banco: %w", err)
	}

	if !exists {
		log.WithField("database", cfg.DBName).Info("Criando banco de dados")

		// Cria o banco de dados
		createQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if _, err := sqlDB.Exec(createQuery); err != nil {
			return fmt.Errorf("falha ao criar banco de dados: %w", err)
		}

		log.WithField("database", cfg.DBName).Info("Banco de dados criado com sucesso")
	} else {
		log.WithField("database", cfg.DBName).Debug("Banco de dados já existe")
	}

	return nil
}

// Migrate executa as migrações automáticas do GORM
func Migrate(db *gorm.DB, log *logrus.Logger) error {
	log.Info("Executando migrações do banco de dados")

	// Passo 1: Migra a tabela de categorias primeiro
	if err := db.AutoMigrate(&models.Categoria{}); err != nil {
		log.WithError(err).Error("Falha ao migrar tabela de categorias")
		return err
	}

	// Passo 2: Verifica se a coluna categoria_id já existe na tabela produtos
	var hasColumn bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'produtos' AND column_name = 'categoria_id'
		)
	`).Scan(&hasColumn).Error

	if err != nil {
		// Tabela pode não existir ainda, continua normalmente
		hasColumn = false
	}

	// Passo 3: Se a coluna não existe, precisamos fazer uma migração especial
	if !hasColumn {
		log.Info("Coluna categoria_id não existe, executando migração especial")

		// Verifica se a tabela produtos existe e tem registros
		var produtosCount int64
		db.Table("produtos").Count(&produtosCount)

		if produtosCount > 0 {
			log.Info("Encontrados produtos existentes, aplicando migração com dados")

			// Adiciona a coluna como NULL primeiro
			if err := db.Exec("ALTER TABLE produtos ADD COLUMN IF NOT EXISTS categoria_id BIGINT").Error; err != nil {
				log.WithError(err).Error("Falha ao adicionar coluna categoria_id")
				return err
			}

			// Executa o seed para criar categoria padrão e atualizar produtos
			if err := Seed(db, log); err != nil {
				return err
			}

			// Agora aplica NOT NULL
			if err := db.Exec("ALTER TABLE produtos ALTER COLUMN categoria_id SET NOT NULL").Error; err != nil {
				log.WithError(err).Error("Falha ao definir NOT NULL em categoria_id")
				return err
			}

			// Adiciona a FK se não existir
			if err := db.Exec(`
				DO $$ 
				BEGIN
					IF NOT EXISTS (
						SELECT 1 FROM information_schema.table_constraints 
						WHERE constraint_name = 'fk_produtos_categoria'
					) THEN
						ALTER TABLE produtos 
						ADD CONSTRAINT fk_produtos_categoria 
						FOREIGN KEY (categoria_id) REFERENCES categorias(id);
					END IF;
				END $$;
			`).Error; err != nil {
				log.WithError(err).Warn("Aviso ao criar FK (pode já existir)")
			}

			log.Info("Migração especial concluída")
			return nil
		}
	}

	// Passo 4: Migração normal do GORM
	if err := db.AutoMigrate(&models.Produto{}); err != nil {
		log.WithError(err).Error("Falha ao migrar tabela de produtos")
		return err
	}

	log.Info("Migrações executadas com sucesso")
	return nil
}
