package controller

import (
	"hp-server-lib/entity"
	"hp-server-lib/service"
	"net/http"
)

type ConfigController struct {
	service.ConfigService
}

func (receiver ConfigController) GetDeviceKey(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserId(w, r)
	if err == nil {
		WriteOk(w, receiver.DeviceKey(id))
	}
}

func (receiver ConfigController) GetConfigList(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserId(w, r)
	if err == nil {
		pageInt, pageSizeInt := ParsePage(r)
		keyword := r.URL.Query().Get("keyword")
		WriteOk(w, receiver.ConfigList(id, pageInt, pageSizeInt, keyword))
	}
}

func (receiver ConfigController) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		configIdInt := ParseIntParam(r, "configId")
		WriteOk(w, receiver.RemoveData(configIdInt))
	}
}

func (receiver ConfigController) Add(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		var msg entity.UserConfigEntity
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = receiver.AddData(msg)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}

func (receiver ConfigController) Keyword(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		keyword := r.URL.Query().Get("keyword")
		data := receiver.KeywordData(userId, keyword)
		if data != nil {
			WriteOk(w, data)
			return
		}
		WriteOk(w, nil)
	}
}

func (receiver ConfigController) RefConfig(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		configIdInt := ParseIntParam(r, "configId")
		err = receiver.RefData(configIdInt)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}

func (receiver ConfigController) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		configIdInt := ParseIntParam(r, "configId")
		err = receiver.ChangeStatusData(configIdInt)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}
