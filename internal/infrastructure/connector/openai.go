package connector

import (
	"errors"

	"github.com/yonisaka/assistant/internal/entities/repository"
)

// OpenAI is a struct to set OpenAI API
type OpenAI struct {
	BaseURL string
	Header  map[string]string
}

// OpenAIFile is a struct to get list file from OpenAI API
type OpenAIFile struct {
	Object  string            `json:"object"`
	HasMore bool              `json:"has_more"`
	Data    []repository.File `json:"data"`
}

// OpenAIMessage is a struct to get response prompt
type OpenAIMessage struct {
	Object  string               `json:"object"`
	Data    []repository.Message `json:"data"`
	FirstID string               `json:"first_id"`
	LastID  string               `json:"last_id"`
	HasMore bool                 `json:"has_more"`
}

// OpenAIThread is a struct to get thread
type OpenAIThread struct {
	Object  string              `json:"object"`
	Data    []repository.Thread `json:"data"`
	FirstID string              `json:"first_id"`
	LastID  string              `json:"last_id"`
	HasMore bool                `json:"has_more"`
}

type OpenAIRun struct {
	ID           string       `json:"id"`
	Object       string       `json:"object"`
	CreatedAt    int64        `json:"created_at"`
	AssistantID  string       `json:"assistant_id"`
	ThreadID     string       `json:"thread_id"`
	Status       string       `json:"status"`
	StartedAt    interface{}  `json:"started_at"`
	ExpiresAt    int64        `json:"expires_at"`
	CancelledAt  interface{}  `json:"cancelled_at"`
	FailedAt     interface{}  `json:"failed_at"`
	CompletedAt  interface{}  `json:"completed_at"`
	LastError    interface{}  `json:"last_error"`
	Model        string       `json:"model"`
	Instructions string       `json:"instructions"`
	Tools        []OpenAITool `json:"tools"`
	FileIDS      []string     `json:"file_ids"`
	Metadata     interface{}  `json:"metadata"`
}

type OpenAITool struct {
	Type string `json:"type"`
}

var (
	ErrCreateThread    = errors.New("failed to create new thread")
	ErrCreateMessage   = errors.New("failed to create new message")
	ErrRunThread       = errors.New("failed to run thread")
	ErrRunStatusThread = errors.New("failed to run status thread")
	ErrGetPrompt       = errors.New("failed to get prompt response")
)

const (
	OpenAIStatusQueued     = "queued"
	OpenAIStatusInProgress = "in_progress"
)
