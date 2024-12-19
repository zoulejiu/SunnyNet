//go:build windows
// +build windows

package CrossCompiled

import "C"
import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Trisia/gosysproxy"
	"github.com/qtgolang/SunnyNet/src/iphlpapi"
	NFapi "github.com/qtgolang/SunnyNet/src/nfapi"
	"github.com/qtgolang/SunnyNet/src/public"
	"golang.org/x/sys/windows"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func NFapi_SunnyPointer(a ...uintptr) uintptr {
	if len(a) > 0 {
		NFapi.SunnyPointer = a[0]
	}
	return NFapi.SunnyPointer
}
func NFapi_IsInit(a ...bool) bool {
	if len(a) > 0 {
		NFapi.IsInit = a[0]
	}
	return NFapi.IsInit
}
func NFapi_ProcessPortInt(a ...uint16) uint16 {
	if len(a) > 0 {
		NFapi.ProcessPortInt = a[0]
	}
	return NFapi.ProcessPortInt
}
func NFapi_ApiInit() bool {
	return NFapi.ApiInit()
}
func NFapi_MessageBox(caption, text string, style uintptr) (result int) {
	return NFapi.MessageBox(caption, text, style)
}
func NFapi_SetHookProcess(open bool) {
	NFapi.SetHookProcess(open)
}
func NFapi_ClosePidTCP(pid int) {
	NFapi.ClosePidTCP(pid)
}
func NFapi_DelName(u string) {
	NFapi.DelName(u)
}
func NFapi_AddName(u string) {
	NFapi.AddName(u)
}
func NFapi_DelPid(pid uint32) {
	NFapi.DelPid(pid)
}
func NFapi_AddPid(pid uint32) {
	NFapi.AddPid(pid)
}
func NFapi_CloseNameTCP(u string) {
	NFapi.CloseNameTCP(u)
}
func NFapi_CancelAll() {
	NFapi.CancelAll()
}
func NFapi_DelTcpConnectInfo(U uint16) {
	NFapi.DelTcpConnectInfo(U)
}
func NFapi_GetTcpConnectInfo(U uint16) *NFapi.ProcessInfo {
	return NFapi.GetTcpConnectInfo(U)
}

func NFapi_API_NfTcpClose(U uint64) {
	NFapi.Api.NfTcpClose(U)
}
func NFapi_UdpSendReceiveFunc(udp func(Type int, Theoni int64, pid uint32, LocalAddress, RemoteAddress string, data []byte) []byte) func(Type int, Theoni int64, pid uint32, LocalAddress, RemoteAddress string, data []byte) []byte {
	NFapi.UdpSendReceiveFunc = udp
	return NFapi.UdpSendReceiveFunc
}
func NFapi_Api_NfUdpPostSend(id uint64, remoteAddress *NFapi.SockaddrInx, buf []byte, option *NFapi.NF_UDP_OPTIONS) (NFapi.NF_STATUS, error) {
	return NFapi.Api.NfUdpPostSend(id, remoteAddress, buf, option)
}

func SetIeProxy(Off bool, Port int) bool {
	// "github.com/Tri sia/gos ysp roxy"
	if Off {
		_ = gosysproxy.Off()
		return true
	}
	ies := "127.0.0.1:" + strconv.Itoa(Port)
	_ = gosysproxy.SetGlobalProxy("http="+ies+";https="+ies, "")
	return true
}

// InstallCert 安装证书 将证书安装到Windows系统内
func InstallCert(certificates []byte) (res string) {

	defer func() {
		CertificateName := public.GetCertificateName(certificates)
		if CertificateName != "" && isInstallSunnyNetCertificates(CertificateName) {
			res = "already in store"
		}
	}()
	tempDir := os.TempDir()
	err := public.WriteBytesToFile(certificates, tempDir+"\\SunnyNet.crt")
	if err != nil {
		return err.Error()
	}
	var args []string
	args = append(args, "-addstore")
	args = append(args, "root")
	args = append(args, tempDir+"\\SunnyNet.crt")
	defer func() { _ = public.RemoveFile(tempDir + "\\SunnyNet.crt") }()
	cmd := exec.Command("certutil", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
	var Buff bytes.Buffer
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadBytes('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		Buff.Write(line)
	}
	return Buff.String()
}

// InstallCert2 有感安装,会提示对话框安装
func InstallCert2(certPEM []byte) string {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return "Invalid certificate"
	}
	storeName, _ := syscall.UTF16PtrFromString("ROOT")
	store, err := windows.CertOpenStore(windows.CERT_STORE_PROV_SYSTEM, 0, 0, windows.CERT_SYSTEM_STORE_CURRENT_USER, uintptr(unsafe.Pointer(storeName)))
	if err != nil {
		return fmt.Sprintf("failed to open certificate store: %v", err)
	}
	defer windows.CertCloseStore(store, 0)
	certContext, err := windows.CertCreateCertificateContext(
		windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		&block.Bytes[0],
		uint32(len(block.Bytes)),
	)
	if err != nil {
		return fmt.Sprintf("failed to create certificate context: %v", err)
	}
	defer windows.CertFreeCertificateContext(certContext)
	// 将证书添加到存储区
	if windows.CertAddCertificateContextToStore(
		store,
		certContext,
		windows.CERT_STORE_ADD_USE_EXISTING,
		nil,
	) != nil {
		return "安装证书失败：用户未授权安装证书"
	}
	return "already in store"
}

const (
	CERT_SYSTEM_STORE_CURRENT_USER  = uint32(1 << 16) // 当前用户证书存储
	CERT_SYSTEM_STORE_LOCAL_MACHINE = uint32(2 << 16) // 本地计算机证书存储
)

// 检查是否安装了包含 "SunnyNet" 的证书
func isInstallSunnyNetCertificates(CertificateName string) bool {
	if !_isInstallSunnyNetCertificates(CERT_SYSTEM_STORE_CURRENT_USER, CertificateName) {
		return _isInstallSunnyNetCertificates(CERT_SYSTEM_STORE_LOCAL_MACHINE, CertificateName)
	}
	return true
}

func _isInstallSunnyNetCertificates(CERT uint32, CertificateName string) bool {
	// 将 "ROOT" 转换为 UTF-16 指针
	storeName, err := syscall.UTF16PtrFromString("ROOT")
	if err != nil {
		return false // 转换失败，返回 false
	}

	// 打开当前用户的根证书存储
	store, err := windows.CertOpenStore(windows.CERT_STORE_PROV_SYSTEM, 0, 0, CERT, uintptr(unsafe.Pointer(storeName)))
	if store == 0 || err != nil {
		return false // 打开证书存储失败，返回 false
	}
	defer windows.CertCloseStore(store, 0) // 确保在函数结束时关闭证书存储

	var cert *windows.CertContext // 声明证书上下文
	for {
		// 枚举证书存储中的证书
		cert, _ = windows.CertEnumCertificatesInStore(store, cert)
		if cert == nil {
			break // 如果没有更多证书，退出循环
		}
		// 获取证书的字节数据
		certBytes := (*[1 << 20]byte)(unsafe.Pointer(cert.EncodedCert))[:cert.Length:cert.Length]
		// 解析证书
		parsedCert, er := x509.ParseCertificate(certBytes)
		if er != nil {
			continue // 如果解析失败，继续下一个证书
		}
		// 检查证书的主题名称是否包含 "CertificateName"
		if strings.Contains(parsedCert.Subject.CommonName, CertificateName) {
			return true // 找到匹配的证书，返回 true
		}
	}

	return false // 未找到匹配的证书，返回 false
}
func SetNetworkConnectNumber() {
	//https://blog.csdn.net/PYJcsdn/article/details/126251054
	//尽量避免这个问题
	var args []string
	args = append(args, "int")
	args = append(args, "ipv4")
	args = append(args, "set")
	args = append(args, "dynamicport")
	args = append(args, "tcp")
	args = append(args, "start=10000")
	args = append(args, "num=55000")
	NFapi.ExecCommand("netsh", args)
	var args1 []string
	args1 = append(args1, "int")
	args1 = append(args1, "ipv6")
	args1 = append(args1, "set")
	args1 = append(args1, "dynamicport")
	args1 = append(args1, "tcp")
	args1 = append(args1, "start=10000")
	args1 = append(args1, "num=55000")
	NFapi.ExecCommand("netsh", args1)
}

// CloseCurrentSocket  关闭指定进程的所有TCP连接
func CloseCurrentSocket(PID int, ulAf uint) {
	iphlpapi.CloseCurrentSocket(PID, ulAf)
}
