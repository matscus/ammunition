package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matscus/ammunition/datapool"
	"github.com/matscus/ammunition/errorImpl"
)

func GetValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	datapool := datapool.Datapool{}
	err := json.NewDecoder(r.Body).Decode(&datapool)
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, err)
		return
	}
	defer r.Body.Close()
	res, err := datapool.Get()
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, err)
		return
	}
	_, err = w.Write([]byte(res))
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, err)
		return
	}
}
