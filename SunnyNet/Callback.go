package SunnyNet

import "C"
import (
	"github.com/qtgolang/SunnyNet/src/Call"
	"github.com/qtgolang/SunnyNet/src/public"
	"strconv"
)

const debug = "Debug"

func GetSceneProxyRequest(MessageId int) (*proxyRequest, bool) {
	messageIdLock.Lock()
	defer messageIdLock.Unlock()
	k := httpStorage[MessageId]
	if k == nil {
		return nil, false
	}
	return k, true
}
func GetSceneWebSocketMsg(MessageId int) (*public.WebsocketMsg, bool) {
	messageIdLock.Lock()
	defer messageIdLock.Unlock()
	k := wsStorage[MessageId]
	if k == nil {
		return nil, false
	}
	return k, true
}

// CallbackTCPRequest TCP请求处理回调
func (s *proxyRequest) CallbackTCPRequest(callType int, msg *public.TcpMsg, RemoteAddr string) {
	if s.Global.disableTCP {
		//由于用户可能在软件中途禁用TCP,所有这里允许触发关闭的回调
		if callType != public.SunnyNetMsgTypeTCPClose {
			//这里如果禁用了TCP,那么这里就不允许触发回调了，并且手动关闭连接
			TcpSceneLock.Lock()
			w := TcpStorage[s.Theology]
			TcpSceneLock.Unlock()
			if w == nil {
				return
			}
			w.L.Lock()
			_ = w.ConnSend.Close()
			_ = w.ConnServer.Close()
			w.L.Unlock()
			return
		}
	}
	LocalAddr := s.Conn.RemoteAddr().String()
	hostname := RemoteAddr
	pid, _ := strconv.Atoi(s.Pid)
	MessageId := NewMessageId()

	messageIdLock.Lock()
	httpStorage[MessageId] = s
	messageIdLock.Unlock()

	defer func() {
		messageIdLock.Lock()
		httpStorage[MessageId] = nil
		delete(httpStorage, MessageId)
		messageIdLock.Unlock()
	}()

	Ams := &tcpConn{c: msg, messageId: MessageId, _type: callType, theology: s.Theology, localAddr: LocalAddr, remoteAddr: hostname, pid: pid, sunnyContext: s.Global.SunnyContext}
	s.Global.scriptTCPCall(Ams)
	if s.TcpCall < 10 {
		if s.TcpGoCall != nil {
			s.TcpGoCall(Ams)
		}
		return
	}
	if callType == public.SunnyNetMsgTypeTCPConnectOK {
		Call.Call(s.TcpCall, s.Global.SunnyContext, LocalAddr, hostname, int(callType), MessageId, msg.Data.Bytes(), msg.Data.Len(), s.Theology, pid)
		return
	}
	if callType == public.SunnyNetMsgTypeTCPClose {
		Call.Call(s.TcpCall, s.Global.SunnyContext, LocalAddr, hostname, int(callType), MessageId, []byte{}, 0, s.Theology, pid)
		return
	}
	if callType == public.SunnyNetMsgTypeTCPClientSend || callType == public.SunnyNetMsgTypeTCPAboutToConnect {
		s.TCP.Send = msg
	} else {
		s.TCP.Receive = msg
	}
	Call.Call(s.TcpCall, s.Global.SunnyContext, LocalAddr, hostname, int(callType), MessageId, msg.Data.Bytes(), msg.Data.Len(), s.Theology, pid)
}

// CallbackBeforeRequest HTTP发起请求处理回调
func (s *proxyRequest) CallbackBeforeRequest() {

	if s.Response.Response != nil {
		if s.Response.Body != nil {
			_ = s.Response.Body.Close()
		}
	}

	pid, _ := strconv.Atoi(s.Pid)
	s.Response.Response = nil
	defer func() {
		if s.Response.Response != nil {
			if s.Response.Response.StatusCode == 0 && len(s.Response.Header) == 0 {
				if s.Response.Body != nil {
					_ = s.Response.Body.Close()
				}
				s.Response.Response = nil
			}
		}
	}()
	MessageId := NewMessageId()
	messageIdLock.Lock()
	httpStorage[MessageId] = s
	messageIdLock.Unlock()
	defer func() {
		messageIdLock.Lock()
		httpStorage[MessageId] = nil
		delete(httpStorage, MessageId)
		messageIdLock.Unlock()
	}()

	m := &httpConn{
		_Theology:   s.Theology,
		_getRawBody: s.RawRequestDataToFile,
		_MessageId:  MessageId,
		_PID:        pid,
		_Context:    s.Global.SunnyContext,
		_Type:       public.HttpSendRequest,
		_request:    s.Request,
		_response:   s.Response.Response,
		_err:        "",
		_proxy:      s.Proxy,
		_ClientIP:   s.Conn.RemoteAddr().String(),
		_Display:    true,
		_Break:      false,
		_tls:        s.TlsConfig,
		_serverIP:   s.Response.ServerIP,
	}
	s.Global.scriptHTTPCall(m)
	s._Display = m._Display
	if s._Display == false {
		return
	}
	err := ""
	if m._Break {
		err = debug
	}

	if s.HttpCall < 10 {
		if s.HttpGoCall != nil {
			s.HttpGoCall(m)
			s.Response.Response = m._response
		}
		return
	}
	Method := s.Request.Method
	Url := s.Request.URL.String()
	Call.Call(s.HttpCall, s.Global.SunnyContext, s.Theology, MessageId, int(public.HttpSendRequest), Method, Url, err, pid)
}

// CallbackBeforeResponse HTTP请求完成处理回调
func (s *proxyRequest) CallbackBeforeResponse() {

	pid, _ := strconv.Atoi(s.Pid)

	MessageId := NewMessageId()

	messageIdLock.Lock()
	httpStorage[MessageId] = s
	messageIdLock.Unlock()
	defer func() {
		messageIdLock.Lock()
		httpStorage[MessageId] = nil
		delete(httpStorage, MessageId)
		messageIdLock.Unlock()
	}()

	m := &httpConn{
		_Theology:   s.Theology,
		_getRawBody: s.RawRequestDataToFile,
		_MessageId:  MessageId,
		_PID:        pid,
		_Context:    s.Global.SunnyContext,
		_Type:       public.HttpResponseOK,
		_request:    s.Request,
		_response:   s.Response.Response,
		_err:        "",
		_ClientIP:   s.Conn.RemoteAddr().String(),
		_Display:    true,
		_Break:      false,
		_tls:        s.TlsConfig,
		_serverIP:   s.Response.ServerIP,
	}
	s.Global.scriptHTTPCall(m)
	if s._Display == false {
		return
	}
	err := ""
	if m._Break {
		err = debug
	}
	if s.HttpCall < 10 {
		if s.HttpGoCall != nil {
			s.HttpGoCall(m)
		}
		return
	}
	Method := s.Request.Method
	Url := s.Request.URL.String()
	Call.Call(s.HttpCall, s.Global.SunnyContext, s.Theology, MessageId, int(public.HttpResponseOK), Method, Url, err, pid)
}

// CallbackError HTTP请求失败处理回调
func (s *proxyRequest) CallbackError(err string) {

	pid, _ := strconv.Atoi(s.Pid)
	MessageId := NewMessageId()
	messageIdLock.Lock()
	httpStorage[MessageId] = s
	messageIdLock.Unlock()
	defer func() {
		messageIdLock.Lock()
		httpStorage[MessageId] = nil
		delete(httpStorage, MessageId)
		messageIdLock.Unlock()
	}()
	m := &httpConn{
		_Theology:   s.Theology,
		_getRawBody: s.RawRequestDataToFile,
		_MessageId:  NewMessageId(),
		_PID:        pid,
		_Context:    s.Global.SunnyContext,
		_Type:       public.HttpRequestFail,
		_request:    s.Request,
		_response:   nil,
		_err:        err,
		_ClientIP:   s.Conn.RemoteAddr().String(),
		_tls:        nil,
		_serverIP:   s.Response.ServerIP,
	}
	s.Global.scriptHTTPCall(m)
	if s._Display == false {
		return
	}
	if s.HttpCall < 10 {
		if s.HttpGoCall != nil {
			s.HttpGoCall(m)
		}
		return
	}
	//请求失败
	Method := s.Request.Method
	Url := "Unknown URL"
	if s.Request.URL != nil {
		Url = s.Request.URL.String()
	}
	Call.Call(s.HttpCall, s.Global.SunnyContext, s.Theology, MessageId, int(public.HttpRequestFail), Method, Url, err, pid)

}

// CallbackWssRequest HTTP->Websocket请求处理回调
func (s *proxyRequest) CallbackWssRequest(State int, Method, Url string, msg *public.WebsocketMsg, MessageId int) {

	if s._Display == false {
		return
	}
	pid, _ := strconv.Atoi(s.Pid)
	m := &wsConn{Pid: pid, _Type: State, SunnyContext: s.Global.SunnyContext, Url: Url, c: msg, _MessageId: MessageId, _Theology: s.Theology, Request: s.Request, _ClientIP: s.Conn.RemoteAddr().String()}
	s.Global.scriptWebsocketCall(m)
	if s.wsCall < 10 {
		if s.wsGoCall != nil {
			s.wsGoCall(m)
		}
		return
	}
	Call.Call(s.wsCall, s.Global.SunnyContext, s.Theology, MessageId, State, Method, Url, pid, msg.Mt)
}
