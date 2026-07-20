package controller

import (
	"hp-server-lib/bean"
	"hp-server-lib/service"
	"net/http"
)

type DeviceController struct {
	service.DeviceService
}

func (receiver DeviceController) Add(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		var msg bean.ReqDeviceInfo
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = receiver.AddData(userId, msg)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}

func (receiver DeviceController) Update(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		var msg bean.ReqDeviceInfo
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = receiver.UpdateData(msg)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}

func (receiver DeviceController) List(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		pageInt, pageSizeInt := ParsePage(r)
		WriteOk(w, receiver.ListData(userId, pageInt, pageSizeInt))
	}
}

func (receiver DeviceController) Del(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	err := receiver.RemoveData(deviceId)
	if err == nil {
		WriteOk(w, nil)
		return
	}
	WriteError(w, err.Error())
}

func (receiver DeviceController) Stop(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("deviceId")
	WriteOk(w, receiver.StopData(deviceId))
}
