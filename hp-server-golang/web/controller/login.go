package controller

import (
	"encoding/json"
	"hp-server-lib/bean"
	"hp-server-lib/config"
	"hp-server-lib/service"
	"hp-server-lib/util"
	"net/http"
	"strconv"
	"strings"
)

type LoginController struct {
	service.UserService
}

func (receiver LoginController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var msg bean.ReqLogin
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := receiver.Login(msg)
	if login != nil {
		json.NewEncoder(w).Encode(bean.ResOk(login))
		return
	} else {
		json.NewEncoder(w).Encode(bean.ResError("登录失败，用户名或密码错误"))
	}
}

func (receiver LoginController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var msg bean.ReqRegister
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = receiver.Register(msg)
	if err == nil {
		json.NewEncoder(w).Encode(bean.ResOk(nil))
		return
	}
	json.NewEncoder(w).Encode(bean.ResError(err.Error()))
}

func (receiver LoginController) SystemConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := receiver.GetSystemConfig()
	json.NewEncoder(w).Encode(bean.ResOk(config))
}

func (receiver LoginController) UpdateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	_, role, _, err := util.DecodeToken(token)
	if err != nil || strings.Compare(role, "ADMIN") != 0 {
		json.NewEncoder(w).Encode(bean.ResErrorCode(-2, "用户权限校验失败"))
		return
	}

	queryParams := r.URL.Query()
	userId := queryParams.Get("userId")
	status := queryParams.Get("status")
	userIdInt, _ := strconv.Atoi(userId)
	statusInt, _ := strconv.Atoi(status)

	err = receiver.UpdateStatus(userIdInt, statusInt)
	if err == nil {
		json.NewEncoder(w).Encode(bean.ResOk(nil))
		return
	}
	json.NewEncoder(w).Encode(bean.ResError(err.Error()))
}

func (receiver LoginController) UpdateSystemConfigHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	_, role, _, err := util.DecodeToken(token)
	if err != nil || strings.Compare(role, "ADMIN") != 0 {
		json.NewEncoder(w).Encode(bean.ResErrorCode(-2, "用户权限校验失败"))
		return
	}

	var sysConfig config.SystemConfig
	err = json.NewDecoder(r.Body).Decode(&sysConfig)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = receiver.UpdateSystemConfig(sysConfig)
	if err == nil {
		json.NewEncoder(w).Encode(bean.ResOk(nil))
		return
	}
	json.NewEncoder(w).Encode(bean.ResError(err.Error()))
}
