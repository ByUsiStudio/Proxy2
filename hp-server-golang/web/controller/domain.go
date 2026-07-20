package controller

import (
	"hp-server-lib/entity"
	"hp-server-lib/service"
	"net/http"
)

type DomainController struct {
	service.DomainService
}

func (receiver DomainController) GetDomainList(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserId(w, r)
	if err == nil {
		pageInt, pageSizeInt := ParsePage(r)
		keyword := r.URL.Query().Get("keyword")
		WriteOk(w, receiver.DomainList(id, pageInt, pageSizeInt, keyword))
	}
}

func (receiver DomainController) RemoveDomain(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		idInt := ParseIntParam(r, "id")
		WriteOk(w, receiver.RemoveData(idInt))
	}
}

func (receiver DomainController) Query(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		keyword := r.URL.Query().Get("keyword")
		WriteOk(w, receiver.DomainListByKey(userId, keyword))
	}
}

func (receiver DomainController) Gen(w http.ResponseWriter, r *http.Request) {
	_, err := GetUserId(w, r)
	if err == nil {
		idInt := ParseIntParam(r, "id")
		WriteOk(w, receiver.GenSsl(false, idInt))
	}
}

func (receiver DomainController) Add(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err == nil {
		var msg entity.UserDomainEntity
		err := DecodeBody(r, &msg)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg.UserId = &userId
		err = receiver.AddData(msg)
		if err == nil {
			WriteOk(w, nil)
			return
		}
		WriteError(w, err.Error())
	}
}
