package server

import (
	"github.com/1k-off/abcd-lite/internal/server/handlers"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/badger/v2"
)

// NewServer creates a new Fiber app and sets up the routes.
func NewServer(storage *badger.Storage, env string) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "abcd-lite",
		BodyLimit:     5 * 1024 * 1024 * 1024,
	})
	app.Use(recover.New())

	if env == "production" {
		app.Get("/*", static.New("./frontend/dist"))
		app.Use("/static", static.New("./frontend/dist/index.html"))
		app.Use(favicon.New(favicon.Config{
			File: "./frontend/dist/favicon.ico",
			URL:  "/favicon.ico",
		}))
		app.Use(compress.New())
	}

	app.Get(healthcheck.DefaultLivenessEndpoint, healthcheck.NewHealthChecker())

	app.Get("/healthz", healthcheck.NewHealthChecker())

	projectService := services.NewProjectService(storage)

	api := app.Group("/api")
	projects := api.Group("/projects")
	projects.Get("/", handlers.GetProjects(projectService))
	projects.Post("/", handlers.CreateProject(projectService))
	projects.Put("/:id", handlers.UpdateProject(projectService))
	projects.Delete("/:id", handlers.DeleteProject(projectService))

	return app
}
