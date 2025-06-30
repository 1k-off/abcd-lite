package main

import (
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server"
	"github.com/gofiber/storage/badger/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg.Database.GeoIPFS = embeddedGeoIPFS
	cfg.CheckServerCountryBlock()

	storage := badger.New(badger.Config{
		Database: cfg.Database.Path + "/abcd.db",
		Reset:    false,
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		storage.Close()
		os.Exit(0)
	}()

	var staticFS fs.FS
	if cfg.App.Env == config.AppEnvProduction {
		staticFS, _ = fs.Sub(embeddedFrontendFS, "frontend/dist")
	}

	app := server.NewServer(server.Config{
		Storage:         storage,
		Env:             cfg.App.Env,
		AllowedOrigins:  cfg.App.AllowedOrigins,
		AdminTokenHash:  cfg.App.AdminToken,
		JwtSecret:       cfg.App.JwtSecret,
		StaticFS:        staticFS,
		GeoIPDB:         cfg.GetGeoIPDB(),
		DeniedCountries: config.DeniedCountries,
	})
	log.Fatal(app.Listen(":" + cfg.App.Port))
}
