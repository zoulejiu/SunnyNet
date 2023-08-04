package main

import "C"
import (
	"SunnyNet/project/SunnyNet"
	"SunnyNet/project/public"
	"SunnyNet/project/src/Certificate"
	"fmt"
	"time"
)

func Test() {

	s := SunnyNet.NewSunny()

	i := CreateCertificate()
	//LoadP12Certificate(i, C.CString("C:\\Users\\qinka\\Desktop\\certificate\\49D4E174.p12"), C.CString("xysj2017"))
	//LoadP12Certificate(i, C.CString("C:\\Users\\qinka\\Desktop\\certificate\\8F9AB6DF.p12"), C.CString("xysj2017"))
	LoadP12Certificate(i, C.CString("C:\\Users\\qinka\\Desktop\\certificate\\F8EF901E.p12"), C.CString("xysj2017"))
	c := Certificate.LoadCertificateContext(i)
	fmt.Println(c.GetServerName())
	AddHttpCertificate(C.CString("ccpay.cib.com.cn"), i, 1)

	//s.SetGlobalProxy("http://192.168.31.173:8888")
	//如果在Go中使用 设置Go的回调地址
	s.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)
	//s.SetIeProxy(false)
	//s.MustTcp(true)
	Port := 2022

	s = s.SetPort(Port).Start()

	// fmt.Println(s.StartProcess())

	// 请注意GoLang调试时候，请不要使用此(ProcessALLName)命令，因为不管开启或关闭，都会将当前所有TCP链接断开一次
	// 因为如果不断开的一次的话,已经建立的TCP链接无法抓包。
	// Go程序调试，是通过TCP连接的，若使用此命令将无法调试。
	// s.ProcessALLName(true)

	// s.ProcessAddName("sunny.exe")
	// s.ProcessAddName("sunny1.exe")
	// s.ProcessAddName("go_build_p.exe")
	err := s.Error
	if err != nil {
		panic(err)
	}
	fmt.Println("Run Port=", Port)
	time.Sleep(24 * time.Hour)
}
func HttpCallback(Conn *SunnyNet.HttpConn) {
	if Conn.Type == public.HttpSendRequest {
		//fmt.Println(Conn.Request.URL.String())
		//发起请求
	} else if Conn.Type == public.HttpResponseOK {
		//请求完成
	} else if Conn.Type == public.HttpRequestFail {
		//请求错误
	}
}
func WSCallback(Conn *SunnyNet.WsConn) {
	//fmt.Println(Conn.Url)
}
func TcpCallback(Conn *SunnyNet.TcpConn) {
	//fmt.Println(Conn.Pid, Conn.LocalAddr, Conn.RemoteAddr, Conn.Type, Conn.GetBodyLen())
}
func UdpCallback(Conn *SunnyNet.UDPConn) {
	if public.SunnyNetUDPTypeReceive == Conn.Type {
		fmt.Println("接收UDP", Conn.LocalAddress, Conn.RemoteAddress, len(Conn.Data))
	}
	if public.SunnyNetUDPTypeSend == Conn.Type {
		fmt.Println("发送UDP", Conn.LocalAddress, Conn.RemoteAddress, len(Conn.Data))
	}
	if public.SunnyNetUDPTypeClosed == Conn.Type {
		fmt.Println("关闭UDP", Conn.LocalAddress, Conn.RemoteAddress)
	}

}
