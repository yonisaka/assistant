package usecases

import (
	"context"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
)

type FileUsecase interface {
	GetListFile(ctx context.Context) (*[]repository.File, error)
}

func NewFileUsecase(connector connector.Connector) FileUsecase {
	return &fileUsecase{
		connector: connector,
	}
}

type fileUsecase struct {
	connector connector.Connector
}
