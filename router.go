package main

import (
	"flag"
	"net"
	"net/http"
	"strconv"

	"github.com/robtuley/httprouter/proxy"
	"github.com/robtuley/report"
)

const (
	port = 8080
)

func main() {
	var log string
	flag.StringVar(&log, "log", "./router.log", "logfile path")
	flag.Parse()

	defer report.Drain()
	report.File(log)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tick := report.Tick()

		proxy.Domain(r.Host).ServeHTTP(w, r)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		report.Tock(tick, "router.request", report.Data{
			"host": r.Host,
			"path": r.URL.Path,
			"ua":   r.UserAgent(),
			"ip":   ip,
		})
	})

	report.Info("router.start", report.Data{"port": port})
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		report.Action("router.error", report.Data{"error": err.Error()})
	}
}
