package usecases

import (
	"context"

	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
)

type PromptUsecase interface {
	SendPrompt(ctx context.Context, message string) (*repository.Message, error)
}

func NewPromptUsecase(connector connector.Connector) PromptUsecase {
	return &promptUsecase{
		connector: connector,
	}
}

type promptUsecase struct {
	connector connector.Connector
}
