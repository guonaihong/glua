package socket

import (
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newSocket,
	})

	RegisterSocketType(mod, L)
	L.Push(mod)
	return 1
}
