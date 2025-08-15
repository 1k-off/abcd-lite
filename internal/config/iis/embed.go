package iis

import "embed"

//go:embed *.config
var configFS embed.FS

func GetConfigFS() embed.FS {
	return configFS
}
