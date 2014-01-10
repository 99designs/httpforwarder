package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	socket := ""
	tcpaddr := ":8080"
	listen := os.Getenv("LISTEN")
	n := strings.Count(listen, ":")
	if n == 0 {
		socket = listen
	} else if n == 1 {
		tcpaddr = listen
	}

	httpForwarder := NewAsyncHttpForwarder()

	log.Println("HTTP forwarder: https://github.com/99designs/httpforwarder")
	if socket != "" {
		// listen on socket
		log.Printf("Listening on socket %s", socket)
		l, err := net.Listen("unix", socket)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(http.Serve(l, httpForwarder))
	} else {
		// listen on TCP address
		log.Printf("Listening on tcp address %s", tcpaddr)
		log.Fatal(http.ListenAndServe(tcpaddr, httpForwarder))
	}
}
