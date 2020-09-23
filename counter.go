package sqlca

import (
	"github.com/civet148/gotools/log"
	"time"
)

type counter struct {
	startTime     int64
	slowQueryTime int
}

func (e *Engine) Counter() *counter {
	return &counter{
		slowQueryTime: e.slowQueryTime,
		startTime:     time.Now().UnixNano(),
	}
}

func (c *counter) Stop(strTip string) {
	elapse := (time.Now().UnixNano() - c.startTime) / 1e6
	if elapse > int64(c.slowQueryTime) {
		log.Warnf("[%v] slow query, elapse %d ms", strTip, elapse)
	}
}
