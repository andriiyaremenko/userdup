package accesslog

import (
	"log"
)

type Duplicates struct {
	Dupes bool `json:"dupes"`
}

type DupesFinder interface {
	CheckDuplicates(fUId, sUId int64) Duplicates
}

type dupesFinder struct {
	repo AccessLogRepo
}

func NewDupesFinder(repo AccessLogRepo) DupesFinder {
	return &dupesFinder{repo: repo}
}

func (df *dupesFinder) CheckDuplicates(fUId, sUId int64) (result Duplicates) {
	result = Duplicates{Dupes: false}
	if fUId == sUId {
		result = Duplicates{Dupes: true}
		return
	}
	mCount := 2
	fAL := make(chan []string)
	sAL := make(chan []string)
	go df.repo.AllIps(fUId, fAL)
	go df.repo.AllIps(sUId, sAL)
	fIPs := <-fAL
	sIPs := <-sAL
	var r []string
	for _, v := range fIPs {
		if contains(sIPs, v) {
			r = append(r, v)
		}
	}
	if len(r) >= mCount {
		result = Duplicates{Dupes: true}
		log.Printf("%v", r)
	}
	return
}

func contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
