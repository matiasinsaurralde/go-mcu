go-mcu
==

`go-mcu` provides an alternative way to work with NodeMCU-based modules like the [ESP8266](https://www.espressif.com/en/products/socs/esp8266). Inspired by [NodeMCU-Tool](https://github.com/andidittrich/NodeMCU-Tool) and [nodemcu-uploader](https://github.com/kmpm/nodemcu-uploader). It can be used as a Go package but also as a standalone CLI tool.

One of the goals is to take advantage of [Go's cross compilation capability](https://dave.cheney.net/tag/cross-compilation) while keeping a minimal set of runtime dependencies.

```go
package main

import (
	"fmt"
	"log"
	"os"

	nodemcu "github.com/matiasinsaurralde/go-mcu/nodemcu"
)

const (
	logPrefix = "go-mcu "
)

func main() {
	// Setup the port
	node, err := nodemcu.NewNodeMCU("/dev/cu.usbserial-1410", 115200)
	if err != nil {
		panic(err)
	}

	// Set a custom logger
	l := log.New(os.Stdout, logPrefix, log.LstdFlags)
	node.SetLogger(l)

	// Initialization
	err = node.Sync()
	if err != nil {
		panic(err)
	}

	// Upload a file
	err = node.SendFile("test.lua")
	if err != nil {
		panic(err)
	}

	// List files
	files, err := node.ListFiles()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d files\n", len(files))
	for i, file := range files {
		fmt.Printf("File #%d, '%s' (%d bytes)\n", i, file.Name, file.Size)
	}

	// Invoke a file
	err = node.Run("blink.lua")
	if err != nil {
		panic(err)
	}

	// Retrieve hardware info
	hwInfo, err := node.HardwareInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(hwInfo)
}
```
