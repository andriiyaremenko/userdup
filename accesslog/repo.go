package accesslog

const (
	timeFormat = "02/01/2006 15:04:05"
)

type AccessLogRepo interface {
	AllIps(ind int64) []string
	AllUsers() map[int64][]string
	Add(als ...AccessLog)
}
