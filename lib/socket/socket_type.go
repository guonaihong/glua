package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/yuin/gopher-lua"
	"io"
	"net"
	_ "strconv"
	"strings"
	"time"
)

type Socket struct {
	net.Conn
}

const luaSocketTypeName = "socket"

func RegisterSocketType(module *lua.LTable, L *lua.LState) {
	mt := L.NewTypeMetatable(luaSocketTypeName)
	//L.SetGlobal("socket", mt)
	//L.SetField(mt, "new", L.NewFunction(newSocket))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), socketMethods))
}

func newSocket(L *lua.LState) int {
	socket := &Socket{}
	ud := L.NewUserData()
	ud.Value = socket
	L.SetMetatable(ud, L.GetTypeMetatable(luaSocketTypeName))
	L.Push(ud)
	return 1
}

func connect(L *lua.LState) int {
	s := checkSocket(L)

	addr := L.CheckString(2)
	var err error
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

func socketClose(L *lua.LState) int {
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

func n2big_bytes(L *lua.LState) int {
	l := L.CheckNumber(2)

	head := bytes.NewBuffer(nil)
	binary.Write(head, binary.BigEndian, l)
	L.Push(lua.LString(head.Bytes()))
	return 1
}

//TODO
func addHeader(L *lua.LState) int {
	return 1
}

var socketMethods = map[string]lua.LGFunction{
	"connect":     connect,
	"read":        read,
	"write":       write,
	"close":       socketClose,
	"ntohl":       ntohl,
	"n2big_bytes": n2big_bytes,
	"add_header":  addHeader,
}

func checkSocket(L *lua.LState) *Socket {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Socket); ok {
		return v
	}

	L.ArgError(1, "socket expected")
	return nil
}
