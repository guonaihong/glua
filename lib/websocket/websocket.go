package websocket

import (
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newWebSocket,
	})

	RegisterSocketType(mod, L)
	L.Push(mod)
	return 1
}
