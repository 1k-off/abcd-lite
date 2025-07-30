//go:build !production
// +build !production

package main

import "embed"

// Empty FS for dev
var embeddedFrontendFS embed.FS
var embeddedGeoIPFS embed.FS
