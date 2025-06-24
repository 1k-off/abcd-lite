package handlers

import (
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/messages"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func IISDeploy(s services.IISDeploymentService) fiber.Handler {
	return func(c fiber.Ctx) error {
		iis := new(domain.IIS)
		if err := c.Bind().Body(iis); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrInvalidRequestBody,
				Error:   err.Error(),
			})
		}
		if err := s.Deploy(*iis); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToDeployIIS,
				Error:   err.Error(),
			})
		}
		return c.JSON(domain.DefaultInfoResponse{
			Message: messages.MsgIISDeployed,
		})
	}
}
