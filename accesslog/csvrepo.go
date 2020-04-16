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

type csvRepo struct {
	sync.RWMutex
	filePath string
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

func (repo *csvRepo) AllIps(id int64) []string {
	repo.RLock()
	defer repo.RUnlock()
	file, err := os.Open(repo.filePath)
	if err != nil {
		log.Fatal(err)
		panic(err)
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
	log.Printf("processing %d finished", id)
	return ips.allById(id)
}
