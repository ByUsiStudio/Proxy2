package controller

import (
	"hp-server-lib/service"
	"net/http"
)

type MonitorController struct {
	service.MonitorService
}

func (receiver MonitorController) List(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserId(w, r)
	if err == nil {
		data := receiver.ListData(id)
		if data != nil {
			WriteOk(w, data)
			return
		} else {
			WriteError(w, "登陆失败")
		}

	}
}

func (receiver MonitorController) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	data := receiver.DetailData(id)
	if data != nil {
		WriteOk(w, data)
		return
	} else {
		WriteError(w, "登陆失败")
	}
}
