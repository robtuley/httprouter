package proxy

import (
	"net/http"
)

var proxy roundRobin

func init() {
	proxy = newRoundRobin(unavailableHandler())

	a := http.HandlerFunc(serveA)
	b := http.HandlerFunc(serveB)
	c := http.HandlerFunc(serveC)

	proxy.Add(&a)
	proxy.Add(&b)
	proxy.Add(&c)
}

func Domain(domain string) http.Handler {
	return proxy.Choose()
}

func serveA(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("A"))
}

func serveB(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("B"))
}

func serveC(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("C"))
}
