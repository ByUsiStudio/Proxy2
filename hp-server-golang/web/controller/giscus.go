package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type GiscusController struct {
}

func (receiver GiscusController) Token(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	if session == "" {
		WriteError(w, "失败")
		return
	}

	url := "https://giscus.app/api/oauth/token"
	payload := map[string]string{"session": session}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		WriteError(w, "请求创建失败")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		WriteError(w, "请求失败")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		WriteError(w, "读取响应失败")
		return
	}
	WriteOk(w, string(body))
}
