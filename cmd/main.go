package main

import (
	"finfluence-chat/internal/transport"
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()
	http.HandleFunc("/message", transport.MessageHandler)
	log.Printf("Server started at http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
