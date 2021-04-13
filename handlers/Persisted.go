package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matscus/ammunition/errorImpl"
	"github.com/matscus/ammunition/pool"
)

func PersistedGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	project, ok := params["project"]
	if ok {
		script, ok := params["script"]
		if ok {
			res, err := pool.PersistedPool{Project: project, Script: script}.GetValue()
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			if res == "" {
				errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("chanel is empty"))
				return
			}
			_, err = w.Write([]byte(res))
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusOK, err)
				return
			}
		}
	}
}

//Manage func from create(method post) or update(method put) or delete (method delete) datapool
func PersistedManageHandler(w http.ResponseWriter, r *http.Request) {
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
	datapool := pool.PersistedPool{Project: project, Script: scriptName}
	switch r.Method {
	case http.MethodGet:

	case http.MethodPost:
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Get form uploadFile error :"+err.Error()))
			return
		}
		err = datapool.Create(&file)
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
			err = datapool.AddValues(&file)
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
