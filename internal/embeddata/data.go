package embeddata

import "embed"

type Config struct {
	GeoIPFS    embed.FS
	FrontendFS embed.FS
}

type EmbedData interface {
	GetGeoIPFS() embed.FS
	GetFrontendFS() embed.FS
}

func New(geoIPFS, frontendFS embed.FS) EmbedData {
	return &Config{
		GeoIPFS:    geoIPFS,
		FrontendFS: frontendFS,
	}
}

func (c *Config) GetGeoIPFS() embed.FS {
	return c.GeoIPFS
}

func (c *Config) GetFrontendFS() embed.FS {
	return c.FrontendFS
}
