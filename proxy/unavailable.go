package proxy

import (
	"net/http"
)

func unavailableHandler() http.Handler {
	h := http.HandlerFunc(serve503)
	return &h
}

func serve503(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("This site is currently unavailable, please try again later."))
}
