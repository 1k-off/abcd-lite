//go:build embed_frontend
// +build embed_frontend

package main

import "embed"

//go:embed frontend/dist/*
var embeddedFrontendFS embed.FS
