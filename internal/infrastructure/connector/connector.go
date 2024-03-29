package connector

import (
	"context"
	"io"
)

//go:generate rm -f ./connector_mock.go
//go:generate mockgen -destination connector_mock.go -package connector -mock_names Connector=GoMockConnector -source connector.go

type connector struct {
	openai *OpenAI
}

// RequestOption is a struct to set HTTP Request Parameter
type RequestOption struct {
	Method       string
	URL          string
	Body         io.Reader
	CustomHeader map[string]string
}

// NewConnector is a function to create new HTTP connector
func NewConnector(openai *OpenAI) Connector {
	return &connector{
		openai: openai,
	}
}

// Connector is an interface to send HTTP Request
type Connector interface {
	Send(ctx context.Context, requestOption *RequestOption, response any) error
}

const (
	HeaderAuthorization = "Authorization"
	BearerAuthType      = "Bearer"
	HeaderKeyOpenAIBeta = "OpenAI-Beta"
	AssistantV1         = "assistants=v1"
)
