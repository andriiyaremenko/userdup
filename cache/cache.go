package cache

type Cache interface {
	Add(uId int64, ips []string)
	Get(uId int64) ([]string, bool)
	Append(uId int64, ip string)
	Restore()
	Clean()
}
