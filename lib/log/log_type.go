package log

import (
	"bytes"
	"fmt"
	"github.com/yuin/gopher-lua"
	"strings"
	"time"
)

const (
	DEBUG = iota + 1
	INFO
	WARN
	ERROR
)

const luaLogTypeName = "log"

type Log struct {
	level int
}

func levelStr2N(str string) int {
	rv := -1
	switch {
	case strings.EqualFold(str, "DEBUG"):
		rv = DEBUG
	case strings.EqualFold(str, "INFO"):
		rv = INFO
	case strings.EqualFold(str, "WARN"):
		rv = WARN
	case strings.EqualFold(str, "ERROR"):
		rv = ERROR
	}
	return rv
}

func levelN2Str(level int) string {
	rv := "unkown"
	switch {
	case level == DEBUG:
		rv = "DEBUG"
	case level == INFO:
		rv = "INFO"
	case level == WARN:
		rv = "WARN"
	case level == ERROR:
		rv = "ERROR"
	}
	return rv
}

func (l *Log) Log(level int, msg string) {

	curLevel := level
	if curLevel < l.level {
		return
	}

	now := time.Now()
	head := fmt.Sprintf("[%s] [%s] ",
		now.Format("2006-01-02 15:04:05.000"),
		levelN2Str(level))

	fmt.Print(head, msg)
}

func (l *Log) debug(msg string) {
	l.Log(DEBUG, msg)
}

func (l *Log) info(msg string) {
	l.Log(INFO, msg)
}

func (l *Log) warn(msg string) {
	l.Log(WARN, msg)
}

func (l *Log) error(msg string) {
	l.Log(ERROR, msg)
}

func RegisterLogType(module *lua.LTable, L *lua.LState) {
	mt := L.NewTypeMetatable(luaLogTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), luaLogMethods))
}

func newLog(L *lua.LState) int {
	log := &Log{
		level: levelStr2N(L.CheckString(1)),
	}

	ud := L.NewUserData()
	ud.Value = log
	L.SetMetatable(ud, L.GetTypeMetatable(luaLogTypeName))
	L.Push(ud)
	return 1
}

func getAll(L *lua.LState) string {
	var out bytes.Buffer

	top := L.GetTop()
	for i := 2; i <= top; i++ {
		out.WriteString(L.ToStringMeta(L.Get(i)).String())
		if i != top {
			out.WriteString("\t")
		}
	}

	return out.String()
}

func debug(L *lua.LState) int {
	log := checkLog(L)
	msg := getAll(L)
	log.debug(msg)
	return 1
}

func info(L *lua.LState) int {
	log := checkLog(L)
	msg := getAll(L)
	log.info(msg)
	return 1
}

func warn(L *lua.LState) int {
	log := checkLog(L)
	msg := getAll(L)
	log.warn(msg)
	return 1
}

func error(L *lua.LState) int {
	log := checkLog(L)
	msg := getAll(L)
	log.error(msg)
	return 1
}

var luaLogMethods = map[string]lua.LGFunction{
	"debug": debug,
	"info":  info,
	"warn":  warn,
	"error": error,
}

func checkLog(L *lua.LState) *Log {
	ud := L.CheckUserData(1)

	if v, ok := ud.Value.(*Log); ok {
		return v
	}

	L.ArgError(1, "log expected")
	return nil
}
