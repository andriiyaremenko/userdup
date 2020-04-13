package accesslog

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type AccessLogRepo interface {
	AllById(ind int64, out chan<- AccessLog, done <-chan interface{})
}

type csvRepo struct {
	filePath string
}

func NewCsvRepo(filePath string) AccessLogRepo {
	return &csvRepo{filePath: filePath}
}

func (repo *csvRepo) AllById(id int64, out chan<- AccessLog, done <-chan interface{}) {
	defer close(out)
	file, err := os.Open(repo.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)

fileSearch:
	for {
		select {
		case <-done:
			break fileSearch
		default:
			record, err := r.Read()
			if err == io.EOF {
				break fileSearch
			}
			if err != nil {
				log.Fatal(err)
				continue fileSearch
			}
			recId, err := strconv.ParseInt(record[0], 10, 64)
			if err != nil {
				log.Fatal(err)
				continue fileSearch
			}
			if recId != id {
				continue fileSearch
			}
			ip := record[1]
			ts, err := time.Parse("02/01/2006 15:04:05", record[2])
			if err != nil {
				log.Fatal(err)
				continue fileSearch
			}
			out <- AccessLog{UserId: id, IPAddr: ip, TS: ts}
		}
	}
}
