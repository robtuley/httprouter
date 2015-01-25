#!/bin/sh

if [ "$GOPATH" = "" ]; then
  printf "Empty GOPATH variable\n"
  exit
fi

# build latest

go build ./webserver/main.go
go build ./router.go

# start 3 example webservers

./main --port 8001 --label A0 &
PID1=$!

./main --port 8002 --label A1 &
PID2=$!

./main --port 8003 --label B0 &
PID3=$!

./main --port 8004 --label B1 &
PID4=$!

# add one route with 2 nodes

etcdctl set /domains/a.example.com:8080/A0 http://127.0.0.1:8001
etcdctl set /domains/a.example.com:8080/A1 http://127.0.0.1:8002

# start router

./router &
PID5=$!

# curl requests to test proxy

for (( i=1; i <= 20; i++ )); do 
 curl --resolve 'a.example.com:8080:127.0.0.1' http://a.example.com:8080/test 
done

read -p "\nPress [Enter] key to terminate..."

# kill processes

printf "wait while killing PIDs %i %i %i %i %i ...\n" $PID1 $PID2 $PID3 $PID4 $PID5

kill $PID1
kill $PID2
kill $PID3
kill $PID4
kill $PID5

sleep 5