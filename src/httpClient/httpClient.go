package httpClient

import (
	"context"
	"errors"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/SunnyProxy"
	tls "github.com/qtgolang/SunnyNet/src/crypto/tls"
	"github.com/qtgolang/SunnyNet/src/dns"
	"github.com/qtgolang/SunnyNet/src/http"
	"github.com/qtgolang/SunnyNet/src/public"
	"net"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type Client struct {
	proxy    *SunnyProxy.Proxy
	redirect bool
	//RandomTLSFingerprint 随机TLS指纹
	RandomTLSFingerprint func() []uint16
	tlsConfig            *tls.Config
	outTime              time.Duration
}

// SetProxy 设置代理
func (w *Client) SetProxy(RequestProxy *SunnyProxy.Proxy) {
	w.proxy = RequestProxy
}

// Redirect 设置是否重定向
func (w *Client) Redirect(allow bool) {
	w.redirect = allow
}

// SetOutTime 设置超时时间 单位毫秒 小于 1 则默认30秒
func (w *Client) SetOutTime(outTime int) {
	if outTime < 1 {
		w.outTime = 0
		return
	}
	w.outTime = time.Duration(outTime) * time.Millisecond
}

// Do 发送请求
func (w *Client) Do(req *http.Request) (Response *http.Response, Conn net.Conn, err error, Close func()) {
	if w.tlsConfig == nil {
		w.tlsConfig = &tls.Config{}
	}
	Response, Conn, err, Close = Do(req, w.proxy, w.redirect, w.tlsConfig, w.outTime, w.RandomTLSFingerprint)
	return Response, Conn, err, Close
}
func Do(req *http.Request, RequestProxy *SunnyProxy.Proxy, CheckRedirect bool, config *tls.Config, outTime time.Duration, GetTLSValues func() []uint16) (Response *http.Response, Conn net.Conn, err error, Close func()) {
	if req.ProtoMajor == 2 {
		Method := req.Method
		switch Method {
		case public.HttpMethodHEAD:
			fallthrough
		case public.HttpMethodGET:
			fallthrough
		case public.HttpMethodTRACE:
			fallthrough
		case public.HttpMethodOPTIONS:
			if req.Body != nil {
				_ = req.Body.Close()
				req.Body = nil
			}
		default:
			break
		}

	}
	cfg := config.Clone()
	if req.URL != nil && req.URL.Scheme != "http" {
		if cfg == nil {
			cfg = &tls.Config{}
		}
		cfg.InsecureSkipVerify = true
	}
	handshakeCount := 0
	for {
		if cfg != nil && GetTLSValues != nil {
			tv := GetTLSValues()
			if len(tv) > 0 {
				cfg.CipherSuites = tv
			}
		}
		Response, Conn, err, Close = do(req, RequestProxy, CheckRedirect, cfg, outTime)
		if cfg == nil {
			return
		}
		if err != nil {
			if Conn != nil {
				_ = Conn.Close()
			}
			ers := err.Error()
			if strings.Contains(ers, "handshake") || strings.Contains(ers, "connection") || strings.Contains(ers, "EOF") {
				handshakeCount++
				if handshakeCount > 10 {
					Close = nil
					return
				}
				continue
			}
		}
		return
	}
}
func do(req *http.Request, RequestProxy *SunnyProxy.Proxy, CheckRedirect bool, config *tls.Config, outTime time.Duration) (*http.Response, net.Conn, error, func()) {
	client := httpClientGet(req, RequestProxy, config, outTime)
	if CheckRedirect {
		client.Client.CheckRedirect = public.HTTPAllowRedirect
	} else {
		client.Client.CheckRedirect = public.HTTPBanRedirect
	}
	reqs, err := client.Client.Do(req)
	var rConn net.Conn
	if reqs != nil {
		if reqs.Request != nil {
			rConn, _ = reqs.Request.Context().Value("rConn").(net.Conn)
		}
		if reqs.Header != nil {
			reqs.Header.Del("Transfer-Encoding")
		}
	}
	return reqs, rConn, err, func() { httpClientPop(client) }
}

var httpLock sync.Mutex
var httpClientMap map[string]clientList

type clientList map[uintptr]*clientPart

func httpClientGet(req *http.Request, Proxy *SunnyProxy.Proxy, cfg *tls.Config, timeout time.Duration) *clientPart {
	httpLock.Lock()
	defer httpLock.Unlock()
	s := ""
	if req != nil && req.URL != nil {
		s = req.URL.Host + "|" + req.Proto + "|" + req.URL.Scheme
	}
	s += "|" + Proxy.String() + "|"
	if cfg != nil {
		s += strings.Join(cfg.NextProtos, "-")
	}
	if clients, ok := httpClientMap[s]; ok {
		if len(clients) > 0 {
			for key, client := range clients {
				delete(clients, key)
				client.RequestProxy = Proxy
				return client
			}
		}
	}
	if cfg != nil {
		if len(cfg.NextProtos) > 0 {
			cfg.GetConfigForServer = func(info *tls.ServerHelloMsg) error {
				for _, proto := range cfg.NextProtos {
					if proto == http.H2Proto && info.SupportedVersion == 772 {
						return nil // 如果支持，则返回 nil
					}
					if proto == http.H11Proto && (info.SupportedVersion == 0 || info.Vers == 771) {
						return nil // 如果支持，则返回 nil
					}
				}
				ver := info.SupportedVersion
				if ver == 0 {
					ver = info.Vers
				}
				Proto, _ := http.ProtoVersions[info.Vers]
				if Proto == "" {
					return fmt.Errorf("服务器不支持您所选HTTP协议版本")
				}
				return fmt.Errorf("服务器不支持您所选HTTP协议版本: 需要协议[%s],请检查您的配置", strings.ToUpper(Proto))
			}
		}
	}
	Tr := &http.Transport{TLSClientConfig: cfg}
	Tr.ResponseHeaderTimeout = timeout // 读取响应头超时
	Tr.IdleConnTimeout = timeout       // 空闲连接超时
	Tr.TLSHandshakeTimeout = timeout   // TLS 握手超时
	/*
		Dial := func(network, addr string) (net.Conn, error) {
			if RequestProxy != nil {
				return RequestProxy.Dial(network, addr)
			}
			var d net.Dialer
			return d.Dial(network, addr)
		}
		Tr.Dial = Dial
		Tr.DialTLS = func(network, addr string) (net.Conn, error) {
			return nil, nil
		}
	*/
	if cfg != nil {
		if len(cfg.NextProtos) < 1 {
			configureHTTP2Transport(Tr, cfg)
		} else {
			for _, proto := range cfg.NextProtos {
				if proto == http.H2Proto {
					configureHTTP2Transport(Tr, cfg)
					break
				}
			}
		}
	}
	var ips []net.IP
	var isLookupIP bool
	var ProxyHost string
	var dial func(network string, addr string) (net.Conn, error)
	if Proxy != nil {
		ProxyHost = Proxy.Host
		dial = Proxy.Dial
	}
	cc := http.Client{Transport: Tr, Timeout: timeout}
	res := &clientPart{Client: cc, s: s, RequestProxy: Proxy}
	Tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {

		address, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		i := net.ParseIP(address)
		if i != nil {
			if len(i) == net.IPv4len {
				return res.RequestProxy.Dial(network, i.String()+":"+port)
			}
			return res.RequestProxy.Dial(network, fmt.Sprintf("[%s]:%s", address, port))
		}
		if strings.ToLower(address) == "localhost" {
			return res.RequestProxy.Dial(network, "127.0.0.1:"+port)
		}
		var retries bool
		for {
			if !isLookupIP {
				isLookupIP = true
				first := dns.GetFirstIP(address, ProxyHost)
				ips, _ = dns.LookupIP(address, ProxyHost, dial)
				if first != nil {
					if first.To4() != nil {
						return res.RequestProxy.Dial(network, fmt.Sprintf("%s:%s", first.String(), port))
					} else {
						return res.RequestProxy.Dial(network, fmt.Sprintf("[%s]:%s", first.String(), port))
					}
				}
			}
			if len(ips) < 1 {
				dns.SetFirstIP(address, ProxyHost, nil)
				if retries {
					return nil, noIP
				}
				isLookupIP = false
				retries = true
				continue
			}
			ip := extractAndRemoveIP(&ips)
			if ip != nil {
				if ip.To4() != nil {
					conn, er := res.RequestProxy.Dial(network, fmt.Sprintf("%s:%s", ip.String(), port))
					if conn != nil {
						dns.SetFirstIP(address, ProxyHost, ip)
					}
					return conn, er
				}
				conn, er := res.RequestProxy.Dial(network, fmt.Sprintf("[%s]:%s", ip.String(), port))
				if conn != nil {
					dns.SetFirstIP(address, ProxyHost, ip)
				}
				return conn, er
			}
		}
	}
	return res
}

// 优先使用IPV4
func extractAndRemoveIP(ips *[]net.IP) net.IP {
	for i, ip := range *ips {
		if ip.To4() != nil { // 检查是否为 IPv4
			// 找到 IPv4，删除并返回
			*ips = append((*ips)[:i], (*ips)[i+1:]...) // 删除
			return ip
		}
	}
	// 如果没有找到 IPv4，查找 IPv6
	for i, ip := range *ips {
		if ip.To16() != nil { // IPv6 的情况
			*ips = append((*ips)[:i], (*ips)[i+1:]...) // 删除
			return ip
		}
	}
	return nil
}

var noIP = errors.New("DNS解析失败,无可用IP地址")

func configureHTTP2Transport(Tr *http.Transport, cfg *tls.Config) {
	// 检查是否配置了 HTTP/2.0 协议
	protoFound := false
	for _, proto := range cfg.NextProtos {
		if proto == http.H2Proto {
			protoFound = true
			break
		}
	}

	// 如果找到了 HTTP/2.0 协议，则配置 HTTP/2.0 传输
	if protoFound || len(cfg.NextProtos) == 0 {
		http.HTTP2configureTransport(Tr)
	}
}

type clientPart struct {
	Client       http.Client
	time         time.Time
	s            string
	RequestProxy *SunnyProxy.Proxy
}

func httpClientPop(client *clientPart) {
	if client == nil || client.s == "" {
		return
	}
	httpLock.Lock()
	defer httpLock.Unlock()
	client.time = time.Now()
	clients := httpClientMap[client.s]
	if clients == nil {
		httpClientMap[client.s] = make(clientList)
		clients = httpClientMap[client.s]
	}
	clients[uintptr(unsafe.Pointer(client))] = client
}
func httpClientClear() {
	httpLock.Lock()
	defer httpLock.Unlock()
	t := time.Now()
	o := 5 * time.Second
	for k, clients := range httpClientMap {
		for key, client := range clients {
			if t.Sub(client.time) > o {
				delete(clients, key)
			}
		}
		if len(clients) == 0 {
			delete(httpClientMap, k)
		}
	}
}
func init() {
	httpClientMap = make(map[string]clientList)
	go func() {
		for {
			time.Sleep(time.Second * 3)
			httpClientClear()
		}
	}()
}
