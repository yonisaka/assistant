package di

import (
	"fmt"
	"os"

	"github.com/yonisaka/assistant/internal/infrastructure/connector"
)

// GetConnector is a function to get connector
func GetConnector() connector.Connector {
	return connector.NewConnector(
		&connector.OpenAI{
			BaseURL: os.Getenv("OPENAI_V1_BASE_URL"),
			Header: map[string]string{
				connector.HeaderAuthorization: fmt.Sprintf("%s %s", connector.BearerAuthType, os.Getenv("OPENAI_API_KEY")),
			},
		},
	)
}
