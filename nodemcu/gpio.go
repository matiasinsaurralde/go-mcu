package nodemcu

import (
	"fmt"
	"github.com/matiasinsaurralde/go-mcu/nodemcu/gpio"
)

// GPIOModule wraps the gpio module
type GPIOModule struct {
	node *NodeMCU
}

// Mode calls gpio.mode, constants are defined in gpio/constants.go
func (m *GPIOModule) Mode(pin int, mode gpio.Mode) error {
	s := fmt.Sprintf("gpio.mode(%d, %s)\r\n", pin, string(mode))
	m.node.logger.Printf("gpio.mode is called: %s\n", s)
	err := m.node.WriteString(s)
	if err != nil {
		return err
	}
	_, err = m.node.ReadStrings()
	return err
}
