package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/errorImpl"
)

func CookiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write(cache.GetCookies())
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		cookies := cache.Cookies{}
		err := json.NewDecoder(r.Body).Decode(&cookies)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, err)
			return
		}
		val, err := json.Marshal(cookies.Values)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, err)
			return
		}
		err = cache.SetCookies(cookies.Key, val)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	}
}
