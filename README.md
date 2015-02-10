Host-based HTTP Request Router Configured via Etcd
==================================================

[![Docker Repository on Quay.io](https://quay.io/repository/robtuley/httprouter/status "Docker Repository on Quay.io")](https://quay.io/repository/robtuley/httprouter)

Configuring the Router
----------------------

The following command line arguments are available:

    -- logfile <filepath>  log to specified file
    -- logurl <url>        POST log in batches to specified URL (e.g. Loggly bulk input)
    -- etcdurl <url>       etcd URL, defaults to http://127.0.0.1:4001
    -- etcdkey <key>       etcd key to discover routes from, defaults to /domains

If no alternative log file/url is specified, the application will log to stdout. All logs are in JSON format.

Discovery
---------

Routes are discovered from etcd from its [http API](https://coreos.com/docs/distributed-configuration/etcd-api/) via the [etcdwatch package](https://github.com/robtuley/etcdwatch). To add the internal URL `http://internal.host:8000` for the `demo.example.com`:

    etcdctl set /domains/demo.example.com/myRouteName http://internal.host:8000

The application will respect any routes in keys `/domains/<domain.name>/<route.name>` which contain routable internal URLs.

Proxy
-----

HTTP requests are reverse proxied to any available domain routes via the `Host` header. If multiple routes are available, requests are round robin load balanced betwen them. If no route is available for a domain a 503 request will be served.

Example Discovery Script
------------------------

An example discovery bash script to perform a basic polling discovery of a health check URL and poppulation of the route into etcd is:

    #!/bin/sh
    if [ "$#" -ne 3 ]; then
      echo "Usage: $0 <domain> <ip> <port>" >&2
      exit 1
    fi
    
    domain=$1
    ip=$2
    port=$3
    printf "domain:> %s ip:> %s port:> %s" $domain $ip $port
    
    while true; do
      curl -s --max-time 4 --connect-timeout 1 -A discovery-health-check -f http://$ip:$port/health-check > /dev/null
      if [ $? -eq 0 ]; then
        etcdctl set /domains/$domain/$ip:$port http://$ip:$port --ttl 10
	    printf "ok:> %s at %s" $domain $ip:$port
      else
        etcdctl rm /domains/$domain/$ip:$port
	    printf "error:> %s at %s" $domain $ip:$port
      fi
      sleep 5
    done 
