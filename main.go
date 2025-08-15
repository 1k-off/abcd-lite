package main

import (
	"log"
	"os"

	"github.com/1k-off/abcd-lite/cmd"
	"github.com/1k-off/abcd-lite/frontend"
	"github.com/1k-off/abcd-lite/geo"
	"github.com/1k-off/abcd-lite/internal/config/iis"
	"github.com/1k-off/abcd-lite/internal/embeddata"
)

var (
	embeddedFrontendFS  = frontend.GetDistFS()
	embeddedGeoIPFS     = geo.GetGeoIPFS()
	embeddedIISConfigFS = iis.GetConfigFS()
)

func main() {
	data := embeddata.New(embeddedGeoIPFS, embeddedFrontendFS, embeddedIISConfigFS)
	cmd.SetEmbeddedData(data)
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
