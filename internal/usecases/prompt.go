package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
)

// SendPrompt is a function to send prompt to OpenAI with several steps below:
// 1. Create Thread
// 2. Create Message
// 3. Run Thread
// 4. Run Status Thread
// 5. Get Prompt Response
func (u *promptUsecase) SendPrompt(ctx context.Context, message string) (*repository.Message, error) {
	threadID, err := u.createThread(ctx)
	if err != nil {
		return nil, err
	}

	err = u.createMessage(ctx, threadID, message)
	if err != nil {
		return nil, err
	}

	runID, err := u.runThread(ctx, threadID)
	if err != nil {
		return nil, err
	}

	err = u.runStatus(ctx, threadID, runID)
	if err != nil {
		return nil, err
	}

	promptResponse, err := u.getPromptResponse(ctx, threadID)
	if err != nil {
		return nil, err
	}

	return promptResponse, err
}

// createThread is a function to create thread in OpenAI
// It will check if last thread is still active or not
// If active, it will return the last thread ID
// If not active, it will create new thread and return the new thread ID
func (u *promptUsecase) createThread(ctx context.Context) (string, error) {
	httpRequestOption := &connector.RequestOption{
		Method: http.MethodGet,
		URL:    "/threads",
		CustomHeader: map[string]string{
			connector.HeaderKeyOpenAIBeta: connector.AssistantV1,
		},
	}

	var result *connector.OpenAIThread
	if err := u.connector.Send(ctx, httpRequestOption, &result); err != nil {
		return "", err
	}

	if result == nil {
		return "", connector.ErrCreateThread
	}

	log.Infow("Thread Created:", "id", result.FirstID)

	return result.FirstID, nil
}

// createMessage is a function to create message in OpenAI
// It will create message in the last thread from createThread function
func (u *promptUsecase) createMessage(ctx context.Context, threadID, message string) error {
	requestBodyMessage := connector.RequestMessage{
		Role:    "user",
		Content: message,
	}

	var bufMessage bytes.Buffer
	if err := json.NewEncoder(&bufMessage).Encode(requestBodyMessage); err != nil {
		return err
	}

	httpRequestOption := &connector.RequestOption{
		Method: http.MethodPost,
		URL:    fmt.Sprintf("/threads/%s/messages", threadID),
		CustomHeader: map[string]string{
			connector.HeaderKeyOpenAIBeta: connector.AssistantV1,
		},
		Body: &bufMessage,
	}

	var result *repository.Message
	if err := u.connector.Send(ctx, httpRequestOption, &result); err != nil {
		return err
	}

	if result == nil {
		return connector.ErrCreateMessage
	}

	log.Infow("Message Created:", "id", result.ID, "content", result.Content)

	return nil
}

// runThread is a function to run thread in OpenAI
// It will run the last thread from createThread function
// Using AssistantID that has been set on openAI before
func (u *promptUsecase) runThread(ctx context.Context, threadID string) (string, error) {
	requestBodyRun := connector.RequestRun{
		AssistantID: os.Getenv("OPENAI_ASSISTANT_ID"),
	}

	var bufRun bytes.Buffer
	if err := json.NewEncoder(&bufRun).Encode(requestBodyRun); err != nil {
		return "", err
	}

	httpRequestOption := &connector.RequestOption{
		Method: http.MethodPost,
		URL:    fmt.Sprintf("/threads/%s/runs", threadID),
		CustomHeader: map[string]string{
			connector.HeaderKeyOpenAIBeta: connector.AssistantV1,
		},
		Body: &bufRun,
	}

	var result *connector.OpenAIRun
	if err := u.connector.Send(ctx, httpRequestOption, &result); err != nil {
		return "", err
	}

	if result == nil {
		return "", connector.ErrRunThread
	}

	log.Infow("Run Created:", "id", result.ID, "status", result.Status)

	return result.ID, nil
}

// runStatus is a function to check run status in OpenAI
// It will check the last run status from runThread function
// If the status is not queued or in progress, it will return the last run status
// If the status is queued or in progress, it will check the status again until it's not queued or in progress with 500ms delay
func (u *promptUsecase) runStatus(ctx context.Context, threadID, runID string) error {
	httpRequestOption := &connector.RequestOption{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/threads/%s/runs/%s", threadID, runID),
		CustomHeader: map[string]string{
			connector.HeaderKeyOpenAIBeta: connector.AssistantV1,
		},
	}

	runStep := 1
	for { //nolint: wsl
		var result *connector.OpenAIRun
		if err := u.connector.Send(ctx, httpRequestOption, &result); err != nil {
			return err
		}

		if result == nil {
			return connector.ErrRunStatusThread
		}

		log.Infow("Run Status:", "id", result.ID, "status", result.Status, "step", runStep)

		if result.Status != "" &&
			result.Status != connector.OpenAIStatusQueued &&
			result.Status != connector.OpenAIStatusInProgress {
			break
		}

		runStep++

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// getPromptResponse is a function to get prompt response in OpenAI
// It will get the last prompt response from runStatus function
// We only need the last prompt response
func (u *promptUsecase) getPromptResponse(ctx context.Context, threadID string) (*repository.Message, error) {
	httpRequestOption := &connector.RequestOption{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/threads/%s/messages", threadID),
		CustomHeader: map[string]string{
			connector.HeaderKeyOpenAIBeta: connector.AssistantV1,
		},
	}

	var result *connector.OpenAIMessage
	if err := u.connector.Send(ctx, httpRequestOption, &result); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, connector.ErrGetPrompt
	}

	log.Infow("Prompt Last Response:", "id", result.Data[0].ID)

	return &result.Data[0], nil
}
