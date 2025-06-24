package handlers

import (
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/messages"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func GetProjects(service services.ProjectService) fiber.Handler {
	return func(c fiber.Ctx) error {
		projects, err := service.GetProjects()
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToGetProjects,
				Error:   err.Error(),
			})
		}
		return c.JSON(domain.ProjectsResponse{
			Projects: projects,
		})
	}
}

func CreateProject(service services.ProjectService) fiber.Handler {
	return func(c fiber.Ctx) error {
		project := new(domain.Project)
		log.Debug("Called CreateProject with body: ", string(c.Body()))
		if err := c.Bind().Body(project); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrInvalidRequestBody,
				Error:   err.Error(),
			})
		}
		createdProject, err := service.CreateProject(*project)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToCreateProject,
				Error:   err.Error(),
			})
		}
		log.Info("Created project: ", createdProject.Name)
		return c.Status(fiber.StatusCreated).JSON(domain.ProjectResponse{
			Project: createdProject,
		})
	}
}

func UpdateProject(service services.ProjectService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		project := new(domain.Project)
		if err := c.Bind().Body(project); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrInvalidRequestBody,
				Error:   err.Error(),
			})
		}
		project.ID = id
		if err := service.UpdateProject(*project); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToUpdateProject,
				Error:   err.Error(),
			})
		}
		log.Info("Updated project: ", project.Name)
		return c.Status(fiber.StatusOK).JSON(domain.ProjectResponse{
			Project: *project,
		})
	}
}

func DeleteProject(service services.ProjectService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		project, err := service.GetProject(id)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToGetProject,
				Error:   err.Error(),
			})
		}
		if err := service.DeleteProject(id); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: messages.ErrFailedToDeleteProject,
				Error:   err.Error(),
			})
		}
		log.Info("Deleted project: ", project.Name)
		return c.Status(fiber.StatusOK).JSON(domain.DefaultErrorResponse{
			Message: messages.MsgProjectDeleted,
		})
	}
}
