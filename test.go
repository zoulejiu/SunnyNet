package main

import "C"
import (
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/GoScriptCode"
	"github.com/qtgolang/SunnyNet/src/public"
	"log"
	"os"
)

func Test() {
	s := SunnyNet.NewSunny()
	//i := CreateCertificate()
	//ok := LoadP12Certificate(i, C.CString("C:\\Users\\qinka\\Desktop\\74fe394a37757545d8cfbd2ea264c7c3.p12"), C.CString("qyrhudhZ"))
	//ok := AddCertPoolPath(i, C.CString("C:\\Users\\qinka\\Desktop\\P12\\certificate.pem"))
	//fmt.Println("载入P12:", ok)
	//c := Certificate.LoadCertificateContext(i)
	//fmt.Println("证书名称：", c.GetCommonName())
	//AddHttpCertificate(C.CString("ws-gateway-odis.volkswagenag.com"), i, 1)
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
	s.SetGlobalProxy("socket://192.168.31.1:4321", 60000)
	s.SetMustTcpRegexp("zz.com", true)
	Port := 2024
	//s.SetMustTcpRegexp("*.baidu.com")
	s = s.SetPort(Port).Start()

	//s.SetIeProxy(true)
	s.SetHTTPRequestMaxUpdateLength(1000000)
	//fmt.Println(s.StartProcess())
	// 请注意GoLang调试时候，请不要使用此(ProcessALLName)命令，因为不管开启或关闭，都会将当前所有TCP链接断开一次
	// 因为如果不断开的一次的话,已经建立的TCP链接无法抓包。
	// Go程序调试，是通过TCP连接的，若使用此命令将无法调试。
	// s.ProcessALLName(true)

	//s.ProcessAddName("GoTest.exe")
	//s.ProcessAddName("msedge.exe")

	//s.ProcessAddName("pop_dd_workbench.exe")
	err := s.Error
	if err != nil {
		panic(err)
	}
	fmt.Println("Run Port=", Port)
	select {}
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
	return
	Conn.Context()
	//fmt.Println(Conn.Url)
}
func TcpCallback(Conn SunnyNet.ConnTCP) {

	return
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
	return
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
