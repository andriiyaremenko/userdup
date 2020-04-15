package accesslog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	timeFormat = "02/01/2006 15:04:05"
)

type AccessLogRepo interface {
	AllIps(ind int64, out chan<- []string)
	AllUsers() map[int64][]string
	Add(als ...AccessLog)
}

type csvRepo struct {
	sync.RWMutex
	filePath string
}

type syncIPs struct {
	sync.RWMutex
	items map[int64][]string
}

func (s *syncIPs) all() map[int64][]string {
	s.RLock()
	defer s.RUnlock()
	return s.items
}

func (s *syncIPs) allById(id int64) []string {
	s.RLock()
	defer s.RUnlock()
	return s.items[id]
}

func (s *syncIPs) addIp(id int64, ip string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.items[id]; !ok {
		s.items[id] = []string{ip}
		return
	}
	has := false
	for _, v := range s.items[id] {
		if v == ip {
			has = true
		}
	}
	if !has {
		s.items[id] = append(s.items[id], ip)
	}
}

func newSyncIPs() *syncIPs {
	return &syncIPs{items: make(map[int64][]string)}
}

func NewCsvRepo(filePath string) AccessLogRepo {
	return &csvRepo{filePath: filePath}
}

func (repo *csvRepo) Add(als ...AccessLog) {
	repo.Lock()
	defer repo.Unlock()
	f, err := os.Create(repo.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, al := range als {
		fmt.Fprintf(w, "%d,%s,%s\n", al.UserId, al.IPAddr, al.TS.Format(timeFormat))
	}
	w.Flush()
}

func (repo *csvRepo) AllUsers() map[int64][]string {
	repo.RLock()
	defer repo.RUnlock()
	file, err := os.Open(repo.filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()
	ips := newSyncIPs()
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		go func() {
			record := strings.Split(string(line), ",")
			id, err := strconv.ParseInt(record[0], 10, 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			ips.addIp(id, record[1])
		}()
	}
	return ips.all()
}

func (repo *csvRepo) AllIps(id int64, out chan<- []string) {
	repo.RLock()
	defer repo.RUnlock()
	defer close(out)
	file, err := os.Open(repo.filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	log.Printf("processing %d started", id)
	ips := newSyncIPs()
	sId := strconv.FormatInt(id, 10)
	for s.Scan() {
		line := s.Text()
		go func() {
			record := strings.Split(string(line), ",")
			if record[0] != sId {
				return
			}
			ips.addIp(id, record[1])
		}()
	}
	out <- ips.allById(id)
	log.Printf("processing %d finished", id)
}
