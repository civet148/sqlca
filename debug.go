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

// judgement: variant is a pointer type ?
func isPtrType(v interface{}) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		return true
	}
	return false
}

// judgement: bool, string, struct, slice, map is nil or false?
func isNilOrFalse(v interface{}) bool {
	switch v.(type) {
	case string:
		if v.(string) == "" {
			return true
		}
	case bool:
		return !v.(bool)
	case int8, int16, int, int32, int64:
		{
			if fmt.Sprintf("%v", v) == "0" {
				return true
			}
		}
	default:
		{
			val := reflect.ValueOf(v)
			return val.IsNil()
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
	log.SetLevel(log.LEVEL_DEBUG)
}

func (e *Engine) isDebug() bool {
	return e.debug
}

func (e *Engine) panic(strFmt string, args ...interface{}) {
	if e.isDebug() {
		panic(fmt.Sprintf(strFmt, args...))
	} else {
		//strFmt = fmtParentCaller(strFmt)
		log.Errorf(strFmt, args...)
	}
}
