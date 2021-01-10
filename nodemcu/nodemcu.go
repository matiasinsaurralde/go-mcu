package nodemcu

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

const (
	defaultBaudRate = 115200
)

var (
	errUnexpectedData = errors.New("Unexpected data")
	errACKFail        = errors.New("ACK failure")
	errNotReady       = errors.New("Device isn't ready to receive data")
)

// NodeMCU is the main data structure
type NodeMCU struct {
	cfg    *serial.Config
	port   *serial.Port
	logger *log.Logger
	ackBuf []byte

	GPIO *GPIOModule
}

// File wraps FS ops
type File struct {
	Name string
	Size int

	node *NodeMCU
}

// Remove removes a file
func (f *File) Remove() error {
	s := fmt.Sprintf("file.remove(\"%s\")\r\n", f.Name)
	f.node.logger.Printf("Run is called: %s\n", s)
	err := f.node.WriteString(s)
	if err != nil {
		return err
	}
	_, err = f.node.ReadStrings()
	return err
}

// Run is an alias for NodeMCU.Run
func (f *File) Run() error {
	return f.node.Run(f.Name)
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
// TODO: add timeout handler
func (n *NodeMCU) Sync() error {
	defer n.port.Flush()
	n.logger.Println("Sync is called")
	for {
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
			if strings.Contains(ln, "2048") {
				n.logger.Println("Sync ok")
				return nil
			}
		}
	}
}

// ListFiles returns a list of NodeMCUFile, including file size
func (n *NodeMCU) ListFiles() ([]File, error) {
	n.logger.Println("ListFiles is called")
	files := make([]File, 0)
	n.WriteString(listFilesCode)
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
		f := File{Name: name, Size: sz, node: n}
		files = append(files, f)
	}
	n.logger.Printf("Found %d files\n", len(files))
	return files, nil
}

// Run invokes an existing Lua script
// TODO: capture output
func (n *NodeMCU) Run(filename string) error {
	s := fmt.Sprintf("dofile(\"%s\")\r\n", filename)
	n.logger.Printf("Run is called: %s\n", s)
	err := n.WriteString(s)
	if err != nil {
		return err
	}
	_, err = n.ReadStrings()
	return err
}

// HardwareInfo gets HW info
func (n *NodeMCU) HardwareInfo() (*HardwareInfo, error) {
	n.logger.Println("HardwareInfo is called")
	n.WriteString(hardwareInfoCode)
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

// ReadACK checks the ACK reply
func (n *NodeMCU) ReadACK() error {
	if n.ackBuf == nil {
		n.ackBuf = make([]byte, 1)
	}
	readBytes, err := n.port.Read(n.ackBuf)
	if err != nil {
		return err
	}
	if readBytes == 0 || err != nil {
		return errACKFail
	}
	if n.ackBuf[0] != 0x06 {
		return errACKFail
	}
	n.logger.Println("ACK ok")
	return nil
}

// ReadyToRecv checks if the node is ready to receive data
func (n *NodeMCU) ReadyToRecv() bool {
	n.logger.Println("ReadyToRecv is called")
	signalBuf := make([]byte, 64)
	signalReadBytes, err := n.port.Read(signalBuf)
	if err != nil || signalReadBytes == 0 {
		return false
	}
	return bytes.ContainsRune(signalBuf, 'C')
}

// SendFile uploads a file to the device
func (n *NodeMCU) SendFile(inputFile string) error {
	n.logger.Println("SendFile is called")
	startTime := time.Now()
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	n.logger.Printf("File opened '%s'\n", inputFile)

	n.logger.Println("SendFile is called, loading recv code")
	n.WriteString(recvCode)
	n.port.Flush()
	time.Sleep(1 * time.Second)
	n.logger.Println("Calling recv()")
	n.WriteString("recv()\r\n")
	n.port.Flush()
	time.Sleep(1 * time.Second)

	if !n.ReadyToRecv() {
		return errNotReady
	}
	n.logger.Println("Device is ready to receive data")

	filename := []byte(inputFile)
	filename = append(filename, 0)
	n.logger.Println("Passing filename to recv()")
	n.port.Write(filename)
	n.port.Flush()

	err = n.ReadACK()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	readBytes := 0
	noChunks := 1
	n.logger.Println("Starting to write file data")
	var buf []byte
	for {
		buf = make([]byte, 128)
		l, err := reader.Read(buf[:cap(buf)])
		data := []byte{0x1, byte(l)}
		data = append(data, buf...)
		n.port.Write(data)
		n.port.Flush()
		err = n.ReadACK()
		if err != nil {
			panic(err)
		}
		n.logger.Printf("Sending chunk %d, initial size is %d bytes\n", noChunks, len(data))
		buf = buf[:l]
		if l == 0 {
			break
		}
		noChunks++
		readBytes += len(data)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	diff := time.Since(startTime)
	n.logger.Printf("SendFile finished, took %d milliseconds, %d chunks sent, total bytes %d\n",
		int64(diff.Milliseconds()), noChunks, readBytes)
	return nil
}

// Restart calls node.restart()
func (n *NodeMCU) Restart() error {
	defer n.port.Flush()
	n.logger.Println("Restart is called")
	err := n.WriteString("node.restart()\r\n")
	if err != nil {
		return err
	}
	return nil
}

// Compile calls node.compile()
func (n *NodeMCU) Compile(filename string) error {
	s := fmt.Sprintf("compile(\"%s\")\r\n", filename)
	n.logger.Printf("Compile is called: %s\n", s)
	err := n.WriteString(s)
	if err != nil {
		return err
	}
	_, err = n.ReadStrings()
	return err
}

// SetLogger sets the logger
func (n *NodeMCU) SetLogger(l *log.Logger) {
	n.logger = l
}

// NewNodeMCU creates a new NodeMCU object and initializes the serial connection
// Logging is disabled by default
func NewNodeMCU(port string, baudRate int) (node *NodeMCU, err error) {
	node = &NodeMCU{
		cfg:    &serial.Config{Name: port, Baud: baudRate},
		logger: log.New(ioutil.Discard, "", log.LstdFlags),
	}

	// Enable GPIO module:
	node.GPIO = &GPIOModule{node: node}
	node.port, err = serial.OpenPort(node.cfg)
	return
}
