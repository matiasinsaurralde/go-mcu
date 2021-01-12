package gpio

type Mode string

const (
	Output Mode = "gpio.output"
	Input  Mode = "gpio.INPUT"
	Low    Mode = "gpio.LOW"
	High   Mode = "gpio.HIGH"
)
