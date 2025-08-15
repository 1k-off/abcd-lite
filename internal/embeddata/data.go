package embeddata

import "embed"

type Config struct {
	GeoIPFS     embed.FS
	FrontendFS  embed.FS
	IISConfigFS embed.FS
}

func New(geoIPFS, frontendFS, iisConfigFS embed.FS) *Config {
	return &Config{
		GeoIPFS:     geoIPFS,
		FrontendFS:  frontendFS,
		IISConfigFS: iisConfigFS,
	}
}
