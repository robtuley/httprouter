package main

import (
	"log"
	"net/url"
)

type Route struct {
	Host string
	UrlC chan []url.URL
}

var discoverC chan Route

func init() {
	discoverC = make(chan Route)
	addC := discoverNewHosts("http://127.0.0.1:4001/v2/keys/hosts")

	go func() {
		for {
			host := <-addC
			log.Println("addC:>", host)
		}
	}()

}

func discoverNewHosts(keyUrl string) chan string {
	addC := make(chan string)
	nodeC := longPollForJson(keyUrl)

	go func() {
		for {
			node := <-nodeC
			log.Println(node)
			addC <- node.Key
		}
	}()

	return addC
}
