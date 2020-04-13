package accesslog

import "time"

type AccessLog struct {
	UserId int64
	IPAddr string
	TS     time.Time
}
