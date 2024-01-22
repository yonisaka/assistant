package repository

// File is a struct of file from OpenAI API
type File struct {
	Object        string      `json:"object"`
	ID            string      `json:"id"`
	Purpose       string      `json:"purpose"`
	Filename      string      `json:"filename"`
	Bytes         int64       `json:"bytes"`
	CreatedAt     int64       `json:"created_at"`
	Status        string      `json:"status"`
	StatusDetails interface{} `json:"status_details"`
}

type Message struct {
	ID          string           `json:"id"`
	Object      string           `json:"object"`
	CreatedAt   int64            `json:"created_at"`
	ThreadID    string           `json:"thread_id"`
	Role        string           `json:"role"`
	Content     []ContentMessage `json:"content"`
	FileIDS     []interface{}    `json:"file_ids"`
	AssistantID string           `json:"assistant_id"`
	RunID       string           `json:"run_id"`
	Metadata    interface{}      `json:"metadata"`
}

type ContentMessage struct {
	Type string      `json:"type"`
	Text TextContent `json:"text"`
}

type TextContent struct {
	Value       string        `json:"value"`
	Annotations []interface{} `json:"annotations"`
}

type Thread struct {
	ID        string      `json:"id"`
	Object    string      `json:"object"`
	CreatedAt int64       `json:"created_at"`
	Metadata  interface{} `json:"metadata"`
}
