package di

import "github.com/yonisaka/assistant/internal/usecases"

// GetFileUsecase is a function to get usecase
func GetFileUsecase() usecases.FileUsecase {
	return usecases.NewFileUsecase(
		GetConnector(),
	)
}

// GetPromptUsecase is a function to get usecase
func GetPromptUsecase() usecases.PromptUsecase {
	return usecases.NewPromptUsecase(
		GetConnector(),
	)
}
