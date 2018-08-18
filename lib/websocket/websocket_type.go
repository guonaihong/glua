package websocket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	gowebsocket "github.com/gorilla/websocket"
	"github.com/yuin/gopher-lua"
	"io"
	"net"
	_ "strconv"
	"strings"
	"time"
)

type WebSocket struct {
	*gowebsocket.Conn
}

const luaSocketTypeName = "websocket"

func RegisterSocketType(module *lua.LTable, L *lua.LState) {
	mt := L.NewTypeMetatable(luaSocketTypeName)
	//L.SetGlobal("websocket", mt)
	//L.SetField(mt, "new", L.NewFunction(newSocket))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), websocketMethods))
}

func newWebSocket(L *lua.LState) int {
	websocket := &WebSocket{}
	ud := L.NewUserData()
	ud.Value = websocket
	L.SetMetatable(ud, L.GetTypeMetatable(luaSocketTypeName))
	L.Push(ud)
	return 1
}

func connect(L *lua.LState) int {
	s := checkSocket(L)

	addr := L.CheckString(2)
	var err error

	c, _, err := websocket.DefaultDialer.Dial(u1.String(), header)
	if err != nil {
		return nil, err
	}

	s.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		L.ArgError(1, err.Error())
		return 0
	}

	return 1
}

func parseTime(s string) time.Duration {
	s = strings.ToLower(s)

	rv := int64(0)

	fmt.Sscanf(s, "%d", &rv)
	switch {
	case strings.HasSuffix(s, "ms"):
		rv = rv * int64(time.Millisecond)
	case strings.HasSuffix(s, "s"):
		rv = rv * int64(time.Second)
	}
	return time.Duration(rv)
}

func read(L *lua.LState) int {
	s := checkSocket(L)

	n := L.ToInt(2)

	if top := L.GetTop(); top == 3 {
		arg3 := L.ToString(3)
		if arg3 != "" {
			err := s.SetReadDeadline(time.Now().Add(parseTime(arg3)))
			if err != nil {
			}
		}
	}

	body := make([]byte, n)
	n1, err := io.ReadFull(s.Conn, body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if top := L.GetTop(); top == 3 {
		s.SetReadDeadline(time.Time{})
	}
	L.Push(lua.LString(body[:n1]))
	return 1
}

func write(L *lua.LState) int {
	s := checkSocket(L)
	data := L.CheckString(2)

	if top := L.GetTop(); top == 3 {
		arg3 := L.ToString(3)
		if arg3 != "" {
			err := s.SetWriteDeadline(time.Now().Add(parseTime(arg3)))
			if err != nil {
			}
		}
	}
	n, err := s.Conn.Write([]byte(data))
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if top := L.GetTop(); top == 3 {
		s.SetWriteDeadline(time.Time{})
	}

	L.Push(lua.LNumber(n))
	return 1
}

func websocketClose(L *lua.LState) int {
	s := checkSocket(L)
	s.Conn.Close()
	return 1
}

func ntohl(L *lua.LState) int {
	len := L.CheckString(2)
	u32 := binary.BigEndian.Uint32([]byte(len))
	L.Push(lua.LNumber(u32))
	return 1
}

var websocketMethods = map[string]lua.LGFunction{
	"connect": connect,
	"read":    read,
	"write":   write,
	"close":   websocketClose,
}

func checkSocket(L *lua.LState) *WebSocket {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*WebSocket); ok {
		return v
	}

	L.ArgError(1, "websocket expected")
	return nil
}
