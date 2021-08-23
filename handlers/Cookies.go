package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/matscus/ammunition/errorImpl"
	"github.com/matscus/ammunition/pool"
)

func CookiesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		cookies := pool.Data{}
		log.Println("start")
		cookies.Value = string(pool.GetCookies())
		log.Println(cookies.Value)
		err := json.NewEncoder(w).Encode(cookies)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		cookies := pool.Data{}
		err := json.NewDecoder(r.Body).Decode(&cookies)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, err)
			return
		}
		log.Println(cookies.Key, cookies.Value)
		err = pool.SetCookies(cookies.Key, cookies.Value)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	}
}
