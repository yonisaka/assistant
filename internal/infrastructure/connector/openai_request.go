package connector

type (
	RequestMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	RequestRun struct {
		AssistantID  string `json:"assistant_id"`
		Instructions string `json:"instructions"`
	}
)
