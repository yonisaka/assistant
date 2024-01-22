package request

type Prompt struct {
	Message string `json:"message" validate:"required"`
}
