package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

const (
	defaultBaudRate = 115200
)

var (
	errUnexpectedData = errors.New("Unexpected data")
)

// NodeMCU is the main data structure
type NodeMCU struct {
	cfg  *serial.Config
	port *serial.Port
}

// File wraps FS ops.
// TODO: implement read/write ops.
type File struct {
	Name string
	Size int
}

// HardwareInfo contains hardware info.
type HardwareInfo struct {
	ChipID     int
	FlashSize  int
	FlashMode  int
	FlashSpeed int
	FlashID    int
}

// WriteString writes stuff
func (n *NodeMCU) WriteString(input string) error {
	_, err := n.port.Write([]byte(input))
	if err != nil {
		return err
	}
	return nil
}

// ReadStrings reads multiline output
func (n *NodeMCU) ReadStrings() ([]string, error) {
	defer n.port.Flush()
	reader := bufio.NewReader(n.port)
	text, err := reader.ReadString('>')
	if err != nil {
		return nil, err
	}
	splits := strings.Split(text, "\r\n")
	return splits, nil
}

func (n *NodeMCU) fixStrings(splits []string) []string {
	return splits[1 : len(splits)-1]
}

// TODO: handle errors
func (n *NodeMCU) parseTab(input []string, intValue bool) (map[string]interface{}, error) {
	tab := make(map[string]interface{}, 0)
	for _, ln := range input {
		if !strings.Contains(ln, "|") {
			continue
		}
		splits := strings.Split(ln, "|")
		key := strings.TrimSpace(splits[0])
		val := strings.TrimSpace(splits[1])
		if !intValue {
			tab[key] = val
			continue
		}
		i, _ := strconv.Atoi(val)
		tab[key] = i
	}
	return tab, nil
}

// Sync runs test code
func (n *NodeMCU) Sync() error {
	// defer n.port.Flush()
	var ready bool
	for !ready {
		if err := n.WriteString("print(1024*2);\r\n"); err != nil {
			return err
		}
		output, err := n.ReadStrings()
		if err != nil {
			return err
		}
		if len(output) == 0 {
			return errUnexpectedData
		}
		for _, ln := range output {
			ln = strings.TrimSpace(ln)
			i, _ := strconv.Atoi(ln)
			if i == 2048 {
				ready = true
			}
		}
	}
	return nil
}

// ListFiles returns a list of NodeMCUFile, including file size
func (n *NodeMCU) ListFiles() ([]File, error) {
	files := make([]File, 0)
	n.WriteString("for key,value in pairs(file.list()) do print(key,\"|\",value) end\r\n")
	output, err := n.ReadStrings()
	if err != nil {
		return nil, err
	}
	output = n.fixStrings(output)
	for _, v := range output {
		splits := strings.Split(v, "|")
		sz, err := strconv.Atoi(strings.TrimSpace(splits[1]))
		if err != nil {
			continue
		}
		name := strings.TrimSpace(splits[0])
		f := File{Name: name, Size: sz}
		files = append(files, f)
	}
	return files, nil
}

// Run invokes an existing Lua script
// TODO: capture output
func (n *NodeMCU) Run(filename string) error {
	s := fmt.Sprintf("dofile(\"%s\")\r\n", filename)
	return n.WriteString(s)
}

// HardwareInfo gets HW info
func (n *NodeMCU) HardwareInfo() (*HardwareInfo, error) {
	n.WriteString("for key,value in pairs(node.info('hw')) do k=tostring(key) print(k, '|', tostring(value)) end\r\n")
	output, err := n.ReadStrings()
	if err != nil {
		return nil, err
	}
	output = n.fixStrings(output)
	m, err := n.parseTab(output, true)
	if err != nil {
		return nil, err
	}
	hwInfo := &HardwareInfo{}
	for k, v := range m {
		switch k {
		case "chip_id":
			hwInfo.ChipID = v.(int)
		case "flash_size":
			hwInfo.FlashSize = v.(int)
		case "flash_mode":
			hwInfo.FlashMode = v.(int)
		case "flash_speed":
			hwInfo.FlashSpeed = v.(int)
		case "flash_id":
			hwInfo.FlashID = v.(int)
		}
	}
	return hwInfo, nil
}

// NewNodeMCU creates a new NodeMCU object and initializes the serial connection
func NewNodeMCU(port string, baudRate int) (node *NodeMCU, err error) {
	node = &NodeMCU{
		cfg: &serial.Config{Name: port, Baud: baudRate},
	}
	node.port, err = serial.OpenPort(node.cfg)
	return
}

func main() {
	node, err := NewNodeMCU("/dev/cu.usbserial-1410", 115200)
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
	hwInfoJSON, _ := json.Marshal(&hwInfo)
	fmt.Println("hwinfo=", string(hwInfoJSON))
}
