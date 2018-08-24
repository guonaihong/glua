package websocket

import (
	"encoding/binary"
	"fmt"
	gowebsocket "github.com/gorilla/websocket"
	"github.com/yuin/gopher-lua"
	"net/http"
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

func headersAdd(header []string, reqHeader http.Header) {

	for _, v := range header {

		headers := strings.Split(v, ":")

		if len(headers) != 2 {
			continue
		}

		headers[0] = strings.TrimSpace(headers[0])
		headers[1] = strings.TrimSpace(headers[1])

		reqHeader.Add(headers[0], headers[1])
	}

	reqHeader.Set("Accept", "*/*")

}

func connect(L *lua.LState) int {
	s := checkSocket(L)

	addr := L.CheckString(2)
	headerTb := L.ToTable(3)
	var err error
	var header []string
	reqHeader := http.Header{}

	headerTb.ForEach(func(_ lua.LValue, value lua.LValue) {
		header = append(header, value.String())
	})

	headersAdd(header, reqHeader)
	c, _, err := gowebsocket.DefaultDialer.Dial(addr, reqHeader)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	s.Conn = c

	return 0
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

	if top := L.GetTop(); top == 2 {
		arg3 := L.ToString(3)
		if arg3 != "" {
			err := s.SetReadDeadline(time.Now().Add(parseTime(arg3)))
			if err != nil {
			}
		}
	}

	mt, message, err := s.Conn.ReadMessage()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 3
	}

	if top := L.GetTop(); top == 2 {
		s.SetReadDeadline(time.Time{})
	}

	mtMessage := "text"
	if mt == gowebsocket.BinaryMessage {
		mtMessage = "binary"
	}

	L.Push(lua.LString(mtMessage))
	L.Push(lua.LString(message))
	return 2
}

func write(L *lua.LState) int {
	s := checkSocket(L)
	mtMessage := L.CheckString(2)
	data := L.CheckString(3)

	if top := L.GetTop(); top == 4 {
		arg3 := L.ToString(3)
		if arg3 != "" {
			err := s.SetWriteDeadline(time.Now().Add(parseTime(arg3)))
			if err != nil {
			}
		}
	}

	mt := gowebsocket.TextMessage
	if mtMessage == "binary" {
		mt = gowebsocket.BinaryMessage
	}

	err := s.Conn.WriteMessage(mt, []byte(data))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	if top := L.GetTop(); top == 4 {
		s.SetWriteDeadline(time.Time{})
	}

	return 0
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
