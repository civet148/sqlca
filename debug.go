package sqlca

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"reflect"
	"runtime"
	"strings"
)

//assert bool and string/struct/slice/map nil, call panic
func assert(v interface{}, strMsg string, args ...interface{}) {
	if isNilOrFalse(v) {
		panic(fmt.Sprintf(strMsg, args...))
	}
}

func isNilOrFalse(v interface{}) bool {
	switch v.(type) {
	case string:
		if v.(string) != "" {
			return true
		}
	case bool:
		return v.(bool)
	default:
		if !reflect.ValueOf(v).IsNil() {
			return true
		}
	}
	return false
}

func getCaller(skip int) (strFile, strFunc string, nLine int) {
	pc, f, n, ok := runtime.Caller(skip)
	if ok {
		strFile = f
		nLine = n
		strFunc = getFuncNameFromPC(pc)
	}
	return
}

// get function name from call stack
func getFuncNameFromPC(pc uintptr) (name string) {

	n := runtime.FuncForPC(pc).Name()
	ns := strings.Split(n, ".")
	name = ns[len(ns)-1]
	return
}

func fmtParentCaller(strFmt string) string {
	strFunc, strFile, nLine := getCaller(1)
	strFmt = fmt.Sprintf("<%v:%v %v()> ", strFile, nLine, strFunc) + strFmt
	return strFmt
}

func (e *Engine) setDebug(ok bool) {
	e.debug = ok
}

func (e *Engine) isDebug() bool {
	return e.debug
}

func (e *Engine) panic(strFmt string, args ...interface{}) {
	if e.isDebug() {
		panic(fmt.Sprintf(strFmt, args...))
	} else {
		strFmt = fmtParentCaller(strFmt)
		log.Errorf(strFmt, args...)
	}
}

func (e *Engine) debugf(strFmt string, args ...interface{}) {
	if e.isDebug() {
		strFmt = fmtParentCaller(strFmt)
		log.Debugf(strFmt, args...)
	}
}
