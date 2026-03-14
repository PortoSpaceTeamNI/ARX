package stream

import "go.bug.st/serial"

type SerialStream struct {
	port serial.Port
}

func NewSerialStream(device string, baud int) (*SerialStream, error) {
	mode := &serial.Mode{BaudRate: baud}
	port, err := serial.Open(device, mode)
	if err != nil {
		return nil, err
	}
	port.SetDTR(true)

	return &SerialStream{
		port: port,
	}, nil
}

func (s *SerialStream) Read(p []byte) (int, error) {
	return s.port.Read(p)
}

func (s *SerialStream) Write(p []byte) (int, error) {
	return s.port.Write(p)
}

func (s *SerialStream) Close() error {
	return s.port.Close()
}
