package server

import (
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/handlers"
	jwtware "github.com/1k-off/abcd-lite/internal/server/middleware/jwt"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/badger/v2"
)

type Config struct {
	Storage        *badger.Storage
	Env            string
	AllowedOrigins []string
	AdminTokenHash string
	JwtSecret      string
}

// NewServer creates a new Fiber app and sets up the routes.
func NewServer(c Config) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "abcd-lite",
	})
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     c.AllowedOrigins,
		AllowCredentials: true,
	}))

	if c.Env == "production" {
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

	projectService := services.NewProjectService(c.Storage)
	iisDeploymentService := services.NewIISDeploymentService(projectService)

	deploy := app.Group("/deploy")
	deploy.Post("/iis", handlers.IISDeploy(iisDeploymentService))

	jwtConfig := jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(c.JwtSecret),
		},
		TokenLookup: "cookie:jwt,header:Authorization",
		AuthScheme:  "Bearer",
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.DefaultErrorResponse{
				Message: "Unauthorized",
				Error:   err.Error(),
			})
		},
	}

	app.Post("/login", handlers.Login(c.AdminTokenHash, c.JwtSecret, c.Env))
	app.Post("/logout", handlers.Logout(c.Env))

	api := app.Group("/api", jwtware.New(jwtConfig))

	projects := api.Group("/projects")
	projects.Get("/", handlers.GetProjects(projectService))
	projects.Post("/", handlers.CreateProject(projectService))
	projects.Put("/:id", handlers.UpdateProject(projectService))
	projects.Delete("/:id", handlers.DeleteProject(projectService))
	projects.Post("/:id/api-keys", handlers.AddAPIKey(projectService))
	projects.Delete("/:id/api-keys/:keyId", handlers.DeleteAPIKey(projectService))

	return app
}
