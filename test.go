package main

import "C"
import (
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/src/GoScriptCode"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"github.com/qtgolang/SunnyNet/src/public"
	"log"
)

func Test() {
	s := SunnyNet.NewSunny()
	//cert := SunnyNet.NewCertManager()
	//ok := cert.LoadP12Certificate("C:\\Users\\Qin\\Desktop\\Cert\\ca6afc5aa40fcbd3.p12", "GXjc75IRAO0T")
	//fmt.Println("载入P12:", ok)
	//fmt.Println("证书名称：", cert.GetCommonName())
	//s.AddHttpCertificate("api.vlightv.com", cert, SunnyNet.HTTPCertRules_Request)

	//如果在Go中使用 设置Go的回调地址
	//s.SetGlobalProxy("socket://137.11.0.11:205", 30000)
	s.SetScriptCall(func(Context int, info ...any) {
		fmt.Println("x脚本日志", fmt.Sprintf("%v", info))
	}, func(Context int, code []byte) {})
	s.SetScriptCode(string(GoScriptCode.DefaultCode))
	//s.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)
	s.SetGoCallback(HttpCallback1, nil, nil, nil)
	//s.SetMustTcpRegexp("*.baidu.com")
	s.CompileProxyRegexp("127.0.0.1;[::1];192.168.*")
	//https://api.pmangplus.com/apps/631/maintenance_banner?device_cd=ANDROID&local_cd=KOR
	//socket5://qj07:123@61.75.40.187:9021
	//s.MustTcp(true)
	//s.DisableTCP(true)
	//s.SetGlobalProxy("socket://192.168.31.1:4321", 60000)
	s.SetMustTcpRegexp("api.pmangplus.com;ip138.com;*.ip138.com;*.ipip.net", false)
	Port := 2225
	//s.SetMustTcpRegexp("*.baidu.com")
	s.SetPort(Port).Start()
	//s.SetIEProxy()
	//s.SetHTTPRequestMaxUpdateLength(10086)
	fmt.Println("加载驱动", s.OpenDrive(false))

	s.ProcessAddName("dnplayer.exe")
	s.ProcessAddName("LdBoxHeadless.exe")
	s.ProcessAddName("LdVBoxHeadless.exe")
	s.ProcessAddName("Ld9BoxHeadless.exe")
	s.ProcessAddName("VBoxNetDHCP.exe")
	s.ProcessAddName("VBoxNetNAT.exe")
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
func HttpCallback1(Conn SunnyNet.ConnHTTP) {
	switch Conn.Type() {
	case public.HttpSendRequest: //发起请求
		Conn.SetAgent("socket5://qj07:123@61.75.40.187:9021", 3000*10)
		fmt.Println("发起请求", Conn.URL())
		return
		Conn.SetResponseBody([]byte("123456"))
		//直接响应,不让其发送请求
		//Conn.StopRequest(200, "Hello Word")
		return
	case public.HttpResponseOK: //请求完成
		//bs := Conn.GetResponseBody()
		//log.Println("请求完成", Conn.URL(), len(bs), Conn.GetResponseHeader())
		return
	case public.HttpRequestFail: //请求错误
		//fmt.Println(time.Now(), Conn.URL(), Conn.Error())
		return
	}
}
func WSCallbackX(Conn SunnyNet.ConnWebSocket) {
	switch Conn.Type() {
	case public.WebsocketConnectionOK: //连接成功
		log.Println("PID", Conn.PID(), "Websocket 连接成功:", Conn.URL())
		return
	case public.WebsocketUserSend: //发送数据
		if Conn.MessageType() < 5 {
			log.Println("PID", Conn.PID(), "Websocket 发送数据:", Conn.MessageType(), "->", hex.EncodeToString(Conn.Body()))
		}
		return
	case public.WebsocketServerSend: //收到数据
		if Conn.MessageType() < 5 {
			log.Println("PID", Conn.PID(), "Websocket 收到数据:", Conn.MessageType(), "->", hex.EncodeToString(Conn.Body()))
		}
		return
	case public.WebsocketDisconnect: //连接关闭
		log.Println("PID", Conn.PID(), "Websocket 连接关闭", Conn.URL())
		return
	default:
		return
	}
}
func TcpCallbackX(Conn SunnyNet.ConnTCP) {

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
func UdpCallbackX(Conn SunnyNet.ConnUDP) {

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
