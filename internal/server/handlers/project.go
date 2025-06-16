package handlers

import (
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/messages"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v2"
)

func GetProjects(service services.ProjectService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		projects, err := service.GetProjects()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrFailedToGetProjects,
				ErrorFull: err,
			})
		}

		return c.JSON(domain.ProjectsResponse{
			Projects: projects,
		})
	}
}

func CreateProject(service services.ProjectService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var project domain.Project
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrInvalidRequestBody,
				ErrorFull: err,
			})
		}
		createdProject, err := service.CreateProject(project)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrFailedToCreateProject,
				ErrorFull: err,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(domain.ProjectResponse{
			Project: createdProject,
		})
	}
}

func UpdateProject(service services.ProjectService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var project domain.Project
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrInvalidRequestBody,
				ErrorFull: err,
			})
		}
		project.ID = id
		if err := service.UpdateProject(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrFailedToUpdateProject,
				ErrorFull: err,
			})
		}
		return c.Status(fiber.StatusOK).JSON(domain.ProjectResponse{
			Project: project,
		})
	}
}

func DeleteProject(service services.ProjectService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.DeleteProject(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ProjectsErrorResponse{
				Error:     messages.ErrFailedToDeleteProject,
				ErrorFull: err,
			})
		}
		return c.Status(fiber.StatusOK).JSON(domain.ProjectInfoResponse{
			Message: messages.MsgProjectDeleted,
		})
	}
}
