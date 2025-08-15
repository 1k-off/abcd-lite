package geo

import "embed"

//go:embed GeoLite2-Country.mmdb
var geoIPFS embed.FS

func GetGeoIPFS() embed.FS {
	return geoIPFS
}
