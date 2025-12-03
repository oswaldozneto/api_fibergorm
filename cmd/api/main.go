package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"api_fibergorm/internal/config"
	"api_fibergorm/internal/database"
	"api_fibergorm/internal/middleware"
	"api_fibergorm/internal/routes"

	"github.com/gofiber/fiber/v2"
)

// @title API Produtos
// @version 1.0
// @description API de Produtos - POC com Fiber e GORM
// @termsOfService http://swagger.io/terms/

// @contact.name Suporte API
// @contact.email suporte@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	// Carrega as configurações
	cfg := config.Load()

	// Configura o logger
	log := config.SetupLogger(cfg.LogLevel)

	log.Info("Iniciando API de Produtos - POC Fiber + GORM")

	// Conecta ao banco de dados
	db, err := database.Connect(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Falha ao conectar ao banco de dados")
	}

	// Executa as migrações
	if err := database.Migrate(db, log); err != nil {
		log.WithError(err).Fatal("Falha ao executar migrações")
	}

	// Executa o seed de dados iniciais (categoria padrão, etc.)
	if err := database.Seed(db, log); err != nil {
		log.WithError(err).Fatal("Falha ao executar seed de dados")
	}

	// Cria a aplicação Fiber
	app := fiber.New(fiber.Config{
		AppName:      "API Produtos v1.0",
		ErrorHandler: customErrorHandler,
	})

	// Configura os middlewares
	middleware.SetupMiddlewares(app, log)

	// Configura as rotas
	routes.SetupRoutes(app, db, log)

	// Canal para capturar sinais de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.ServerPort)
		log.WithField("port", cfg.ServerPort).Info("Servidor iniciado")
		if err := app.Listen(addr); err != nil {
			log.WithError(err).Fatal("Erro ao iniciar o servidor")
		}
	}()

	// Aguarda sinal de shutdown
	<-quit
	log.Info("Encerrando servidor...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		log.WithError(err).Error("Erro ao encerrar servidor")
	}

	log.Info("Servidor encerrado com sucesso")
}

// customErrorHandler trata erros globais da aplicação
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
