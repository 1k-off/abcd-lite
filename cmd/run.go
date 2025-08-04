package cmd

import (
	"context"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1k-off/abcd-lite/internal/config"
	"github.com/1k-off/abcd-lite/internal/server"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the ABCD Lite server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		cfg.Database.GeoIPFS = embeddedData.GetGeoIPFS()
		cfg.CheckServerCountryBlock()

		limits, licenseDeactivationFunc, err := config.Limits(context.Background(), cfg.App.LicenseKey, cfg.App.LicenseKey)
		if err != nil {
			log.Printf("Failed to load license. Personal edition will be used. Error: %v", err)
		}

		cfg.App.DeactivationFunc = licenseDeactivationFunc

		storage := cfg.Storage()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Println("Shutting down...")
			if licenseDeactivationFunc != nil {
				err := licenseDeactivationFunc()
				if err != nil {
					log.Printf("Failed to deactivate license. Error: %v", err)
				}
			}
			storage.Close()
			os.Exit(0)
		}()

		var staticFS fs.FS
		if cfg.App.Env == config.AppEnvProduction {
			staticFS, _ = fs.Sub(embeddedData.GetFrontendFS(), "frontend/dist")
		} else {
			staticFS = os.DirFS("frontend/dist")
		}

		app := server.NewServer(server.Config{
			Storage:         storage,
			Env:             cfg.App.Env,
			AllowedOrigins:  cfg.App.AllowedOrigins,
			StaticFS:        staticFS,
			GeoIPDB:         cfg.GetGeoIPDB(),
			DeniedCountries: config.DeniedCountries,
			PackageLimits:   limits,
		})
		log.Fatal(app.Listen(":" + cfg.App.Port))
	},
}
