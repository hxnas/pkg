package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hxnas/pkg/lod"
)

func Blob[T ~[]byte | ~string](body T, contentType string, status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		w.Write([]byte(body))
	}
}

func Redirect(to string, permanent ...bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, lod.Iif(lod.Select(permanent...), http.StatusPermanentRedirect, http.StatusTemporaryRedirect))
	}
}

func JSON(value any, status int) http.HandlerFunc {
	data, err := json.Marshal(value)
	if err != nil {
		return Blob(fmt.Sprintf(`{"code": 500, "err": %s}`, strconv.Quote(err.Error())), "application/json", 500)
	}
	return Blob(data, "application/json", status)
}
