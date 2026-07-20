package controller

import (
	"hp-server-lib/entity"
	"hp-server-lib/service"
	"net/http"
)

type ClientUserController struct {
	service.UserCustomService
}

func (receiver ClientUserController) Add(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}
	var msg entity.UserCustomEntity
	err := DecodeBody(r, &msg)
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = receiver.AddData(msg)
	if err != nil {
		WriteError(w, err.Error())
		return
	}
	WriteOk(w, nil)
}

func (receiver ClientUserController) List(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}
	pageInt, pageSizeInt := ParsePage(r)
	WriteOk(w, receiver.ListData(pageInt, pageSizeInt))
}

func (receiver ClientUserController) Del(w http.ResponseWriter, r *http.Request) {
	if !CheckAdmin(w, r) {
		return
	}
	idInt := ParseIntParam(r, "id")
	receiver.RemoveData(idInt)
	WriteOk(w, nil)
}
