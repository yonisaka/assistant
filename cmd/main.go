package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/yonisaka/assistant/internal/di"
)

func main() {
	// Create new Fiber instance
	app := fiber.New()

	// Logging Request ID
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	// Provide healthcheck
	app.Use(healthcheck.New())

	// Initialize Router
	di.GetRouter(app)
	// Listen on port 3000
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
