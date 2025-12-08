package routes

import (
	"api_fibergorm/internal/handler"
	"api_fibergorm/internal/metrics"
	"api_fibergorm/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	_ "api_fibergorm/docs" // Importa a documentação gerada pelo swag
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(app *fiber.App, db *gorm.DB, log *logrus.Logger) {
	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Prometheus metrics endpoint
	app.Get("/metrics", metrics.MetricsHandler())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "api-produtos",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// Setup das rotas usando a nova arquitetura
	setupCategoriaRoutes(api, db, log)
	setupProdutoRoutes(api, db, log)
}

// setupCategoriaRoutes configura as rotas de categorias
func setupCategoriaRoutes(router fiber.Router, db *gorm.DB, log *logrus.Logger) {
	// Cria o serviço (que já configura repositório, mapper e validator internamente)
	categoriaService := service.NewCategoriaService(db, log)

	// Cria o handler
	categoriaHandler := handler.NewCategoriaHandler(categoriaService, log)

	// Registra as rotas
	categorias := router.Group("/categorias")
	categoriaHandler.RegisterRoutes(categorias)
}

// setupProdutoRoutes configura as rotas de produtos
func setupProdutoRoutes(router fiber.Router, db *gorm.DB, log *logrus.Logger) {
	// Cria o serviço (que já configura repositório, mapper e validator internamente)
	produtoService := service.NewProdutoService(db, log)

	// Cria o handler
	produtoHandler := handler.NewProdutoHandler(produtoService, log)

	// Registra as rotas
	produtos := router.Group("/produtos")
	produtoHandler.RegisterRoutes(produtos)
}
