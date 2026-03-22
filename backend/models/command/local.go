package command

type UpdateSerialPortCommand struct {
	LocalCommand
	SerialPort string `json:"serial_port"`
}

func (c *UpdateSerialPortCommand) ToString() string {
	return "Update Serial Port to " + c.SerialPort
}

type LocalLogCommand struct {
	LocalCommand
	Command LogAction `json:"command"`
}

func (c *LocalLogCommand) ToString() string {
	if c.Command == LogStart {
		return "Start Local Log"
	}

	return "Stop Local Log"
}
