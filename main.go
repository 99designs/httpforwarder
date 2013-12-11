package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	socket := os.Getenv("SOCKPATH")

	httpForwarder := NewAsyncHttpForwarder()

	if socket != "" {
		// listen on socket
		log.Printf("HTTP forwarder listening on socket %s", socket)
		l, err := net.Listen("unix", socket)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(http.Serve(l, httpForwarder))
	} else {
		if port == "" {
			port = "8080"
		}
		// listen on TCP port
		log.Printf("HTTP forwarder listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, httpForwarder))
	}
}
