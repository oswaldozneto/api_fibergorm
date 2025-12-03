package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

// SetupMiddlewares configura os middlewares globais da aplicação
func SetupMiddlewares(app *fiber.App, log *logrus.Logger) {
	// Recover middleware para capturar panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Request ID para rastreamento
	app.Use(requestid.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Logger middleware customizado com Logrus
	app.Use(LoggerMiddleware(log))
}

// LoggerMiddleware middleware para logging das requisições
func LoggerMiddleware(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Processa a requisição
		err := c.Next()

		// Calcula o tempo de resposta
		latency := time.Since(start)

		// Log da requisição
		log.WithFields(logrus.Fields{
			"request_id": c.Locals("requestid"),
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"latency":    latency.String(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		}).Info("Requisição HTTP")

		return err
	}
}
