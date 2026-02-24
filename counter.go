package sqlca

import (
	"bytes"
	"fmt"
	"github.com/civet148/log"
	"time"
)

const (
	printMaxSqlCount = 4096
	printSqlSuffix   = "..."
)

func (e *Engine) SqlCounter() func(msg string, args ...any) {
	startTime := time.Now().UnixNano()
	return func(msg string, args ...any) {
		elapse := (time.Now().UnixNano() - startTime) / 1e6
		buf := bytes.NewBufferString(msg)
		if bytes.Count(buf.Bytes(), nil) > printMaxSqlCount {
			buf.Truncate(printMaxSqlCount)
			msg = buf.String() + printSqlSuffix
		}

		if e.slowQueryOn {
			if e.slowQueryTime == 0 {
				log.Debugf("query elapse %dms: %s", elapse, fmt.Sprintf(msg, args...))
			} else if elapse > int64(e.slowQueryTime) {
				log.Warnf("slow query elapse %dms: %s", elapse, fmt.Sprintf(msg, args...))
			}
		}
	}
}
