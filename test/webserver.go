package main

import (
	"flag"
	"io"
	"net/http"
	"strings"

	"github.com/robtuley/httpserver"
)

func main() {
	var label string
	var port int
	flag.StringVar(&label, "label", "Demo App", "label to echo from webserver")
	flag.IntVar(&port, "port", 8001, "port to run on")
	flag.Parse()

	service, err := httpserver.New(port)
	if err != nil {
		panic(err)
	}

	service.HandleFunc("/dump", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(res, "HEADERS:")
		for k, v := range req.Header {
			io.WriteString(res, "\n"+k+": "+strings.Join(v, " "))
		}
	})

	service.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(res, label)
	})

	service.Start()
	service.WaitStop()
}
