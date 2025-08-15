//go:build !production

package frontend

import "embed"

// Empty FS for development/testing
var distFS embed.FS

func GetDistFS() embed.FS {
	return distFS
}
