package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/robtuley/httprouter/proxy"
	"github.com/robtuley/httpserver"
	"github.com/robtuley/report"
)

const (
	PORT = 8080
)

func main() {
	defer report.Drain()
	report.Global(report.Data{"service": "httprouter"})

	var logfile, logurl, etcdurl, etcdkey string
	flag.StringVar(&logfile, "logfile", "", "log file path e.g. /var/log/xxx.log")
	flag.StringVar(&logurl, "logurl", "", "log URL where data is POSTed to")
	flag.StringVar(&etcdurl, "etcdurl", "http://127.0.0.1:4001", "Etcd URL")
	flag.StringVar(&etcdkey, "etcdkey", "/domains", "etcd key to discovery routes from")
	flag.Parse()

	initLogOutput(logfile, logurl)
	proxy.Listen(etcdurl, etcdkey)

	report.Info("service.start", report.Data{"port": PORT})
	service, err := httpserver.New(PORT)
	if err != nil {
		report.Action("service.start.error", report.Data{"error": err.Error()})
		return
	}

	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	service.Start()

	sig := service.WaitStop()
	report.Info("service.stopped", report.Data{"port": PORT, "signal": sig})
	if err = service.Err(); err != nil {
		report.Action("service.error", report.Data{"error": err.Error()})
	}
}

func initLogOutput(logfile string, logurl string) {
	switch {
	case len(logfile) > 0:
		report.File(logfile)
	case len(logurl) > 0:
		report.BatchPostToUrl(logurl)
	default:
		report.StdOut()
	}
}
