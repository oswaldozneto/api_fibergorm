package routes

import (
	"api_fibergorm/internal/handler"
	"api_fibergorm/internal/repository"
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

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "api-produtos",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// Setup produto routes
	setupProdutoRoutes(api, db, log)
}

// setupProdutoRoutes configura as rotas de produtos
func setupProdutoRoutes(router fiber.Router, db *gorm.DB, log *logrus.Logger) {
	// Inicializa as camadas
	produtoRepo := repository.NewProdutoRepository(db)
	produtoService := service.NewProdutoService(produtoRepo, log)
	produtoHandler := handler.NewProdutoHandler(produtoService, log)

	// Grupo de rotas de produtos
	produtos := router.Group("/produtos")
	{
		produtos.Post("/", produtoHandler.Create)
		produtos.Get("/", produtoHandler.GetAll)
		produtos.Get("/:id", produtoHandler.GetByID)
		produtos.Put("/:id", produtoHandler.Update)
		produtos.Delete("/:id", produtoHandler.Delete)
	}
}
