//go:build !embed_frontend
// +build !embed_frontend

package main

import "embed"

// Empty FS for dev
var embeddedFrontendFS embed.FS
