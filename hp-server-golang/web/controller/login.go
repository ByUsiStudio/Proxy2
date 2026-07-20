package controller

import (
	"hp-server-lib/bean"
	"hp-server-lib/config"
	"hp-server-lib/service"
	"net/http"
)

type LoginController struct {
	service.UserService
}

func (receiver LoginController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var msg bean.ReqLogin
	err := DecodeBody(r, &msg)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := receiver.Login(msg)
	if login != nil {
		WriteOk(w, login)
		return
	} else {
		WriteError(w, "登录失败，用户名或密码错误")
	}
}

func (receiver LoginController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var msg bean.ReqRegister
	err := DecodeBody(r, &msg)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = receiver.Register(msg)
	if err == nil {
		WriteOk(w, nil)
		return
	}
	WriteError(w, err.Error())
}

func (receiver LoginController) SystemConfigHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}
	config := receiver.GetSystemConfig()
	WriteOk(w, config)
}

func (receiver LoginController) PublicConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := receiver.GetPublicConfig()
	WriteOk(w, config)
}

func (receiver LoginController) UpdateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}

	userIdInt := ParseIntParam(r, "userId")
	statusInt := ParseIntParam(r, "status")

	err := receiver.UpdateStatus(userIdInt, statusInt)
	if err == nil {
		WriteOk(w, nil)
		return
	}
	WriteError(w, err.Error())
}

func (receiver LoginController) UpdateSystemConfigHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}

	var sysConfig config.SystemConfig
	err := DecodeBody(r, &sysConfig)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = receiver.UpdateSystemConfig(sysConfig)
	if err == nil {
		WriteOk(w, nil)
		return
	}
	WriteError(w, err.Error())
}
