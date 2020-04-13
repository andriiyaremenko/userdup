package main

import (
	"flag"
	"github.com/andriiyaremenko/userdup/handlers"
	"log"
	"net/http"
)

func main() {
	address := flag.String("address", "localhost:8080", "host (and port) to serve app")
	flag.Parse()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handlers.Dublicates))
	log.Fatal(http.ListenAndServe(*address, mux))
}
