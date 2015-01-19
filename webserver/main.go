package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	var label string
	var port int
	flag.StringVar(&label, "label", "Demo App", "label to echo from webserver")
	flag.IntVar(&port, "port", 8001, "port to run on")
	flag.Parse()

	log.Println("started:> port:", port, " label:", label)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, label)
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
