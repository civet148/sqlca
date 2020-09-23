package sqlca

import (
	"github.com/civet148/gotools/log"
	"time"
)

type counter struct {
	slowQueryOn   bool
	startTime     int64
	slowQueryTime int
}

func (e *Engine) Counter() *counter {
	return &counter{
		slowQueryOn:   e.slowQueryOn,
		slowQueryTime: e.slowQueryTime,
		startTime:     time.Now().UnixNano(),
	}
}

func (c *counter) Stop(strTip string) {
	elapse := (time.Now().UnixNano() - c.startTime) / 1e6
	if c.slowQueryOn && c.slowQueryTime > 0 && elapse > int64(c.slowQueryTime) {
		log.Warnf("[%v] slow query, elapse %d ms", strTip, elapse)
	}
}
