//go:build production

package frontend

import "embed"

//go:embed dist
var distFS embed.FS

func GetDistFS() embed.FS {
	return distFS
}
