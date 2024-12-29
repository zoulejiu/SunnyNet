//go:build !mini
// +build !mini

package Resource

import (
	"embed"
	_ "embed"
	"strings"
)

//go:embed SunnyNetScriptEdit/assets
var frontendAssets embed.FS

//go:embed SunnyNetScriptEdit/index.html
var FrontendIndex []byte

func ReadVueFile(name string) ([]byte, error) {
	fullPath := "SunnyNetScriptEdit/" + name
	if strings.HasPrefix(name, "/") {
		fullPath = "SunnyNetScriptEdit" + name
	}
	return frontendAssets.ReadFile(fullPath)
}
