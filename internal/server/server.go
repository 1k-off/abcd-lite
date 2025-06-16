package server

import (
	"github.com/1k-off/abcd-lite/internal/server/handlers"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/badger/v2"
)

// NewServer creates a new Fiber app and sets up the routes.
func NewServer(storage *badger.Storage, env string) *fiber.App {
	app := fiber.New()

	// Configure CORS only in development
	if env != "production" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:5173",
			AllowHeaders: "Origin, Content-Type, Accept",
			AllowMethods: "GET, POST, PUT, DELETE",
		}))
	}

	projectService := services.NewProjectService(storage)

	// API routes
	api := app.Group("/api")
	projects := api.Group("/projects")
	projects.Get("/", handlers.GetProjects(projectService))
	projects.Post("/", handlers.CreateProject(projectService))
	projects.Put("/:id", handlers.UpdateProject(projectService))
	projects.Delete("/:id", handlers.DeleteProject(projectService))

	// Health check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Serve static files in production
	if env == "production" {
		// Serve static files from the embedded frontend build
		app.Static("/", "./frontend/dist")
		// Handle SPA routing by serving index.html for all other routes
		app.Get("*", func(c *fiber.Ctx) error {
			return c.SendFile("./frontend/dist/index.html")
		})
	}

	return app
}
