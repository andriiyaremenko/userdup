package handlers

import (
	"fmt"
	"github.com/andriiyaremenko/userdup/accesslog"
	"net/http"
	"path"
	"strconv"
)

func Dublicates(w http.ResponseWriter, req *http.Request) {
	url := req.URL.RequestURI()
	second := path.Base(url)
	first := path.Base(url[:len(url)-len(second)])
	sUserId, err := strconv.ParseInt(second, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")
		return
	}
	fUserId, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")
		return
	}

	fmt.Fprintf(w, "%t", accesslog.CheckDublicates(fUserId, sUserId))
}
