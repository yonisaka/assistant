package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yonisaka/assistant/internal/entities/repository"
	"github.com/yonisaka/assistant/internal/infrastructure/connector"
	"net/http"
	"os"
	"time"
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

	var threadResponse *connector.OpenAIThread
	err := u.connector.Send(ctx, httpRequestOption, &threadResponse)
	if err != nil {
		return "", err
	}

	if threadResponse == nil {
		return "", connector.ErrCreateThread
	}

	log.Infow("Thread Created:", "id", threadResponse.FirstID)

	return threadResponse.FirstID, err
}

// createMessage is a function to create message in OpenAI
// It will create message in the last thread from createThread function
func (u *promptUsecase) createMessage(ctx context.Context, threadID, message string) error {
	requestBodyMessage := connector.RequestMessage{
		Role:    "user",
		Content: message,
	}

	var bufMessage bytes.Buffer
	err := json.NewEncoder(&bufMessage).Encode(requestBodyMessage)
	if err != nil {
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

	var messageResponse *repository.Message
	err = u.connector.Send(ctx, httpRequestOption, &messageResponse)
	if err != nil {
		return err
	}

	if messageResponse == nil {
		return connector.ErrCreateMessage
	}

	log.Infow("Message Created:", "id", messageResponse.ID, "content", messageResponse.Content)

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
	err := json.NewEncoder(&bufRun).Encode(requestBodyRun)
	if err != nil {
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

	var runResponse *connector.OpenAIRun
	err = u.connector.Send(ctx, httpRequestOption, &runResponse)
	if err != nil {
		return "", err
	}

	if runResponse == nil {
		return "", connector.ErrRunThread
	}

	log.Infow("Run Created:", "id", runResponse.ID, "status", runResponse.Status)

	return runResponse.ID, err
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
	for {
		var runStatusResponse *connector.OpenAIRun
		err := u.connector.Send(ctx, httpRequestOption, &runStatusResponse)
		if err != nil {
			return err
		}

		if runStatusResponse == nil {
			return connector.ErrRunStatusThread
		}

		log.Infow("Run Status:", "id", runStatusResponse.ID, "status", runStatusResponse.Status, "step", runStep)

		if runStatusResponse.Status != "" &&
			runStatusResponse.Status != connector.OpenAIStatusQueued &&
			runStatusResponse.Status != connector.OpenAIStatusInProgress {
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

	var promptResponse *connector.OpenAIMessage
	err := u.connector.Send(ctx, httpRequestOption, &promptResponse)
	if err != nil {
		return nil, err
	}

	if promptResponse == nil {
		return nil, connector.ErrGetPrompt
	}

	log.Infow("Prompt Last Response:", "id", promptResponse.Data[0].ID)

	return &promptResponse.Data[0], err
}
