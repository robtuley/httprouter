package proxy

import (
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/robtuley/httprouter/discover"
	"github.com/robtuley/report"
)

var domainMap struct {
	mu    sync.RWMutex
	proxy map[string]*roundRobin
}

func init() {
	domainMap.proxy = map[string]*roundRobin{}
}

// start listening to a particular etcd key
func Listen(etcdUrl string, etcdKey string) {
	go discoveryToProxyDomainMap(etcdUrl, etcdKey)
}

func Domain(domain string) http.Handler {
	domainMap.mu.RLock()
	defer domainMap.mu.RUnlock()

	p, exists := domainMap.proxy[domain]
	if !exists {
		report.Info("proxy.miss", report.Data{"domain": domain})
		return unavailableHandler()
	}
	return p.Choose()
}

func discoveryToProxyDomainMap(etcdUrl string, etcdKey string) {
	routeC := discover.Etcd(etcdUrl, etcdKey)

	for {
		route, more := <-routeC
		if !more {
			report.Info("proxy.route.closed", report.Data{})
			return
		}
		addRoute(route)
	}
}

func addRoute(route discover.Route) {
	domainMap.mu.Lock()
	defer domainMap.mu.Unlock()

	rrPointer, exists := domainMap.proxy[route.Domain]
	if !exists {
		rr := newRoundRobin(unavailableHandler())
		rrPointer = &rr
		domainMap.proxy[route.Domain] = rrPointer
	}

	proxy := httputil.NewSingleHostReverseProxy(route.URL)
	rrPointer.Add(proxy)
	report.Info("proxy.route.add", report.Data{
		"domain": route.Domain,
		"url":    route.URL.String(),
	})

	// handle signal channel
	go func() {
		for {
			_, more := <-route.C
			if !more {
				rrPointer.Remove(proxy)
				report.Info("proxy.route.remove", report.Data{
					"domain": route.Domain,
					"url":    route.URL.String(),
				})
				return
			}
		}
	}()
}
