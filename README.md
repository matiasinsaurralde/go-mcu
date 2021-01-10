go-mcu
==

`go-mcu` provides an alternative way to work with NodeMCU-based modules like the [ESP8266](https://www.espressif.com/en/products/socs/esp8266). Inspired by [NodeMCU-Tool](https://github.com/andidittrich/NodeMCU-Tool) and [nodemcu-uploader](https://github.com/kmpm/nodemcu-uploader). It can be used as a Go package but also as a standalone CLI tool.

One of the goals is to take advantage of [Go's cross compilation capability](https://dave.cheney.net/tag/cross-compilation) while keeping a minimal set of runtime dependencies.

```go
import(
	"fmt"

	nodemcu "github.com/matiasinsaurralde/go-mcu"
)
func main() {
	node, err := nodemcu.NewNodeMCU("/dev/cu.usbserial-1410", 115200)
	if err != nil {
		panic(err)
	}

	// Initial interaction:
	err = node.Sync()
	if err != nil {
		panic(err)
	}

	// List files:
	files, err := node.ListFiles()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d files\n", len(files))
	for i, file := range files {
		fmt.Printf("File #%d, '%s' (%d bytes)\n", i, file.Name, file.Size)
	}

	// Invoke file:
	err = node.Run("blink.lua")
	if err != nil {
		panic(err)
	}

	// Hardware info:
	hwInfo, err := node.HardwareInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(hwInfo)
}

```
