package controller

import (
	"encoding/json"
	"hp-server-lib/bean"
	"hp-server-lib/util"
	"net/http"
	"strconv"
	"strings"
)

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	token := r.Header.Get("token")
	userId, _, _, err := util.DecodeToken(token)
	if err != nil {
		WriteJSON(w, bean.ResErrorCode(-2, "用户权限校验失败"))
		return 0, err
	}
	return userId, nil
}

func CheckAdmin(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("token")
	_, role, _, err := util.DecodeToken(token)
	if err != nil || !strings.EqualFold(role, "ADMIN") {
		WriteJSON(w, bean.ResErrorCode(-2, "用户权限校验失败"))
		return false
	}
	return true
}

func GetRole(w http.ResponseWriter, r *http.Request) (string, error) {
	token := r.Header.Get("token")
	_, role, _, err := util.DecodeToken(token)
	if err != nil {
		WriteJSON(w, bean.ResErrorCode(-2, "用户权限校验失败"))
		return "", err
	}
	return role, nil
}

func ParsePage(r *http.Request) (int, int) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("current"))
	pageSize, _ := strconv.Atoi(queryParams.Get("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func ParseIntParam(r *http.Request, key string) int {
	val, _ := strconv.Atoi(r.URL.Query().Get(key))
	return val
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(data)
}

func WriteOk(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, bean.ResOk(data))
}

func WriteError(w http.ResponseWriter, msg string) {
	WriteJSON(w, bean.ResError(msg))
}

func WriteErrorCode(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, bean.ResErrorCode(code, msg))
}

func DecodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
