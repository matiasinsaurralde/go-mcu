package nodemcu

import "fmt"

type GPIOMode string

const (
	GPIO_OUTPUT = "gpio.OUTPUT"
	GPIO_INPUT  = "gpio.INPUT"
	GPIO_LOW    = "gpio.LOW"
	GPIO_HIGH   = "gpio.HIGH"
)

// GPIOModule wraps the gpio module
type GPIOModule struct {
	node *NodeMCU
}

// Mode calls gpio.mode
func (m *GPIOModule) Mode(pin int, mode GPIOMode) error {
	s := fmt.Sprintf("gpio.mode(%d, %s)\r\n", pin, string(mode))
	m.node.logger.Printf("gpio.mode is called: %s\n", s)
	err := m.node.WriteString(s)
	if err != nil {
		return err
	}
	_, err = m.node.ReadStrings()
	return err
}
