package command

type UpdateSerialPortCommand struct {
	LocalCommand
	SerialPort string `json:"serial_port"`
}

func (c *UpdateSerialPortCommand) ToString() string {
	return "Update Serial Port to " + c.SerialPort
}
