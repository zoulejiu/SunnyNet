//go:build mini
// +build mini

package SunnyNet

import (
	"github.com/qtgolang/SunnyNet/src/http"
)

// 是否是用户自定义脚本编辑请求
func (s *proxyRequest) isUserScriptCodeEditRequest(request *http.Request) bool {
	return false
}

// SetScriptPage 设置脚本页面
func (s *Sunny) SetScriptPage(Page string) string {
	return "no"
}
