package controller

import (
	"hp-server-lib/entity"
	"hp-server-lib/service"
	"net/http"
)

type ForwardController struct {
	service.ForwardService
}

func (receiver ForwardController) Add(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		var msg entity.UserFwdEntity
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg.UserId = &userId
		err = receiver.AddData(msg)
		if err != nil {
			WriteError(w, err.Error())
			return
		}
		WriteOk(w, nil)
	}
}

func (receiver ForwardController) List(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		pageInt, pageSizeInt := ParsePage(r)
		WriteOk(w, receiver.ListData(userId, pageInt, pageSizeInt))
	}
}

func (receiver ForwardController) Del(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		idInt := ParseIntParam(r, "id")
		receiver.RemoveData(idInt)
		WriteOk(w, nil)
	}
}
