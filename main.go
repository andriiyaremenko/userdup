package main

import (
	"flag"
	"github.com/andriiyaremenko/userdup/accesslog"
	"github.com/andriiyaremenko/userdup/handlers"
	"log"
	"net/http"
)

func main() {
	address := flag.String("address", "localhost:8080", "host (and port) to serve app")
	filePath := flag.String("f", "./access_log.csv", "file storing access log")
	flag.Parse()
	repo := accesslog.NewCsvRepo(*filePath)
	finder := accesslog.NewDupesFinder(repo)
	mux := http.NewServeMux()
	mux.Handle("/", handlers.NewDublicatesHandler(finder))
	log.Printf("serving on %s", *address)
	log.Fatal(http.ListenAndServe(*address, mux))
}
