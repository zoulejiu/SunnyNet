//go:build !mini
// +build !mini

package Resource

import (
	"embed"
	_ "embed"
	"strings"
)

//go:embed Script/dist/assets
var frontendAssets embed.FS

//go:embed Script/dist/index.html
var FrontendIndex []byte

func ReadVueFile(name string) ([]byte, error) {
	fullPath := "Script/dist/" + name
	if strings.HasPrefix(name, "/") {
		fullPath = "Script/dist" + name
	}
	return frontendAssets.ReadFile(fullPath)
}
