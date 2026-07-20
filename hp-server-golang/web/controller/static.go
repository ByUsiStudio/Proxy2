package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const staticDir = "static"

type StaticController struct{}

func (s StaticController) Static(w http.ResponseWriter, r *http.Request) {
	requestPath := strings.TrimPrefix(r.URL.Path, "/")
	fullPath := filepath.Join(staticDir, requestPath)

	info, err := os.Stat(fullPath)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/user/") || strings.HasPrefix(r.URL.Path, "/client/") {
		http.NotFound(w, r)
		return
	}

	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, indexPath)
}
