//go:build !embed_files
// +build !embed_files

package main

import "embed"

// Empty FS for dev
var embeddedFrontendFS embed.FS
var embeddedGeoIPFS embed.FS
