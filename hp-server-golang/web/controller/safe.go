package controller

import (
	"hp-server-lib/entity"
	"hp-server-lib/service"
	"net/http"
)

type SafeController struct {
	service.UserSafeService
}

func (receiver SafeController) Add(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		var msg entity.UserSafeEntity
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = receiver.AddData(userId, msg)
		if err != nil {
			WriteError(w, err.Error())
			return
		}
		WriteOk(w, nil)
	}
}

func (receiver SafeController) List(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		pageInt, pageSizeInt := ParsePage(r)
		WriteOk(w, receiver.ListData(userId, pageInt, pageSizeInt))
	}
}

func (receiver SafeController) Del(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		idInt := ParseIntParam(r, "id")
		receiver.RemoveData(idInt)
		WriteOk(w, nil)
	}
}

func (receiver SafeController) Query(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		keyword := r.URL.Query().Get("keyword")
		WriteOk(w, receiver.SafeListByKey(userId, keyword))
	}
}
