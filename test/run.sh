#!/bin/sh

function errExit
{
  echo "$1" 1>&2
  exit 1
}

if [ "$GOPATH" = "" ]; then
  errExit "Empty GOPATH variable\n"
fi

DIR=$GOPATH/src/github.com/robtuley/httprouter
cd $DIR

PIDS=()

LOGFILE=$DIR/test/run.log

# delete previous keys 

etcdctl rm /domains/a.example.com:8080/A0 > /dev/null 2>&1
etcdctl rm /domains/a.example.com:8080/A1 > /dev/null 2>&1
etcdctl rm /domains/b.example.com:8080/B0 > /dev/null 2>&1
etcdctl rm /domains/b.example.com:8080/B1 > /dev/null 2>&1

# build latest

go build $DIR/test/webserver.go || errExit "webserver build failed\n"
go build $DIR/router.go || errExit "router build failed\n"

# start example webservers

function startWebserver
{
  ./webserver --port $1 --label $2 &
  PIDS+=($!)
}

startWebserver 8001 A0
startWebserver 8002 A1
startWebserver 8003 B0
startWebserver 8004 B1

# add one route with 2 nodes

etcdctl set /domains/a.example.com:8080/A0 http://127.0.0.1:8001 > /dev/null
etcdctl set /domains/a.example.com:8080/A1 http://127.0.0.1:8002 > /dev/null

# start router

printf "Starting router (logging to %s)" $LOGFILE
./router --log=$LOGFILE &
PIDS+=($!)
sleep 2 

# test init

printf "\nRunning tests "
typeset -i NPASS=0
typeset -i NFAIL=0
function try { 
  this="$1"; 
}

function assert {
  [ "$1" = "$2" ] && { printf "."; let NPASS+=1; return; }
  let NFAIL+=1
  printf "\nFAIL: $this\n'$1' != '$2'\n"
}

# tests

try "None routable host gets a 503 response"
 
out=$(curl -s -w "%{http_code}" -o /dev/null --resolve 'b.example.com:8080:127.0.0.1' http://b.example.com:8080/)
assert "503" "$out"

# kill processes

printf "\nKilling runnings PIDs:"
for i in "${PIDS[@]}"
do
  printf " %i" $i
  kill $i
done

# cleanup compiled bins

rm ./webserver
rm ./router

# allow PIDs to close down

printf "\nClosing down "
for (( i=1; i <= 5; i++ )); do 
  printf "."
  sleep 1 
done

# summary

printf "\n\nPASSED: %i FAILED: %i\n\n" $NPASS $NFAIL