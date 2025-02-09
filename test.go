package main

import "C"
import (
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/GoScriptCode"
	"github.com/qtgolang/SunnyNet/src/public"
	"log"
	"os"
	"time"
)

func Test() {
	s := SunnyNet.NewSunny()
	cert := SunnyNet.NewCertManager()
	ok := cert.LoadP12Certificate("C:\\Users\\Qin\\Desktop\\Cert\\ca6afc5aa40fcbd3.p12", "GXjc75IRAO0T")
	fmt.Println("载入P12:", ok)
	fmt.Println("证书名称：", cert.GetCommonName())
	s.AddHttpCertificate("api.vlightv.com", cert, SunnyNet.HTTPCertRules_Request)

	//如果在Go中使用 设置Go的回调地址
	//s.SetGlobalProxy("socket://192.168.31.1:4321", 30000)
	//s.SetScriptCall(func(info ...any) {
	//	fmt.Println("x脚本日志", fmt.Sprintf("%v", info))
	//}, func(code []byte) {})
	s.SetScriptCode(string(GoScriptCode.DefaultCode))
	s.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)
	//s.SetMustTcpRegexp("*.baidu.com")
	s.CompileProxyRegexp("127.0.0.1;[::1];192.168.*")

	//s.MustTcp(true)
	//s.DisableTCP(true)
	//s.SetGlobalProxy("socket://192.168.31.1:4321", 60000)
	s.SetMustTcpRegexp("tpstelemetry.tencent.com", true)
	Port := 2025
	//s.SetMustTcpRegexp("*.baidu.com")
	s.SetPort(Port).Start()
	//s.SetIEProxy()
	s.SetHTTPRequestMaxUpdateLength(100000000)
	fmt.Println(s.OpenDrive(false))
	//s.ProcessALLName(true, false)

	//s.ProcessAddName("WeChat.exe")
	err := s.Error
	if err != nil {
		panic(err)
	}
	fmt.Println("Run Port=", Port)
	//阻止程序退出
	select {}
}
func HttpCallback(Conn SunnyNet.ConnHTTP) {
	if Conn.Type() == public.HttpSendRequest {
		fmt.Println("发起请求", Conn.URL())
		//发起请求
		Conn.SetResponseBody([]byte("123456"))
		//直接响应,不让其发送请求
		//Conn.StopRequest(200, "Hello Word")

	} else if Conn.Type() == public.HttpResponseOK {
		//请求完成
		bs := Conn.GetResponseBody()
		log.Println("请求完成", Conn.URL(), len(bs), Conn.GetResponseHeader())
	} else if Conn.Type() == public.HttpRequestFail {
		//请求错误
		fmt.Println(time.Now(), Conn.URL(), Conn.Error())
	}
}
func WSCallback(Conn SunnyNet.ConnWebSocket) {
	fmt.Println("WebSocket", Conn.URL())
}
func TcpCallback(Conn SunnyNet.ConnTCP) {
	if Conn.Type() == public.SunnyNetMsgTypeTCPAboutToConnect {
		//即将连接
		mode := string(Conn.Body())
		info.Println("PID", Conn.PID(), "TCP 即将连接到:", mode, Conn.LocalAddress(), "->", Conn.RemoteAddress())
		//修改目标连接地址
		//Conn.SetNewAddress("8.8.8.8:8080")
		return
	}

	if Conn.Type() == public.SunnyNetMsgTypeTCPConnectOK {
		info.Println("PID", Conn.PID(), "TCP 连接到:", Conn.LocalAddress(), "->", Conn.RemoteAddress(), "成功")
		return
	}

	if Conn.Type() == public.SunnyNetMsgTypeTCPClose {
		info.Println("PID", Conn.PID(), "TCP 断开连接:", Conn.LocalAddress(), "->", Conn.RemoteAddress())
		return
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientSend {
		info.Println("PID", Conn.PID(), "发送数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	}
	if Conn.Type() == public.SunnyNetMsgTypeTCPClientReceive {
		info.Println("PID", Conn.PID(), "收到数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	}
}
func UdpCallback(Conn SunnyNet.ConnUDP) {

	if Conn.Type() == public.SunnyNetUDPTypeSend {
		//客户端向服务器端发送数据
		info.Println("PID", Conn.PID(), "发送UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改发送的数据
		//Conn.SetBody([]byte("Hello Word"))

		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeReceive {
		//服务器端向客户端发送数据
		info.Println("PID", Conn.PID(), "接收UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改响应的数据
		//Conn.SetBody([]byte("Hello Word"))
		return
	}
	if Conn.Type() == public.SunnyNetUDPTypeClosed {

		info.Println("PID", Conn.PID(), "关闭UDP", Conn.LocalAddress(), Conn.RemoteAddress())
		return
	}

}

var info = log.New(os.Stdout, "", log.LstdFlags)
