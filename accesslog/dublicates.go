package accesslog

type Duplicates struct {
	Dupes bool `json:"dupes"`
}

type DupesFinder interface {
	CheckDublicates(fUId, sUId int64) Duplicates
}

type dupesFinder struct {
	repo AccessLogRepo
}

func NewDupesFinder(repo AccessLogRepo) DupesFinder {
	return &dupesFinder{repo: repo}
}

func (df *dupesFinder) CheckDublicates(fUId, sUId int64) (result Duplicates) {
	result = Duplicates{Dupes: false}
	if fUId == sUId {
		return
	}
	mCount := 2
	fIPs := make(map[string]struct{})
	sIPs := make(map[string]struct{})
	done := make(chan interface{})
	defer close(done)
	fAL := make(chan AccessLog)
	sAL := make(chan AccessLog)
	go df.repo.AllById(fUId, fAL, done)
	go df.repo.AllById(sUId, sAL, done)
loop:
	for {
		select {
		case f, ok := <-fAL:
			if ok {
				fIPs[f.IPAddr] = struct{}{}
			} else {
				fAL = nil
			}
		case s, ok := <-sAL:
			if ok {
				sIPs[s.IPAddr] = struct{}{}
			} else {
				sAL = nil
			}
		}
		found := 0
		for k, _ := range fIPs {
			if _, ok := sIPs[k]; ok {
				found++
			}
			if found >= mCount {
				result = Duplicates{Dupes: true}
				break loop
			}
		}
		if sAL == nil && fAL == nil {
			break loop
		}
	}
	return
}
