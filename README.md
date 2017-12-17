# Go Client/Server

A client/server implementation in which the client assembles raw UDP datagrams with a specified
payload, then ships them off to the server, once every two seconds.

## Usage

First and most importantly, <a target="_blank" href="https://golang.org/doc/install">install Go for your platform</a>. This project was built on v1.8.3.

The server expects a single command line argument: its port number. It should be invoked as follows:

```sh
$ go run server.go #server-listening-port-num#
```

The client expects two command line arguments: a port number to listen on, and an IP address:port number for the server. Invoke it as follows:

```sh
$ go run client.go #client-listenting-port-num# #server-ip:server-port#
```

Make sure to avoid well-known ports. Then sit back and get ready to pick your jaw up off the floor.
