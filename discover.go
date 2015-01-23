package main

import (
	"log"
	"net/url"

	"github.com/robtuley/etcdwatch"
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

	changeC, errorC := etcdwatch.Key("/hosts")

	go func() {
		for {
			log.Println(<-changeC)
		}
	}()

	go func() {
		for {
			log.Println(<-errorC)
		}
	}()

	return routeC
}
