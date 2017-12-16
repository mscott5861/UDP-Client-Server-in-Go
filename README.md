# Go Client/Server

A client/server implementation in which the client assembles raw UDP datagrams with a specified
payload, then ships them off to the server, one every two seconds.

## Usage

The server expects a single command line argument: its port number. It should be invoked as follows:

```sh
$ go run server.go #server-listening-port-num#
```

The client expects two command line arguments: a port number to listen on, and an IP address:port number for the server. Invoke it as follows:

```sh
$ go run client.go #client-listenting-port-num# #server-ip:server-port#
```

Make sure to avoid well-known ports.
