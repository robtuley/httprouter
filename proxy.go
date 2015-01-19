package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/robtuley/report"
)

func main() {
	defer report.Drain()
	report.StdOut()

	url_a, err := url.Parse("http://localhost:8001")
	if err != nil {
		report.Action("proxy.endpoint.fail", report.Data{"error": err.Error()})
		return
	}
	proxy_a := httputil.NewSingleHostReverseProxy(url_a)

	url_b, err := url.Parse("http://localhost:8002")
	proxy_b := httputil.NewSingleHostReverseProxy(url_b)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		switch {
		case r.Host == "a.example.com:8000":
			proxy_a.ServeHTTP(w, r)
		case r.Host == "b.example.com:8000":
			proxy_b.ServeHTTP(w, r)
		default:
			http.Error(w, "No route for "+r.Host, http.StatusNotFound)
		}

	})

	report.Info("proxy.start", report.Data{"port": 8000})
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		report.Action("proxy.startup.fail", report.Data{"error": err.Error()})
	}
}
