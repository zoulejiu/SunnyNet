package Api

import (
	"SunnyNet/project/src/Certificate"
	"SunnyNet/project/src/HttpCertificate"
	"crypto/x509"
)

// AddHttpCertificate 创建 Http证书管理器 对象 实现指定Host使用指定证书
func AddHttpCertificate(host string, CertManagerId int, Rules uint8) bool {
	HttpCertificate.Lock.Lock()
	defer HttpCertificate.Lock.Unlock()
	w := Certificate.LoadCertificateContext(CertManagerId)
	if w == nil {
		return false
	}
	ca := w.ExportCA()
	key := w.ExportKEY()
	var RootCAs *x509.CertPool
	if w.Tls != nil {
		if w.Tls.ClientCAs != nil {
			RootCAs = w.Tls.ClientCAs
		}
	}
	if ca == "" || key == "" && RootCAs != nil {
		c := &HttpCertificate.CertificateRequestManager{Rules: Rules}
		c.AddRootCAs(RootCAs)
		HttpCertificate.Map[HttpCertificate.ParsingHost(host)] = c
		return true
	}
	if ca == "" {
		return false
	}
	if key == "" {
		return false
	}
	c := &HttpCertificate.CertificateRequestManager{Rules: Rules}
	if c.Load(ca, key) {
		c.AddRootCAs(RootCAs)
		HttpCertificate.Map[HttpCertificate.ParsingHost(host)] = c
		return true
	}
	return false
}

// DelHttpCertificate 删除 Http证书管理器 对象
func DelHttpCertificate(host string) {
	HttpCertificate.Lock.Lock()
	delete(HttpCertificate.Map, HttpCertificate.ParsingHost(host))
	HttpCertificate.Lock.Unlock()
}
