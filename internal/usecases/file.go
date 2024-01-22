package usecases

import (
	"context"
	"net/http"

	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
)

// GetListFile is a function to get list file from OpenAI API
func (u *fileUsecase) GetListFile(ctx context.Context) (*[]repository.File, error) {
	// Set HTTP Request Parameter
	httpRequestOption := &connector.RequestOption{
		Method: http.MethodGet,
		URL:    "/files",
	}

	var response *connector.OpenAIFile
	// Do HTTP Request
	err := u.connector.Send(ctx, httpRequestOption, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
