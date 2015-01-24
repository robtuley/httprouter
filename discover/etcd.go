// Package discover enables discovery and tracking of traffic routing
// rules from etcd
package discover

import (
	"net/url"
	"strings"

	"github.com/robtuley/etcdwatch"
	"github.com/robtuley/report"
)

// Route summarises a desired HTTP traffic flow
type Route struct {
	Domain string
	URL    *url.URL
	C      chan int
}

func (r *Route) Close() {
	report.Info("route.close", report.Data{
		"domain": r.Domain,
		"url":    r.URL.String(),
	})
	close(r.C)
}

// Discovers routes from an etcd directory when in the form:
//
//   key:   /domains/demo.example.com/<name>
//   value: http://internal.host:8000
//
func Etcd(etcdKey string) chan Route {
	routeC := make(chan Route)
	changeC, errorC := etcdwatch.Key(etcdKey)

	go func() {
		// domainMap[domainKey{"example.com" "localhost:8000"}] == Route
		type domainKey struct {
			domain, name string
		}
		domainMap := make(map[domainKey]Route)

	nextChange:
		for {
			change := <-changeC
			report.Info("etcd.change", report.Data{
				"key":    change.Key,
				"value":  change.Value,
				"action": change.Action,
			})

			domain, name := keyToDomainAndName(etcdKey, change.Key)
			k := domainKey{domain, name}
			existingRoute, isset := domainMap[k]

			// remove route
			if isset && (change.Action == "expire" || change.Action == "delete") {
				delete(domainMap, k)
				existingRoute.Close()
			}

			if change.Action == "get" || change.Action == "set" {

				newRoute, err := createRoute(domain, change.Value)
				if err != nil {
					report.Action("route.error", report.Data{
						"error":  err.Error(),
						"change": change,
					})
					continue nextChange
				}

				// no existing route, publish new route
				if !isset {
					domainMap[k] = newRoute
					routeC <- newRoute
					report.Info("route.new", report.Data{
						"domain": newRoute.Domain,
						"url":    newRoute.URL.String(),
					})
				}

				// route already exists, but changed: close and open new
				if isset && existingRoute.URL != newRoute.URL {
					existingRoute.Close()
					domainMap[k] = newRoute
					routeC <- newRoute
					report.Info("route.change", report.Data{
						"domain":  newRoute.Domain,
						"url":     newRoute.URL.String(),
						"prevUrl": existingRoute.URL.String(),
					})
				}

			}

		}
	}()

	go func() {
		for {
			err := <-errorC
			report.Action("etcd.error", report.Data{"error": err.Error()})
		}
	}()

	return routeC
}

func keyToDomainAndName(rootDir string, key string) (string, string) {
	split := strings.Split(strings.TrimPrefix(key, rootDir+"/"), "/")
	return split[0], split[1]
}

func createRoute(domain string, urlStr string) (Route, error) {
	r := Route{Domain: domain, C: make(chan int)}
	var err error
	r.URL, err = url.Parse(urlStr)
	return r, err
}
