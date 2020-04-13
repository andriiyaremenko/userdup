package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/andriiyaremenko/userdup/accesslog"
	"net/http"
	"path"
	"strconv"
)

type duplicates struct {
	finder accesslog.DupesFinder
}

func NewDuplicatesHandler(finder accesslog.DupesFinder) http.Handler {
	return &duplicates{finder: finder}
}

func (h *duplicates) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.RequestURI()
	second := path.Base(url)
	first := path.Base(url[:len(url)-len(second)])
	sUserId, err := strconv.ParseInt(second, 10, 64)
	if err != nil {
		notFound(w)
		return
	}
	fUserId, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		notFound(w)
		return
	}
	result, err := json.Marshal(h.finder.CheckDuplicates(fUserId, sUserId))
	if err != nil {
		internalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "not found")
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err)
}
