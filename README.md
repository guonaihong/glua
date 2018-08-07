# glua
存放导出给lua用的库

#### usage
```go
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
```

#### example

uuid.lua
```lua
local uuid = require("uuid")
print(uuid:newv4())
```

```bash
env GOPATH=`pwd` go build github.com/guonaihong/glua/cmd/glua
glua uuid.lua
```
