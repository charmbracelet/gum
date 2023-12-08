package progress

import (
	"math"
	"time"
)

type barInfo struct {
	iter        uint
	title       string
	limit       uint
	incrementTs []time.Time
}

func newBarInfo(title string, limit uint) *barInfo {
	info := &barInfo{
		title:       title,
		limit:       limit,
		incrementTs: make([]time.Time, 0, limit),
	}
	info.incrementTs = append(info.incrementTs, time.Now())
	return info
}

func (self *barInfo) Update(progressAmount uint) {
	self.iter += progressAmount

	now := time.Now()
	for i := uint(0); i < progressAmount; i++ {
		self.incrementTs = append(self.incrementTs, now)
	}
}

func (self *barInfo) Elapsed() time.Duration {
	return time.Now().Sub(self.incrementTs[0]).Truncate(time.Second)
}

func (self *barInfo) Pct() int {
	pct := math.Round(safeDivide(float64(self.iter), float64(self.limit)) * 100)
	return int(pct)
}

func (self *barInfo) Avg() time.Duration {
	if len(self.incrementTs) < 2 {
		return time.Now().Sub(self.incrementTs[0])
	}
	var sum time.Duration
	for i := 1; i < len(self.incrementTs); i++ {
		sum += self.incrementTs[i].Sub(self.incrementTs[i-1])
	}
	return (sum / time.Duration(self.iter))
}

func (self *barInfo) Eta() time.Duration {
	if self.iter >= self.limit {
		return 0
	}
	avg := self.Avg()
	if avg == 0 {
		return 0
	}

	return time.Duration(self.limit-self.iter) * avg
}
