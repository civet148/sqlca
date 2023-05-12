package sqlca

import (
	"bytes"
	"github.com/civet148/log"
	"time"
)

const (
	printMaxSqlCount = 4096
	printSqlSuffix   = "..."
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

	buf := bytes.NewBufferString(strTip)
	if bytes.Count(buf.Bytes(), nil) > printMaxSqlCount {
		buf.Truncate(printMaxSqlCount)
		strTip = buf.String() + printSqlSuffix
	}

	if c.slowQueryOn {
		if c.slowQueryTime == 0 {
			log.Debugf("query elapse %d ms %s", elapse, strTip)
		} else if elapse > int64(c.slowQueryTime) {
			log.Warnf("slow query elapse %d ms %s", elapse, strTip)
		}
	}
}
