package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Router struct {
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) *Router {
	return &Router{logger: logger}
}

func (r *Router) MustInit() *fiber.App {
	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, world!")
	})

	return app
}
