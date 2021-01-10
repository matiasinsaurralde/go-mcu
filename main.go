package main

import (
	"fmt"

	nodemcu "github.com/matiasinsaurralde/go-mcu/nodemcu"
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
