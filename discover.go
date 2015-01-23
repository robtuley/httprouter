package main

import (
	"log"
	"net/url"

	"github.com/robtuley/proxy/etcd"
	"github.com/robtuley/report"
)

type Route struct {
	Host    string
	URL     url.URL
	SignalC chan int
}

var discover struct {
	C chan Route
}

func init() {
	discover.C = discoverRoutesFromEtcD()
}

func discoverRoutesFromEtcD() chan Route {
	routeC := make(chan Route)

	respC, reportC := etcd.Watch("/hosts")

	go func() {
		for {
			log.Println(<-respC)
		}
	}()

	go func() {
		for {
			report.Action("etcd.response", <-reportC)
		}
	}()

	return routeC
}
