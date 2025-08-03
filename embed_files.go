//go:build production
// +build production

package main

import "embed"

//go:embed frontend/dist/*
var embeddedFrontendFS embed.FS

//go:embed data/geo/GeoLite2-Country.mmdb
var embeddedGeoIPFS embed.FS

//go:embed internal/config/iis/*.config
var embeddedIISConfigFS embed.FS
