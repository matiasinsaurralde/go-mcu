go-mcu
==

[![GoDoc](https://godoc.org/github.com/urfave/cli?status.svg)](https://godoc.org/github.com/matiasinsaurralde/go-mcu/nodemcu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


`go-mcu` provides an alternative way to work with NodeMCU-based modules like the [ESP8266](https://www.espressif.com/en/products/socs/esp8266). Inspired by [NodeMCU-Tool](https://github.com/andidittrich/NodeMCU-Tool) and [nodemcu-uploader](https://github.com/kmpm/nodemcu-uploader). It can be used as a Go package but also as a standalone CLI tool.

One of the goals is to take advantage of [Go's cross compilation capability](https://dave.cheney.net/tag/cross-compilation) while keeping a minimal set of runtime dependencies.

## Getting started

To download the package use:

```
go get -u github.com/matiasinsaurralde/go-mcu
```

The CLI tool should be available afterwards:

```
$ go-mcu
NAME:
   go-mcu - NodeMCU tool (in Golang)

USAGE:
   go-mcu [global options] command [command options] [arguments...]

COMMANDS:
   hwinfo   retrieves hardware info
   upload   upload a file
   run      invoke a script
   restart  trigger node restart
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value  Serial port device
   --baud value  Baud rate (default: 115200)
   --help, -h    show help (default: false)

```

Binaries are available in the [releases page](https://github.com/matiasinsaurralde/go-mcu/releases).

## Supported/tested platforms

- Mac
- Windows
- Linux

## Sample code

Full sample [here](https://github.com/matiasinsaurralde/go-mcu/blob/master/samples/sample.go).

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

## Additional features

### GPIO module

```go
node.GPIO.Mode(4, nodemcu.GPIO_OUTPUT)
time.Sleep(1 * time.Second)
node.GPIO.Mode(4, nodemcu.GPIO_HIGH)
time.Sleep(1 * time.Second)
node.GPIO.Mode(4, nodemcu.GPIO_LOW)
```

## License

[MIT](README.md)
