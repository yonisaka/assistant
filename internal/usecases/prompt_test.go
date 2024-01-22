package usecases_test

import (
	"context"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"testing"
)

func TestOpenAIUsecase_GetListFile(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type test struct {
		fields  promptFields
		args    args
		want    *[]repository.File
		wantErr bool
	}
}
