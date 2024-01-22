package usecases_test

import (
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"github.com/yonisaka/assistant/internal/usecases"
)

type fileFields struct {
	connector connector.Connector
}

func fileSut(f fileFields) usecases.FileUsecase {
	return usecases.NewFileUsecase(
		f.connector,
	)
}
