package httphandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yonisaka/assistant/internal/usecases"
)

type fileHandler struct {
	fileUsecase usecases.FileUsecase
}

func NewFileHandler(fileUsecase usecases.FileUsecase) FileHandler {
	return &fileHandler{
		fileUsecase: fileUsecase,
	}
}

type FileHandler interface {
	GetListFile(c *fiber.Ctx) error
}

func (h *fileHandler) GetListFile(c *fiber.Ctx) error {
	result, err := h.fileUsecase.GetListFile(c.Context())
	if err != nil {
		log.Warn(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(result)
}
