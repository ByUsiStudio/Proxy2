package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticController 静态资源控制器
// 不再通过 go:embed 集成前端产物，改为运行时从本地 static 目录读取，
// 这样编译产物不再依赖前端文件，部署时按需将前端构建产物放到同级 static 目录即可。
type StaticController struct{}

func (s StaticController) Static(w http.ResponseWriter, r *http.Request) {
	// 确保路径处理正确（去除可能的前缀，基于实际项目调整）
	filePath := r.URL.Path
	// 拼接静态文件目录（静态文件都在 "static" 目录下）
	fullPath := filepath.Join("static", filePath)

	// 尝试读取请求的文件
	data, err := os.ReadFile(fullPath)
	if err != nil {
		// 文件不存在时，尝试返回 index.html（支持前端路由）
		indexData, indexErr := os.ReadFile(filepath.Join("static", "index.html"))
		if indexErr != nil {
			// 如果 index.html 也不存在，才返回 404
			http.NotFound(w, r)
			return
		}
		// 返回 index.html 内容
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexData)
		return
	}

	// 文件存在时，设置正确的 Content-Type 并返回
	switch {
	case strings.HasSuffix(filePath, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(filePath, ".css"):
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case strings.HasSuffix(filePath, ".js"):
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case strings.HasSuffix(filePath, ".json"):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case strings.HasSuffix(filePath, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(filePath, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	w.Write(data)
}
