package http

import (
	"crypto/tls"
	"fmt"
	"hp-server-lib/config"
	"hp-server-lib/entity"
	"hp-server-lib/log"
	"hp-server-lib/service"
	"net/http"
	"os"
	"strconv"
)

func getCertificateAndTargetForDomain(domain string) (*tls.Certificate, error) {
	value, ok := service.DOMAIN_INFO.Load(domain)
	if !ok {
		return nil, fmt.Errorf("域名找不到证书: %s", domain)
	}
	info := value.(*entity.UserDomainEntity)
	certificate, err := tls.X509KeyPair([]byte(info.CertificateContent), []byte(info.CertificateKey))
	if err != nil {
		return nil, fmt.Errorf("证书解析失败 %s: %v", domain, err)
	}
	return &certificate, nil
}

func StartHttpsServer() {
	mux := http.NewServeMux()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			domain := clientHello.ServerName
			cert, err := getCertificateAndTargetForDomain(domain)
			if err != nil {
				return nil, err
			}
			return cert, nil
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r)
	})
	port := config.ConfigData.Tunnel.HttpsPort
	if port <= 0 {
		port = 443
	}
	server := &http.Server{
		Addr:      ":" + strconv.Itoa(port),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
	log.Info("HTTPS代理服务启动")
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Errorf("HTTPS代理服务启动失败: %v", err)
		os.Exit(1)
	}
}
