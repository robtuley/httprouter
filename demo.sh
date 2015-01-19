#!/bin/sh

# build latest

go build ./webserver/main.go

# start 2 example webservers

./webserver/webserver --port 8001 --label a.example.com &
PID1=$!

./webserver/webserver --port 8002 --label b.example.com &
PID2=$!

# start proxy

#go run ./proxy.go &
#PID3=$!

# curl requests to test proxy

#curl --resolve 'a.example.com:8000:127.0.0.1' http://a.example.com:8000/test --verbose
#curl --resolve 'b.example.com:8000:127.0.0.1' http://a.example.com:8000/test --verbose

read -p "Press [Enter] key to terminate..."


# kill processes

printf "killing PIDs %i %i\n" $PID1 $PID2

kill $PID1
kill $PID2
#kill $PID3

sleep 2