package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/errorImpl"
)

type kv struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

func KVHahdler(w http.ResponseWriter, r *http.Request) {
	kv := kv{}
	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		errorImpl.WriteHTTPError(w, http.StatusOK, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		res, err := cache.KV.Get(kv.Key)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusOK, err)
			return
		}
		kv.Value = string(res)
		err = json.NewEncoder(w).Encode(kv)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		err = cache.KV.Set(kv.Key, []byte(kv.Value))
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodDelete:
		err = cache.KV.Delete(kv.Key)
		if err != nil {
			errorImpl.WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	}
}
