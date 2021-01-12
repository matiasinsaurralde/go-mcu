package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	nodemcu "github.com/matiasinsaurralde/go-mcu/nodemcu"
	"github.com/urfave/cli/v2"
)

const (
	logPrefix       = "go-mcu "
	defaultBaudRate = 115200
)

var (
	portDevice string
	baudRate   int
)

func initNode() (*nodemcu.NodeMCU, error) {
	node, err := nodemcu.NewNodeMCU(portDevice, baudRate)
	if err != nil {
		return nil, err
	}
	l := log.New(os.Stdout, logPrefix, log.LstdFlags)
	node.SetLogger(l)
	err = node.Sync()
	if err != nil {
		return nil, err
	}
	return node, err
}

func main() {
	app := &cli.App{
		Name:  "go-mcu",
		Usage: "NodeMCU tool (in Golang)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Usage:       "Serial port device",
				Destination: &portDevice,
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "baud",
				Usage:       "Baud rate",
				Destination: &baudRate,
				Required:    true,
				Value:       defaultBaudRate,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "hwinfo",
				Usage: "retrieves hardware info",
				Action: func(c *cli.Context) error {
					node, err := initNode()
					if err != nil {
						return err
					}
					hwInfo, err := node.HardwareInfo()
					if err != nil {
						return err
					}
					hwInfoJSON, err := json.Marshal(hwInfo)
					if err != nil {
						return err
					}
					fmt.Println(string(hwInfoJSON))
					return nil
				},
			},
			{
				Name:  "upload",
				Usage: "upload a file",
				Action: func(c *cli.Context) error {
					node, err := initNode()
					if err != nil {
						return err
					}
					filename := c.Args().First()
					err = node.SendFile(filename)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "invoke a script",
				Action: func(c *cli.Context) error {
					node, err := initNode()
					if err != nil {
						return err
					}
					filename := c.Args().First()
					err = node.Run(filename)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "restart",
				Usage: "trigger node restart",
				Action: func(c *cli.Context) error {
					node, err := initNode()
					if err != nil {
						return err
					}
					err = node.Restart()
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
