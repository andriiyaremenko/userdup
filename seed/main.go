package main

import (
	"flag"
	"fmt"
	"github.com/andriiyaremenko/userdup/accesslog"
	"log"
	"math/rand"
	"time"
)

func main() {
	filePath := flag.String("f", "./access_log.csv", "file storing access log")
	log.Print("started seeding...")
	log.Printf("Seeding file %s...", *filePath)
	repo := accesslog.NewCsvRepo(*filePath)
	var als []accesslog.AccessLog
	var counter int64
	for i := 1; i <= 10000; i++ {
		n := rand.Int63n(2500000 / int64(i))
		ips := [6]string{
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
			fmt.Sprintf("127.0.0.%d", rand.Int31n(10)),
		}
		var j int64 = 0
		for ; j < n; j++ {
			var year int64 = 60 * 60 * 24 * 365
			randomTime := rand.Int63n(time.Now().Unix()-year) + year
			randomNow := time.Unix(randomTime, 0)
			randomIp := ips[rand.Int31n(6)]
			als = append(als, accesslog.AccessLog{UserId: int64(i), IPAddr: randomIp, TS: randomNow})
			counter++
		}
	}
	log.Printf("Saving %d records...", counter)
	repo.Add(als...)
	log.Printf("Saved %d records", counter)
	log.Print("Seeding finished")
}
