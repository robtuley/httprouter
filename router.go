package main

import (
	"net/http"

	"github.com/robtuley/httprouter/proxy"
	"github.com/robtuley/report"
)

func main() {
	defer report.Drain()
	report.StdOut()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.Domain("test").ServeHTTP(w, r)
	})

	report.Info("router.start", report.Data{"port": 8080})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		report.Action("router.start.fail", report.Data{"error": err.Error()})
	}

	/*
		routeC := discover.Etcd("/domains")
		go func() {
			for {
				<-routeC
			}
		}()

		// host, _, _ := net.SplitHostPort(req.RemoteAddr)

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
	*/
}
