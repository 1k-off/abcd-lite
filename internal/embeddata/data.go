package embeddata

import "embed"

type Config struct {
	GeoIPFS     embed.FS
	FrontendFS  embed.FS
	IISConfigFS embed.FS
}

type EmbedData interface {
	GetGeoIPFS() embed.FS
	GetFrontendFS() embed.FS
	GetIISConfigFS() embed.FS
}

func New(geoIPFS, frontendFS, iisConfigFS embed.FS) EmbedData {
	return &Config{
		GeoIPFS:     geoIPFS,
		FrontendFS:  frontendFS,
		IISConfigFS: iisConfigFS,
	}
}

func (c *Config) GetGeoIPFS() embed.FS {
	return c.GeoIPFS
}

func (c *Config) GetFrontendFS() embed.FS {
	return c.FrontendFS
}

func (c *Config) GetIISConfigFS() embed.FS {
	return c.IISConfigFS
}
