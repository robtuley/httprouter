package main

import (
	"log"
	"net/url"
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

	respC, err := longPollForKeyChanges("/hosts")
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			log.Println(<-respC)
		}
	}()

	return routeC
}
