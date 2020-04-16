package cache

import (
	"log"
	"sync"

	"github.com/andriiyaremenko/userdup/accesslog"
)

type memCache struct {
	sync.RWMutex
	items map[int64][]string
	repo  accesslog.AccessLogRepo
}

func (mc *memCache) Add(uId int64, ips []string) {
	mc.Lock()
	mc.items[uId] = ips
	mc.Unlock()
}

func (mc *memCache) Get(uId int64) (ips []string, ok bool) {
	mc.RLock()
	defer mc.RUnlock()
	ips, ok = mc.items[uId]
	return
}

func (mc *memCache) Append(uId int64, ip string) {
	mc.Lock()
	defer mc.Unlock()
	if ips, ok := mc.items[uId]; ok {
		mc.items[uId] = append(ips, ip)
		return
	}
	mc.items[uId] = []string{ip}
}

func (mc *memCache) Clean() {
	mc.Lock()
	mc.items = make(map[int64][]string)
	mc.Unlock()
}

func (mc *memCache) Restore() {
	mc.Lock()
	defer mc.Unlock()
	log.Print("Started cache restoration")
	mc.items = mc.repo.AllUsers()
	log.Print("Finished cache restoration")
}

func NewMemCahce(repo accesslog.AccessLogRepo) Cache {
	return &memCache{items: make(map[int64][]string), repo: repo}
}
