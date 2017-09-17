/* This tool puts a web server directly behind a chanlistener. In practice this
is silly but shows the ability to filter connections before they go to the
back end or decouple a web server from a listening socket. */
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/AgentZombie/chanlisten"
)

const (
	addr = ":8080"
)

func fatalIfError(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %s", msg, err.Error())
	}
}

func main() {
	cl := chanlisten.New(0)
	go func(l net.Listener) {
		s := http.Server{Handler: http.FileServer(http.Dir("/usr/share/doc"))}
		s.Serve(cl)
	}(cl)

	l, err := net.Listen("tcp", addr)
	fatalIfError(err, "listening on "+addr)

	log.Print("waiting for connections on", addr)
	for {
		conn, err := l.Accept()
		fatalIfError(err, "accepting connection")
		log.Print("accepted connection from ", conn.RemoteAddr())

		// connections could be conditionally rerouted or closed here.
		fatalIfError(cl.Queue(conn), "queuing connection")
	}
}
