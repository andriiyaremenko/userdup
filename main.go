package main

import (
	"flag"
	"github.com/andriiyaremenko/userdup/accesslog"
	"github.com/andriiyaremenko/userdup/cache"
	"github.com/andriiyaremenko/userdup/duplicates"
	"github.com/andriiyaremenko/userdup/handlers"
	"log"
	"net/http"
)

func main() {
	address := flag.String("addr", "localhost:8080", "host (and port) to serve app")
	filePath := flag.String("f", "./access_log.csv", "file storing access log")
	flag.Parse()
	repo := accesslog.NewCsvRepo(*filePath)
	c := cache.NewMemCahce(repo)
	c.Restore()
	defer c.Clean()
	finder := duplicates.NewDupesFinder(repo, c)
	mux := http.NewServeMux()
	mux.Handle("/", handlers.NewDuplicatesHandler(finder))
	log.Printf("serving on %s", *address)
	log.Fatal(http.ListenAndServe(*address, mux))
}
