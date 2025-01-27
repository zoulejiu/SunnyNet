package main

import "C"
import (
	"fmt"
	"github.com/qtgolang/SunnyNet/src/http"
	_ "github.com/qtgolang/SunnyNet/src/http/pprof"
	"github.com/qtgolang/SunnyNet/src/public"
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
}
