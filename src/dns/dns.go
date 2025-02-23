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
var dnsServer = "localhost" //223.5.5.5:853  阿里云公共DNS解析服务器
const dnsServerLocal = "localhost"

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
				if proxy == "" {
					return dialer.DialContext(context, network_, address)
				}
				//使用代理进行查询，代理仅支持TCP
				return Dial("tcp", address)
			}
			_tlsTCP := strings.HasSuffix(_dnsServer, ":853")

			if _tlsTCP {
				if proxy == "" {
					conn, err = dialer.DialContext(context, network_, _dnsServer)
				} else {
					//使用代理连接到自定义DNS服务器，代理仅支持TCP
					conn, err = Dial("tcp", _dnsServer)
				}
				if err != nil {
					return nil, err
				}
				_ = conn.(*net.TCPConn).SetKeepAlive(true)
				_ = conn.(*net.TCPConn).SetKeepAlivePeriod(10 * time.Second)
				return tls.Client(conn, dnsConfig), nil
			}

			if proxy == "" {
				return dialer.DialContext(context, network_, _dnsServer)
			}
			//使用代理连接到自定义DNS服务器，代理仅支持TCP
			return Dial("tcp", _dnsServer)
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
func GetDnsServer() string {
	return dnsServer
}
func SetFirstIP(host string, proxyHost string, ip net.IP) {
	key := ""
	if proxyHost == "" {
		key = "_default_" + host
	} else {
		key = proxyHost + "|" + host
	}
	dnsLock.Lock()
	if ip == nil {
		delete(dnsList, key)
	} else {
		ips := dnsList[key]
		if ips != nil {
			ips.first = ip
			ips.time = time.Now()
		}
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

// deepCopyIPs 进行深拷贝
func deepCopyIPs(src []net.IP) []net.IP {
	dst := make([]net.IP, len(src)) // 创建与源数组相同大小的目标切片
	for i, ip := range src {
		if ip != nil { // 避免空值
			dst[i] = make(net.IP, len(ip)) // 为每个 IP 重新分配内存
			copy(dst[i], ip)               // 复制数据
		}
	}
	return dst
}
func LookupIP(host string, proxy string, Dial func(network, address string) (net.Conn, error)) ([]net.IP, error) {
	dnsLock.Lock()
	if dnsServer == dnsServerLocal {
		dnsLock.Unlock()
		return localLookupIP(host, proxy)
	}
	dnsLock.Unlock()
	ips, err := lookupIP(host, proxy, Dial, "ip4")
	if len(ips) > 0 {
		return deepCopyIPs(ips), err
	}
	ips, err = lookupIP(host, proxy, Dial, "ip")
	if len(ips) > 0 {
		return deepCopyIPs(ips), err
	}
	if proxy == "" {
		return deepCopyIPs(ips), err
	}
	//如果远程没有解析成功,则使用本地DNS解析一次
	return localLookupIP(host, proxy)
}
func lookupIP(host string, proxy string, Dial func(network, address string) (net.Conn, error), Net string) ([]net.IP, error) {
	if proxy == "" {
		return localLookupIP(host, proxy)
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
	_ips_, _err := resolver.rs.LookupIP(context.Background(), Net, host)
	_ips := deepCopyIPs(_ips_)
	if len(_ips) > 0 {
		t := &rsIps{ips: _ips, time: time.Now()}
		dnsLock.Lock()
		dnsList[key] = t
		dnsLock.Unlock()
	}
	return _ips, _err
}
func localLookupIP(host, proxyHost string) ([]net.IP, error) {
	key := ""
	if proxyHost == "" {
		key = "_default_" + host
	} else {
		key = proxyHost + "|" + host
	}
	dnsLock.Lock()
	ips := dnsList[key]
	if ips != nil {
		ips.time = time.Now()
		dnsLock.Unlock()
		return ips.ips, nil
	}
	dnsLock.Unlock()
	_ips, _err := net.LookupIP(host)
	_ips_ := deepCopyIPs(_ips)
	if len(_ips_) > 0 {
		t := &rsIps{ips: _ips_, time: time.Now()}
		dnsLock.Lock()
		dnsList[key] = t
		dnsLock.Unlock()
	}
	return _ips_, _err
}
