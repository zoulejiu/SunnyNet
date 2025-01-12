package main

import "C"
import (
	"bytes"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/http"
	_ "github.com/qtgolang/SunnyNet/src/http/pprof"
	"github.com/qtgolang/SunnyNet/src/public"
	"github.com/qtgolang/SunnyNet/src/tlsClient/srt"
	tlsClient "github.com/qtgolang/SunnyNet/src/tlsClient/tlsClient"
	"github.com/qtgolang/SunnyNet/src/tlsClient/tlsClient/profiles"
	"io"
	"log"
	"os"
)

const information = `
------------------------------------------------------
       欢迎使用 SunnyNet 网络中间件 - V` + public.SunnyVersion + `   
                 本项目为开源项目  
            仅用于技术交流学习和研究的目的 
          请遵守法律法规,请勿用作任何非法用途 
               否则造成一切后果自负 
           若您下载并使用即视为您知晓并同意
------------------------------------------------------
        Sunny开源项目网站：https://esunny.vip
           Sunny QQ交流群(一群)：751406884
           Sunny QQ交流群(二群)：545120699
           Sunny QQ交流群(三群)：170902713
       QQ频道：https://pd.qq.com/g/SunnyNetV5
------------------------------------------------------

`

func init() {
	fmt.Println(information)
	go http.ListenAndServe("0.0.0.0:6001", nil)
}

func main() {
	Test()
	//阻止程序退出
	select {}
}

func main3() {
	defer os.Exit(0)
	// Create a SpoofedRoundTripper that implements the http.RoundTripper interface
	tr, err := srt.NewSpoofedRoundTripper(
		// Reference for more: https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/client-options
		tlsClient.WithRandomTLSExtensionOrder(), // needed for Chrome 107+
		tlsClient.WithClientProfile(profiles.Chrome_124),
	)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: tr,
	}

	Body := io.NopCloser(bytes.NewBuffer([]byte("")))
	defer func() { _ = Body.Close() }()

	req, err := http.NewRequest("GET", ("https://tls.browserleaks.com/json"), Body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Host", "www.qq.com")
	req.Header.Set("Token", "123456789")
	req.Header.Set("aBC", "6666xxxxxx66")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("响应状态码", res.StatusCode)
	for k, v := range res.Header {
		fmt.Println(k, v)
	}
	fmt.Println(string(body))
}
func main2() {
	h1s := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(h1s.ListenAndServeTLS("D:\\Go\\PATH\\pkg\\mod\\github.com\\globalsign\\mgo@v0.0.0-20181015135952-eeefdecb41b8\\harness\\certs\\server.crt", "D:\\Go\\PATH\\pkg\\mod\\github.com\\globalsign\\mgo@v0.0.0-20181015135952-eeefdecb41b8\\harness\\certs\\server.key"))
}
