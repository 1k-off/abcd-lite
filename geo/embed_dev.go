//go:build !production

package geo

import "embed"

// Empty FS for development/testing
var geoIPFS embed.FS

func GetGeoIPFS() embed.FS {
	return geoIPFS
}
