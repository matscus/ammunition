package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matscus/ammunition/errorImpl"
	"github.com/matscus/ammunition/pool"
)

func PersistedGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	project, ok := params["project"]
	if ok {
		script, ok := params["name"]
		if ok {
			res, err := pool.PersistedPool{Project: project, Name: script}.GetValue()
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
	name := r.FormValue("name")
	if name == "" {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form scriptname is nil"))
		return
	}
	bufferLenStr := r.FormValue("bufferlen")
	if name == "" {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form bufferlen is nil"))
		return
	}
	bufferLen, err := strconv.Atoi(bufferLenStr)
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Atoi buffer len error: "+err.Error()))
		return
	}
	workersStr := r.FormValue("workers")
	if name == "" {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form bufferlen is nil"))
		return
	}
	workers, err := strconv.Atoi(workersStr)
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Atoi workers error: "+err.Error()))
		return
	}
	pool := pool.PersistedPool{Project: project, Name: name, BufferLen: bufferLen, Workers: workers}
	switch r.Method {
	case http.MethodPost:
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Get form uploadFile error :"+err.Error()))
			return
		}
		err = pool.Create(&file)
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
		if action == "" {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Form action is nil"))
			return
		}
		switch action {
		case "update":
			err = pool.Update(&file)
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Update datapool error: "+err.Error()))
				return
			}
		case "add":
			err = pool.AddValues(&file)
			if err != nil {
				errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Add values from datapool error: "+err.Error()))
				return
			}
		default:
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("invalid value. possible values Update or Add"))
			return
		}
	case http.MethodDelete:
		err := pool.Delete()
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, errors.New("Delete datapool error: "+err.Error()))
			return
		}
	}

}
