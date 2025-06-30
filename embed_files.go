//go:build embed_files
// +build embed_files

package main

import "embed"

//go:embed frontend/dist/*
var embeddedFrontendFS embed.FS

//go:embed data/geo/GeoLite2-Country.mmdb
var embeddedGeoIPFS embed.FS
