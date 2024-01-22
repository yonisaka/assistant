package httphandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yonisaka/assistant/internal/entities/request"
	"github.com/yonisaka/assistant/internal/usecases"
)

type promptHandler struct {
	promptUsecase usecases.PromptUsecase
}

func NewPromptHandler(promptUsecase usecases.PromptUsecase) PromptHandler {
	return &promptHandler{
		promptUsecase: promptUsecase,
	}
}

type PromptHandler interface {
	SendPrompt(c *fiber.Ctx) error
}

func (h *promptHandler) SendPrompt(c *fiber.Ctx) error {
	prompt := new(request.Prompt)

	if err := c.BodyParser(prompt); err != nil {
		log.Warn(err)
		return fiber.ErrBadRequest
	}

	result, err := h.promptUsecase.SendPrompt(c.Context(), prompt.Message)
	if err != nil {
		log.Warn(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(result)
}
