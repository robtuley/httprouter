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
./router --logfile=$LOGFILE &
PIDS+=($!)
sleep 2 

# test init

printf "\nRunning tests "

typeset -i NPASS=0
typeset -i NFAIL=0

function try {
  TEST_OUTPUT="" 
  TEST_NAME="$1"; 
}

function assertOutputIs {
  [ "$1" = "$TEST_OUTPUT" ] && { printf "."; let NPASS+=1; return; }
  let NFAIL+=1
  printf "\nFAIL: $TEST_NAME\n'$1' != '$TEST_OUTPUT'\n"
}

function assertOutputContains {
  [[ "$TEST_OUTPUT" = *"$1"* ]] && { printf "."; let NPASS+=1; return; }
  let NFAIL+=1
  printf "\nFAIL: $TEST_NAME\n'$TEST_OUTPUT' does not contain '$1'\n"
}

function httpStatusCodeForDomain {
  TEST_OUTPUT=$TEST_OUTPUT$(curl -s -w "%{http_code}" -o /dev/null --resolve "$1:8080:127.0.0.1" http://$1:8080/)
}

function makeRequestToDomain {
  TEST_OUTPUT=$TEST_OUTPUT$(curl -s --resolve "$1:8080:127.0.0.1" http://$1:8080/)
}

function repeat {
  number=$1
  shift
  for i in `seq $number`; do
    $@
  done
}

# tests

try "Non routable host gets a 503 response"
 
httpStatusCodeForDomain "b.example.com"
assertOutputIs "503"

try "Router host with 2 backends round robins between hosts"

repeat 5 makeRequestToDomain "a.example.com"
assertOutputContains "A0A1A0A1"

try "When etcd key deleted, backend removed from round robin pool"

etcdctl rm /domains/a.example.com:8080/A0 > /dev/null
sleep 1

repeat 4 makeRequestToDomain "a.example.com"
assertOutputIs "A1A1A1A1"

try "Etcd key update changes backend route"

etcdctl set /domains/a.example.com:8080/A1 http://127.0.0.1:8001 --ttl 5 > /dev/null
sleep 1

repeat 4 makeRequestToDomain "a.example.com"
assertOutputIs "A0A0A0A0"

try "Etcd key expiry removed route"

sleep 5
httpStatusCodeForDomain "a.example.com"
assertOutputIs "503"

for i in `seq 3`; do

  try "Rapid succession of etcd key changes renewing TTL leases #$i"

  etcdctl set /domains/a.example.com:8080/A0 http://127.0.0.1:8001 --ttl 10 > /dev/null
  etcdctl set /domains/a.example.com:8080/A1 http://127.0.0.1:8002 --ttl 10 > /dev/null
  etcdctl set /domains/b.example.com:8080/B0 http://127.0.0.1:8003 --ttl 10 > /dev/null
  etcdctl set /domains/b.example.com:8080/B1 http://127.0.0.1:8004 --ttl 10 > /dev/null
  sleep 4

  repeat 5 makeRequestToDomain "a.example.com"
  repeat 5 makeRequestToDomain "b.example.com"

  assertOutputContains "A0A1A0A1"
  assertOutputContains "B0B1B0B1"

done

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