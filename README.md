# Sunny网络中间件

---

> Sunny网络中间件 和 Fiddler 类似。 是可跨平台的网络分析组件
 ```log 
 可用于HTTP/HTTPS/WS/WSS/TCP/UDP网络分析 为二次开发量身制作
 
 支持 获取/修改 HTTP/HTTPS/WS/WSS/TCP/TLS-TCP/UDP 发送及返回数据
 
 支持 对 HTTP/HTTPS/WS/WSS 指定连接使用指定代理
 
 支持 对 HTTP/HTTPS/WS/WSS/TCP/TLS-TCP 链接重定向
 
 支持 gzip, deflate, br, ZSTD 解码
 
 支持 WS/WSS/TCP/TLS-TCP/UDP 主动发送数据 
 
```

---
* # 由于代码主要是做DLL使用,部分功能未封装给Go使用，请自行探索！
* # 如需支持Win7系统
* # 请使用Go1.21以下版本编译,例如 go 1.20.4版本 
* # <a href="https://github.com/jmeubank/tdm-gcc/releases/download/v10.3.0-tdm64-2/tdm64-gcc-10.3.0-2.exe">编译请使用 TDM-GCC</a>
<center><h2><a style="color: red;">BUG 反馈</a></center></h2></center>
<center><h3>QQ群:751406884</center></h3></center>
<center><h3>二群：545120699</center></h3></center>
<center><h3>网址：<a href="https://esunny.vip/">https://esunny.vip/</a></center></h3></center>

---

### <center><h3>各语言,示例文件以及抓包工具 下载地址 </center>
<div style="text-align: center;"><h3>https://wwxa.lanzouu.com/b02p4aet8j</h3></div>
<div style="text-align: center;"><h3>密码:4h7r</h3></div>
<div style="text-align: center;"><h3></h3></div>


---
- > GoLang使用示例代码

```golang
package main

import (
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"time"
)

var Sunny = SunnyNet.NewSunny()

func main() {
	//绑定回调函数
	Sunny.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)

	//绑定端口号并启动
	Sunny.SetPort(2025).Start() 
	
	if Sunny.Error != nil {
		panic(Sunny.Error)
	}
	//避免程序退出
	time.Sleep(24 * time.Hour)
}

func HttpCallback(Conn SunnyNet.ConnHTTP) {

	if Conn.Type() == public.HttpSendRequest {
		//fmt.Println(Conn.URL())
		//发起请求

		//直接响应,不让其发送请求
		//Conn.StopRequest(200, "Hello Word")

	} else if Conn.Type() == public.HttpResponseOK {
		//请求完成
		//log.Println("Call", Conn.URL())
	} else if Conn.Type() == public.HttpRequestFail {
		//请求错误
		/*	fmt.Println(Conn.Request.URL.String(), Conn.GetError())
		 */
	}
}
func WSCallback(Conn SunnyNet.ConnWebSocket) {
	log.Println("WebSocket", Conn.URL())
}
func TcpCallback(Conn SunnyNet.ConnTCP) {

	if Conn.Type() == public.SunnyNetMsgTypeTCPAboutToConnect {
		//即将连接
		mode := string(Conn.Body())
		log.Println("PID", Conn.PID(), "TCP 即将连接到:", mode, Conn.LocalAddress(), "->", Conn.RemoteAddress())
		//修改目标连接地址
		//Conn.SetNewAddress("8.8.8.8:8080")
		return
	}

	if Conn.Type() == public.SunnyNetMsgTypeTCPConnectOK {
		log.Println("PID", Conn.PID(), "TCP 连接到:", Conn.LocalAddress(), "->", Conn.RemoteAddress(), "成功")
		return
	}

	if Conn.Type() == public.SunnyNetMsgTypeTCPClose {
		log.Println("PID", Conn.PID(), "TCP 断开连接:", Conn.LocalAddress(), "->", Conn.RemoteAddress())
		return
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientSend {
		log.Println("PID", Conn.PID(), "发送数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientReceive {
		log.Println("PID", Conn.PID(), "收到数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	}
}
func UdpCallback(Conn SunnyNet.ConnUDP) {

	if Conn.Type() == public.SunnyNetUDPTypeSend {
		//客户端向服务器端发送数据
		log.Println("PID", Conn.PID(), "发送UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改发送的数据
		//Conn.SetBody([]byte("Hello Word"))

		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeReceive {
		//服务器端向客户端发送数据
		log.Println("PID", Conn.PID(), "接收UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改响应的数据
		//Conn.SetBody([]byte("Hello Word"))
		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeClosed {

		log.Println("PID", Conn.PID(), "关闭UDP", Conn.LocalAddress(), Conn.RemoteAddress())
		return
	}

}
```