package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/errorImpl"
)

func CookiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cookies := cache.Data{}
		log.Println("start")
		cookies.Value = string(cache.GetCookies())
		log.Println(cookies.Value)
		err := json.NewEncoder(w).Encode(cookies)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		cookies := cache.Data{}
		err := json.NewDecoder(r.Body).Decode(&cookies)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, err)
			return
		}
		err = cache.SetCookies(cookies.Key, cookies.Value)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	}
}
