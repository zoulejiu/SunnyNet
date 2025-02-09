package main

import "C"
import (
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/GoScriptCode"
	"github.com/qtgolang/SunnyNet/src/public"
	"log"
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
	switch Conn.Type() {
	case public.HttpSendRequest: //发起请求
		fmt.Println("发起请求", Conn.URL())
		Conn.SetResponseBody([]byte("123456"))
		//直接响应,不让其发送请求
		//Conn.StopRequest(200, "Hello Word")
		return
	case public.HttpResponseOK: //请求完成
		bs := Conn.GetResponseBody()
		log.Println("请求完成", Conn.URL(), len(bs), Conn.GetResponseHeader())
		return
	case public.HttpRequestFail: //请求错误
		fmt.Println(time.Now(), Conn.URL(), Conn.Error())
		return
	}
}
func WSCallback(Conn SunnyNet.ConnWebSocket) {
	fmt.Println("WebSocket", Conn.URL())
}
func TcpCallback(Conn SunnyNet.ConnTCP) {
	switch Conn.Type() {
	case public.SunnyNetMsgTypeTCPAboutToConnect: //即将连接
		mode := string(Conn.Body())
		log.Println("PID", Conn.PID(), "TCP 即将连接到:", mode, Conn.LocalAddress(), "->", Conn.RemoteAddress())
		//修改目标连接地址
		//Conn.SetNewAddress("8.8.8.8:8080")
		return
	case public.SunnyNetMsgTypeTCPConnectOK: //连接成功
		log.Println("PID", Conn.PID(), "TCP 连接到:", Conn.LocalAddress(), "->", Conn.RemoteAddress(), "成功")
		return
	case public.SunnyNetMsgTypeTCPClose: //连接关闭
		log.Println("PID", Conn.PID(), "TCP 断开连接:", Conn.LocalAddress(), "->", Conn.RemoteAddress())
		return
	case public.SunnyNetMsgTypeTCPClientSend: //客户端发送数据
		log.Println("PID", Conn.PID(), "TCP 发送数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	case public.SunnyNetMsgTypeTCPClientReceive: //客户端收到数据

		log.Println("PID", Conn.PID(), "收到数据", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.Type(), Conn.BodyLen(), Conn.Body())
		return
	default:
		return
	}
}
func UdpCallback(Conn SunnyNet.ConnUDP) {
	switch Conn.Type() {
	case public.SunnyNetUDPTypeSend: //客户端向服务器端发送数据

		log.Println("PID", Conn.PID(), "发送UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改发送的数据
		//Conn.SetBody([]byte("Hello Word"))

		return
	case public.SunnyNetUDPTypeReceive: //服务器端向客户端发送数据
		log.Println("PID", Conn.PID(), "接收UDP", Conn.LocalAddress(), Conn.RemoteAddress(), Conn.BodyLen())
		//修改响应的数据
		//Conn.SetBody([]byte("Hello Word"))
		return
	case public.SunnyNetUDPTypeClosed: //关闭会话
		log.Println("PID", Conn.PID(), "关闭UDP", Conn.LocalAddress(), Conn.RemoteAddress())
		return
	}

}
