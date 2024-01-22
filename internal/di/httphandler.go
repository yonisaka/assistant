package di

import "github.com/yonisaka/assistant/internal/adapters/httphandler"

// GetFileHandler is a function to get http openAI handler
func GetFileHandler() httphandler.FileHandler {
	return httphandler.NewFileHandler(
		GetFileUsecase(),
	)
}

// GetPromptHandler is a function to get http openAI handler
func GetPromptHandler() httphandler.PromptHandler {
	return httphandler.NewPromptHandler(
		GetPromptUsecase(),
	)
}
