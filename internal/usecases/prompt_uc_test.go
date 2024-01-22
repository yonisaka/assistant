package usecases_test

import (
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"github.com/yonisaka/assistant/internal/usecases"
)

type promptFields struct {
	connector connector.Connector
}

func sut(f promptFields) usecases.PromptUsecase {
	return usecases.NewPromptUsecase(
		f.connector,
	)
}
