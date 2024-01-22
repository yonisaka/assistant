package di

import (
	"github.com/gofiber/fiber/v2"
)

func GetRouter(app *fiber.App) {
	// API Group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	fileHandler := GetFileHandler()
	v1.Get("/files", fileHandler.GetListFile)

	promptHandler := GetPromptHandler()
	v1.Post("/prompt", promptHandler.SendPrompt)
}
