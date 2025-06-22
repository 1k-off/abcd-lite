package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/deps"
	"github.com/1k-off/abcd-lite/internal/server"
	"github.com/gofiber/storage/badger/v2"
)

func init() {
	if err := deps.GetDependencies(); err != nil {
		log.Fatalf("Failed to get dependencies: %v", err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	storage := badger.New(badger.Config{
		Database: cfg.Database.Path + "/abcd.db",
		Reset:    false,
	})

	// ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		storage.Close()
		// cancel()
		os.Exit(0)
	}()

	app := server.NewServer(storage, cfg.App.Env)
	log.Fatal(app.Listen(":" + cfg.App.Port))
}
