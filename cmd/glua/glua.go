package main

import (
	"fmt"
	"github.com/guonaihong/glua/lib/cmdparse"
	"github.com/guonaihong/glua/lib/json"
	"github.com/guonaihong/glua/lib/log"
	"github.com/guonaihong/glua/lib/socket"
	"github.com/guonaihong/glua/lib/strings"
	"github.com/guonaihong/glua/lib/time"
	"github.com/guonaihong/glua/lib/uuid"
	"github.com/guonaihong/glua/lib/websocket"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
)

func main() {
	L := lua.NewState()
	L.PreloadModule("socket", socket.Loader)
	L.PreloadModule("cmd", cmdparse.Loader)
	L.PreloadModule("time", time.Loader)
	L.PreloadModule("strings", strings.Loader)
	L.PreloadModule("log", log.Loader)
	L.PreloadModule("uuid", uuid.Loader)
	L.PreloadModule("json", json.Loader)
	L.PreloadModule("websocket", websocket.Loader)

	for _, v := range os.Args[1:] {
		all, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		err = L.DoString(string(all))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}
}
