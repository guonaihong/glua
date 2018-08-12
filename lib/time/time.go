package time

import (
	"github.com/yuin/gopher-lua"
	gotime "time"
)

func parseTime(t string) (rv gotime.Duration) {

	t0 := 0
	for k := 0; k < len(t); k++ {
		v := int(t[k])
		switch {
		case v >= '0' && v <= '9':
			t0 = t0*10 + (v - '0')
		case v == 's':
			rv += gotime.Duration(t0) * gotime.Second
			t0 = 0
		case v == 'm':
			if k+1 < len(t) && t[k+1] == 's' {
				rv += gotime.Duration(t0) * gotime.Millisecond
				t0 = 0
				k++
				continue
			}

			rv += gotime.Duration(t0*60) * gotime.Second
			t0 = 0
		case v == 'h':
			rv += gotime.Duration(t0*60*60) * gotime.Second
			t0 = 0
		case v == 'd':
			rv += gotime.Duration(t0*60*60*24) * gotime.Second
			t0 = 0
		case v == 'w':
			rv += gotime.Duration(t0*60*60*24*7) * gotime.Second
			t0 = 0
		case v == 'M':
			rv += gotime.Duration(t0*60*60*24*7*31) * gotime.Second
			t0 = 0
		case v == 'y':
			rv += gotime.Duration(t0*60*60*24*7*31*365) * gotime.Second
			t0 = 0
		}
	}

	return
}

func sleep(L *lua.LState) int {
	n1 := L.CheckString(1)

	n := parseTime(n1)
	gotime.Sleep(n)
	return 1
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"sleep": sleep,
	})

	L.Push(mod)
	return 1
}
