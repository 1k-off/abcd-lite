package handlers

import (
	"encoding/json"
	"time"

	"github.com/1k-off/abcd/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/badger/v2"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	storage *badger.Storage
}

func NewProjectHandler(storage *badger.Storage) *ProjectHandler {
	return &ProjectHandler{storage: storage}
}

func (h *ProjectHandler) GetProjects(c *fiber.Ctx) error {
	projects := make([]models.Project, 0)

	// Get all projects from storage
	keys, err := h.storage.Get("project:keys")
	if err != nil {
		// If no keys exist yet, return empty array
		return c.JSON(projects)
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		// If keys are invalid, return empty array
		return c.JSON(projects)
	}

	for _, key := range projectKeys {
		data, err := h.storage.Get("project:" + key)
		if err != nil {
			continue
		}

		var project models.Project
		if err := json.Unmarshal(data, &project); err != nil {
			continue
		}
		projects = append(projects, project)
	}

	return c.JSON(projects)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate ID and timestamps
	project.ID = uuid.New().String()
	project.CreatedAt = time.Now().Format(time.RFC3339)
	project.UpdatedAt = project.CreatedAt

	// Store project
	data, err := json.Marshal(project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create project",
		})
	}

	if err := h.storage.Set("project:"+project.ID, data, 0); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store project",
		})
	}

	// Update project keys list
	keys, err := h.storage.Get("project:keys")
	if err != nil {
		keys = []byte("[]")
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		projectKeys = []string{}
	}

	projectKeys = append(projectKeys, project.ID)
	keysData, err := json.Marshal(projectKeys)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update project keys",
		})
	}

	if err := h.storage.Set("project:keys", keysData, 0); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store project keys",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get existing project
	existingData, err := h.storage.Get("project:" + id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	var existingProject models.Project
	if err := json.Unmarshal(existingData, &existingProject); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse existing project",
		})
	}

	// Update fields
	project.ID = id
	project.CreatedAt = existingProject.CreatedAt
	project.UpdatedAt = time.Now().Format(time.RFC3339)

	// Store updated project
	data, err := json.Marshal(project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update project",
		})
	}

	if err := h.storage.Set("project:"+id, data, 0); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store updated project",
		})
	}

	return c.JSON(project)
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	if err := h.storage.Delete("project:" + id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete project",
		})
	}

	// Update project keys list
	keys, err := h.storage.Get("project:keys")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get project keys",
		})
	}

	var projectKeys []string
	if err := json.Unmarshal(keys, &projectKeys); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse project keys",
		})
	}

	// Remove the deleted project ID from the keys list
	newKeys := make([]string, 0)
	for _, key := range projectKeys {
		if key != id {
			newKeys = append(newKeys, key)
		}
	}

	keysData, err := json.Marshal(newKeys)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update project keys",
		})
	}

	if err := h.storage.Set("project:keys", keysData, 0); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store project keys",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
