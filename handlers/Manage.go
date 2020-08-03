package handlers

import (
	"errors"
	"net/http"

	"github.com/matscus/ammunition/datapool"
	"github.com/matscus/ammunition/errorImpl"
)

//Manage func from create(method post) or update(method put) or delete (method delete) datapool
func Manage(w http.ResponseWriter, r *http.Request) {
	project := r.FormValue("project")
	if project == "" {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form project is nil"))
		return
	}
	scriptName := r.FormValue("scriptname")
	if scriptName == "" {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form scriptname is nil"))
		return
	}
	datapool := datapool.Datapool{ProjectName: project, ScriptName: scriptName}
	switch r.Method {
	case http.MethodPost:
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Get form uploadFile error :"+err.Error()))
			return
		}
		err = datapool.New(&file)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Create datapool error: "+err.Error()))
			return
		}
	case http.MethodPut:
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Get form uploadFile error :"+err.Error()))
			return
		}
		action := r.FormValue("action")
		if scriptName == "" {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form action is nil"))
			return
		}
		switch action {
		case "update":
			err = datapool.Update(&file)
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Update datapool error: "+err.Error()))
				return
			}
		case "add":
			err = datapool.Add(&file)
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Add values from datapool error: "+err.Error()))
				return
			}
		default:
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("invalid value. possible values Update or Add"))
			return
		}
	case http.MethodDelete:
		err := datapool.Delete()
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Delete datapool error: "+err.Error()))
			return
		}
	}

}
