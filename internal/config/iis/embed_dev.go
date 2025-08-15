//go:build !production

package iis

import "embed"

// Empty FS for development/testing
var configFS embed.FS

func GetConfigFS() embed.FS {
	return configFS
}
