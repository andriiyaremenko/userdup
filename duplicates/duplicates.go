package duplicates

import (
	"github.com/andriiyaremenko/userdup/accesslog"
	"github.com/andriiyaremenko/userdup/cache"
	"log"
)

type Duplicates struct {
	Dupes bool `json:"dupes"`
}

type DupesFinder interface {
	CheckDuplicates(fUId, sUId int64) Duplicates
}

type dupesFinder struct {
	cache cache.Cache
	repo  accesslog.AccessLogRepo
}

func NewDupesFinder(repo accesslog.AccessLogRepo, cache cache.Cache) DupesFinder {
	return &dupesFinder{repo: repo, cache: cache}
}

func (df *dupesFinder) CheckDuplicates(fUId, sUId int64) (result Duplicates) {
	result = Duplicates{Dupes: false}
	if fUId == sUId {
		result = Duplicates{Dupes: true}
		return
	}
	mCount := 2
	f, fChan := df.getIps(fUId)
	s, sChan := df.getIps(sUId)
	if f == nil {
		f = <-fChan
	}
	if s == nil {
		s = <-sChan
	}
	var r []string
	for _, v := range f {
		if contains(s, v) {
			r = append(r, v)
		}
	}
	if len(r) >= mCount {
		result = Duplicates{Dupes: true}
		log.Printf("%v", r)
	}
	return
}

func (df *dupesFinder) getIps(id int64) ([]string, chan []string) {
	if ips, ok := df.cache.Get(id); ok {
		return ips, nil
	}
	al := make(chan []string)
	go df.repo.AllIps(id, al)
	return nil, al
}

func contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
