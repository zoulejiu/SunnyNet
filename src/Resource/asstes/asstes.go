//go:build windows
// +build windows

package asstes

import (
	"embed"
	"fmt"
	"github.com/gorilla/websocket"
	"io/fs"
	"log"
	"net"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨源请求
	},
}

//go:embed vs
var Assets embed.FS

//go:embed code_temp.html
var CodeTemp []byte

//go:embed vs/codicon.ttf
var codicon []byte

//go:embed vs/index.70d8fe28.js
var indexJs []byte

//go:embed vs/index.460da0d4.css
var indexCSS []byte

func Listen(handleWebSocket func(w http.ResponseWriter, r *http.Request)) int {
	vsFS, err := fs.Sub(Assets, "vs")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(vsFS))
	http.Handle("/vs/", http.StripPrefix("/vs/", fileServer))
	http.HandleFunc("/vs/base/browser/ui/codicons/codicon/codicon.ttf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "font/ttf")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Timing-Allow-Origin", "*")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("X-Cache", "HIT, HIT")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-JSD-Version", "0.28.1")
		w.Header().Set("X-JSD-Version-Type", "version")
		w.WriteHeader(200)
		_, _ = w.Write(codicon)
		return
	})
	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/vs/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-javascript")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Timing-Allow-Origin", "*")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("X-Cache", "HIT, HIT")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-JSD-Version", "0.28.1")
		w.Header().Set("X-JSD-Version-Type", "version")
		w.WriteHeader(200)
		_, _ = w.Write(indexJs)
		return
	})
	http.HandleFunc("/vs/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Timing-Allow-Origin", "*")
		w.Header().Set("Accept-Ranges", "bytes")
		w.WriteHeader(200)
		_, _ = w.Write(indexCSS)
		return
	})
	p := getAvailablePort()
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf("localhost:%d", p), nil)
	}()
	return p
}
func getAvailablePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	p := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return p
}

// WebSocket 处理函数
