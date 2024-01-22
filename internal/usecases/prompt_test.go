package usecases_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPromptUsecase_SendPrompt(t *testing.T) {
	type args struct {
		ctx     context.Context
		message string
	}

	type test struct {
		fields  promptFields
		args    args
		want    *repository.Message
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Send Prompt, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:     ctx,
				message: "Hello World",
			}

			mockConnector := connector.NewGoMockConnector(ctrl)

			mockConnector.EXPECT().Send(args.ctx, gomock.Any(), gomock.Any()).Return(nil).SetArg(2, &connector.OpenAIThread{
				FirstID: "thread-1",
			})

			mockConnector.EXPECT().Send(args.ctx, gomock.Any(), gomock.Any()).Return(nil).SetArg(2, &repository.Message{
				ID: "message-1",
			})

			mockConnector.EXPECT().Send(args.ctx, gomock.Any(), gomock.Any()).Return(nil).SetArg(2, &connector.OpenAIRun{
				ID: "run-1",
			})

			mockConnector.EXPECT().Send(args.ctx, gomock.Any(), gomock.Any()).Return(nil).SetArg(2, &connector.OpenAIRun{
				ID:     "run-1",
				Status: "completed",
			})

			expected := []repository.Message{
				{
					ID:        "message-1",
					Object:    "message",
					CreatedAt: 1234567890,
					ThreadID:  "thread-1",
					Content: []repository.ContentMessage{
						{
							Type: "text",
							Text: repository.TextContent{
								Value: "Hello World",
							},
						},
					},
				},
			}

			mockConnector.EXPECT().Send(args.ctx, gomock.Any(), gomock.Any()).Return(nil).SetArg(2, &connector.OpenAIMessage{
				Data: expected,
			})

			return test{
				fields: promptFields{
					connector: mockConnector,
				},
				args:    args,
				want:    &expected[0],
				wantErr: nil,
			}
		},
	}

	for name, testFn := range tests {

		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := promptSut(tt.fields)

			got, err := sut.SendPrompt(tt.args.ctx, tt.args.message)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
