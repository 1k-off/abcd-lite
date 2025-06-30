package main

import (
	"log"
	"os"

	"github.com/1k-off/abcd-lite/cmd"
	"github.com/1k-off/abcd-lite/internal/embeddata"
)

func main() {
	data := embeddata.New(embeddedGeoIPFS, embeddedFrontendFS)
	cmd.SetEmbeddedData(data)
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
