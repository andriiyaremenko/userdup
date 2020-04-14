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
	Add(als ...AccessLog)
}

type csvRepo struct {
	sync.RWMutex
	filePath string
}

type syncIPs struct {
	sync.RWMutex
	ips []string
}

func (s *syncIPs) all() []string {
	s.RLock()
	defer s.RUnlock()
	return s.ips
}

func (s *syncIPs) add(ip string) {
	s.Lock()
	s.ips = append(s.ips, ip)
	s.Unlock()
}

func NewCsvRepo(filePath string) AccessLogRepo {
	return &csvRepo{filePath: filePath}
}

func (repo *csvRepo) Add(als ...AccessLog) {
	repo.Lock()
	defer repo.Unlock()
	f, err := os.Create(repo.filePath)
	w := bufio.NewWriter(f)
	defer f.Close()
	for _, al := range als {
		fmt.Fprintf(w, "%d,%s,%s\n", al.UserId, al.IPAddr, al.TS.Format(timeFormat))
	}
	if err != nil {
		log.Fatal(err)
	}
	w.Flush()
}

func (repo *csvRepo) AllIps(id int64, out chan<- []string) {
	repo.RLock()
	defer repo.RUnlock()
	defer close(out)
	file, err := os.Open(repo.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	log.Printf("processing %d started", id)
	ips := new(syncIPs)
	sId := strconv.FormatInt(id, 10)
	for s.Scan() {
		line := s.Text()
		go func() {
			record := strings.Split(string(line), ",")
			if record[0] != sId {
				return
			}
			ip := record[1]
			has := false
			for _, v := range ips.all() {
				if v == ip {
					has = true
				}
			}
			if !has {
				ips.add(ip)
			}
		}()
	}
	out <- ips.all()
	log.Printf("processing %d finished", id)
}
