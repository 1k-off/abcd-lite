package server

import (
	"fmt"
	"io/fs"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/1k-off/abcd-lite/internal/server/handlers"
	"github.com/1k-off/abcd-lite/internal/server/middleware"
	jwtware "github.com/1k-off/abcd-lite/internal/server/middleware/jwt"
	"github.com/1k-off/abcd-lite/internal/server/services"
	"github.com/1k-off/abcd-lite/internal/storage"
	"github.com/1k-off/abcd-lite/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/oschwald/geoip2-golang"
)

type Config struct {
	ProjectStorage  storage.ProjectStorage
	SettingsStorage storage.SettingsStorage
	Env             string
	AllowedOrigins  []string
	StaticFS        fs.FS
	GeoIPDB         *geoip2.Reader
	DeniedCountries map[string]bool
	PackageLimits   config.PackageLimits
}

// NewServer creates a new Fiber app and sets up the routes.
func NewServer(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "abcd-lite",
	})

	app.Use(middleware.CountryBlockMiddleware(cfg.GeoIPDB, cfg.DeniedCountries))

	projectService := services.NewProjectService(cfg.ProjectStorage, cfg.PackageLimits)
	settingsService := services.NewSettingsService(cfg.SettingsStorage)
	iisDeploymentService := services.NewIISDeploymentService(projectService)

	// Check if admin token is set
	adminToken, err := settingsService.GetAdminToken()
	if err != nil || adminToken == "" {
		log.Fatal("Admin token is not set. Please run `abcd-lite config generate`")
	}

	// Get JWT secret from settings service
	jwtSecret, err := settingsService.GetJwtSecret()
	if err != nil {
		// If JWT secret doesn't exist in datastore, generate a new one
		jwtSecret = util.GenerateRandomToken(32)
		if err := settingsService.SetJwtSecret(jwtSecret); err != nil {
			panic(fmt.Sprintf("Failed to set JWT secret: %v", err))
		}
	}

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Next: func(c fiber.Ctx) bool {
			return c.Path() == "/api/auth/status"
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowCredentials: true,
	}))

	if cfg.Env == config.AppEnvProduction {
		app.Use("/*", static.New("", static.Config{
			FS: cfg.StaticFS,
		}))
		app.Get("/", func(c fiber.Ctx) error {
			return c.SendFile("dist/index.html")
		})
		app.Use(compress.New())
	}

	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get("/healthz", healthcheck.New())

	deploy := app.Group("/deploy")
	deploy.Post("/iis", handlers.IISDeploy(iisDeploymentService))

	jwtConfig := jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(jwtSecret),
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

	app.Post("/login", handlers.Login(settingsService, cfg.Env))
	app.Post("/logout", handlers.Logout(cfg.Env))

	api := app.Group("/api", jwtware.New(jwtConfig))
	api.Get("/auth/status", handlers.AuthStatus(jwtSecret))

	projects := api.Group("/projects")
	projects.Get("/", handlers.GetProjects(projectService))
	projects.Post("/", handlers.CreateProject(projectService))
	projects.Put("/:id", handlers.UpdateProject(projectService))
	projects.Delete("/:id", handlers.DeleteProject(projectService))
	projects.Post("/:id/api-keys", handlers.AddAPIKey(projectService))
	projects.Delete("/:id/api-keys/:keyId", handlers.DeleteAPIKey(projectService))

	return app
}
