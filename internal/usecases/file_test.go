package usecases_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestFileUsecase_GetListFile(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type test struct {
		fields  fileFields
		args    args
		want    *[]repository.File
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Get List File, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
			}

			mockConnector := connector.NewGoMockConnector(ctrl)

			httpRequestOption := &connector.RequestOption{
				Method: http.MethodGet,
				URL:    "/files",
			}

			expected := []repository.File{
				{
					ID:        "file-1",
					Object:    "file",
					CreatedAt: 1234567890,
					Bytes:     1234567890,
					Status:    "uploaded",
					Filename:  "file-1",
				},
			}

			var result *connector.OpenAIFile
			mockConnector.EXPECT().Send(args.ctx, httpRequestOption, &result).Return(nil).SetArg(2, &connector.OpenAIFile{
				Data: expected,
			})

			return test{
				fields: fileFields{
					connector: mockConnector,
				},
				args:    args,
				want:    &expected,
				wantErr: nil,
			}
		},
	}

	for name, testFn := range tests {

		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := fileSut(tt.fields)

			got, err := sut.GetListFile(tt.args.ctx)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
