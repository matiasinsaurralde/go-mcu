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
