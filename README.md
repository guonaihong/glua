# glua
存放导出给lua用的库

#### usage
cmd/glua/glua.go
```go

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
```

#### `build glua`
```
env GOPATH=`pwd` go build github.com/guonaihong/glua/cmd/glua
```

#### example
uuid.lua
```lua
local uuid = require("uuid")
print(uuid:newv4())
```

```bash
glua uuid.lua
```

time.lua
```lua
local time = require("time")
time.sleep("1m1s500ms")
```

```bash
time glua ./time.lua

real    1m1.506s
user    0m0.000s
sys     0m0.004s

```
