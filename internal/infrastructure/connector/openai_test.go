package connector_test

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"net/http"
	"os"
	"testing"
)

type fields struct {
	openai *connector.OpenAI
}

func sut(f fields) connector.Connector {
	return connector.NewConnector(
		f.openai,
	)
}

func TestMain(m *testing.M) {
	var code int

	defer func() {
		os.Exit(code)
	}()

	_ = os.Setenv("OPENAI_API_URL", "test")
	_ = os.Setenv("OPENAI_API_KEY", "test")

	code = m.Run()
}

func TestConnector_Send(t *testing.T) {
	type args struct {
		ctx           context.Context
		requestOption *connector.RequestOption
		response      any
	}

	type test struct {
		fields  fields
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"success": func(t *testing.T) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				requestOption: &connector.RequestOption{
					Method: http.MethodGet,
					URL:    "/files",
				},
			}

			return test{
				fields: fields{
					openai: &connector.OpenAI{
						BaseURL: os.Getenv("OPENAI_API_URL"),
						Header: map[string]string{
							"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")),
						},
					},
				},
				args:    args,
				wantErr: nil,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := testFn(t)

			sut := sut(tt.fields)

			var result connector.OpenAIFile
			err := sut.Send(tt.args.ctx, tt.args.requestOption, &result)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			log.Info(result)
		})
	}
}
