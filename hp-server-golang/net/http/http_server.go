package http

import (
	"hp-server-lib/config"
	"hp-server-lib/log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func StartHttpServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r)
	})
	port := config.ConfigData.Tunnel.HttpPort
	if port <= 0 {
		port = 80
	}
	log.Info("HTTP代理服务启动")
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		log.Errorf("HTTP代理服务启动失败: %v", err)
		os.Exit(1)
	}
}
