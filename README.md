Host-based HTTP Request Router Configured via Etcd
==================================================


Discovery
---------

Routes are discovered from etcd from its [http API](https://coreos.com/docs/distributed-configuration/etcd-api/).

## Bootstrapping Current State

{"action":"get","node":{"key":"/hosts","dir":true,"nodes":[{"key":"/hosts/a.example.com","dir":true,"nodes":[{"key":"/hosts/a.example.com/127.0.0.1:8002","value":"127.0.0.1:8002","expiration":"2015-01-23T12:39:11.102727157Z","ttl":151,"modifiedIndex":6,"createdIndex":6}],"modifiedIndex":3,"createdIndex":3}],"modifiedIndex":3,"createdIndex":3}}

## Listening for Changes

{"action":"set","node":{"key":"/hosts/a.example.com/127.0.0.1:8003","value":"127.0.0.1:8003","expiration":"2015-01-23T12:31:23.273412205Z","ttl":10,"modifiedIndex":8,"createdIndex":8},"prevNode":{"key":"/hosts/a.example.com/127.0.0.1:8003","value":"127.0.0.1:8003","expiration":"2015-01-23T12:41:11.250241165Z","ttl":598,"modifiedIndex":7,"createdIndex":7}}

{"action":"expire","node":{"key":"/hosts/a.example.com/127.0.0.1:8003","modifiedIndex":9,"createdIndex":8},"prevNode":{"key":"/hosts/a.example.com/127.0.0.1:8003","value":"127.0.0.1:8003","expiration":"2015-01-23T12:31:23.273412205Z","modifiedIndex":8,"createdIndex":8}}

{"action":"delete","node":{"key":"/hosts/a.example.com/127.0.0.1:8001","modifiedIndex":10,"createdIndex":5},"prevNode":{"key":"/hosts/a.example.com/127.0.0.1:8001","value":"127.0.0.1:8001","expiration":"2015-01-23T12:37:21.590004293Z","ttl":329,"modifiedIndex":5,"createdIndex":5}}





