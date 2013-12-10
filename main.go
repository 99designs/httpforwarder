package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

const DEFAULT_SENTRY_URL = "http://google.com.au/"

func main() {

	sentryUrl, _ := url.Parse(os.Getenv("FORWARDURL"))
	if sentryUrl.String() == "" {
		sentryUrl, _ = url.Parse(DEFAULT_SENTRY_URL)
	}
	port := os.Getenv("PORT")
	socket := os.Getenv("SOCKPATH")

	httpForwarder := NewAsyncHttpForwarder(sentryUrl)

	if socket != "" {
		// listen on socket
		log.Printf("Listening on socket %s, forwarding requests to %s", socket, sentryUrl)
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
		log.Printf("Listening on port %s, forwarding requests to %s", port, sentryUrl)
		http.Handle("/", httpForwarder)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
