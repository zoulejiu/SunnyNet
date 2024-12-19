package dns

import (
	"context"
	"github.com/qtgolang/SunnyNet/src/crypto/tls"
	"net"
	"strings"
	"sync"
	"time"
)

var dnsConfig = &tls.Config{
	ClientSessionCache: tls.NewLRUClientSessionCache(32),
	InsecureSkipVerify: true,
}
var dnsList = make(map[string]*rsIps)
var dnsLock sync.Mutex
var dnsTools = make(map[string]*tools)
var dnsServer = "223.5.5.5:853" //阿里云公共DNS解析服务器

func init() {
	go clean()
}
func newResolver(proxy string, Dial func(network, address string) (net.Conn, error)) *net.Resolver {
	var dialer net.Dialer
	_default_ := &net.Resolver{
		PreferGo: true,
		Dial: func(context context.Context, network_, address string) (net.Conn, error) {
			dnsLock.Lock()
			_dnsServer := dnsServer + ""
			dnsLock.Unlock()
			var conn net.Conn
			var err error
			if _dnsServer == "" {
				//本地解析
				conn, err = dialer.DialContext(context, network_, address)
			} else if proxy == "" {
				//本地使用自定义DNS服务器解析
				conn, err = dialer.DialContext(context, "tcp", _dnsServer)
			} else {
				//使用代理并使用DNS服务器解析
				conn, err = Dial("tcp", _dnsServer)
			}
			if err != nil {
				return nil, err
			}
			_ = conn.(*net.TCPConn).SetKeepAlive(true)
			_ = conn.(*net.TCPConn).SetKeepAlivePeriod(10 * time.Minute)
			if !strings.HasSuffix(_dnsServer, ":853") {
				return conn, nil
			}
			return tls.Client(conn, dnsConfig), nil
		},
	}
	return _default_
}
func clean() {
	for {
		time.Sleep(time.Minute)
		dnsLock.Lock()
		for key, value := range dnsTools {
			if time.Now().Sub(value.time) > time.Minute*10 {
				delete(dnsTools, key)
			}
		}
		for key, value := range dnsList {
			if time.Now().Sub(value.time) > time.Minute*10 {
				delete(dnsList, key)
			}
		}
		dnsLock.Unlock()
	}
}

type tools struct {
	rs   *net.Resolver
	time time.Time
}
type rsIps struct {
	ips   []net.IP
	first net.IP
	time  time.Time
}

func SetDnsServer(server string) {
	dnsLock.Lock()
	dnsServer = server
	dnsList = make(map[string]*rsIps)
	dnsTools = make(map[string]*tools)
	dnsLock.Unlock()
}
func SetFirstIP(host string, proxyHost string, ip net.IP) {
	key := ""
	if proxyHost == "" {
		key = "_default_" + host
	} else {
		key = proxyHost + "|" + host
	}
	dnsLock.Lock()
	ips := dnsList[key]
	if ips != nil {
		ips.first = ip
		ips.time = time.Now()
	}
	dnsLock.Unlock()
}
func GetFirstIP(host string, proxyHost string) net.IP {
	key := ""
	if proxyHost == "" {
		key = "_default_" + host
	} else {
		key = proxyHost + "|" + host
	}
	var ip net.IP
	dnsLock.Lock()
	ips := dnsList[key]
	if ips != nil {
		ip = ips.first
		ips.time = time.Now()
	}
	dnsLock.Unlock()
	return ip
}
func LookupIP(host string, proxy string, Dial func(network, address string) (net.Conn, error)) ([]net.IP, error) {
	ips, err := lookupIP(host, proxy, Dial, "ip4")
	if len(ips) > 0 {
		return ips, err
	}
	return lookupIP(host, proxy, Dial, "ip")
}
func lookupIP(host string, proxy string, Dial func(network, address string) (net.Conn, error), Net string) ([]net.IP, error) {
	if proxy == "" {
		key := "_default_" + host
		dnsLock.Lock()
		ips := dnsList[key]
		if ips != nil {
			ips.time = time.Now()
			dnsLock.Unlock()
			return ips.ips, nil
		}
		dnsLock.Unlock()
		_ips, _err := net.LookupIP(host)
		if len(_ips) > 0 {
			t := &rsIps{ips: _ips, time: time.Now()}
			dnsLock.Lock()
			dnsList[key] = t
			dnsLock.Unlock()
		}
		return _ips, _err
	}
	key := proxy + "|" + host
	dnsLock.Lock()
	ips := dnsList[key]
	if ips != nil {
		ips.time = time.Now()
		dnsLock.Unlock()
		return ips.ips, nil
	}
	resolver := dnsTools[proxy]
	if resolver == nil {
		t := &tools{rs: newResolver(proxy, Dial)}
		dnsTools[proxy] = t
	}
	resolver = dnsTools[proxy]
	resolver.time = time.Now()
	dnsLock.Unlock()
	_ips, _err := resolver.rs.LookupIP(context.Background(), Net, host)
	if len(_ips) > 0 {
		t := &rsIps{ips: _ips, time: time.Now()}
		dnsLock.Lock()
		dnsList[key] = t
		dnsLock.Unlock()
	}
	return _ips, _err
}
