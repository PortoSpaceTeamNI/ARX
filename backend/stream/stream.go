package stream

import (
	"io"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"sync"
	"time"

	"go.bug.st/serial"
)

type Stream interface {
	io.ReadWriteCloser
}

type StreamManager struct {
	mu            sync.RWMutex
	activePort    serial.Port
	currentConfig *serial.Mode
}

func NewStreamManager() *StreamManager {
	return &StreamManager{}
}

func (m *StreamManager) OpenPort(portName string, baud int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activePort != nil {
		m.activePort.Close()
	}

	mode := &serial.Mode{
		BaudRate: baud,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}

	m.activePort = port
	log.Printf("Switched to serial port: %s at %d baud", portName, baud)
	return nil
}

func (m *StreamManager) Run() {
	go func() {
		buf := make([]byte, 256)
		for {
			m.mu.RLock()
			port := m.activePort
			m.mu.RUnlock()

			if port == nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			n, err := port.Read(buf)
			if n > 0 {
				data := make([]byte, n)
				copy(data, buf[:n])
				globals.ByteChannel <- data
			}

			if err != nil {
				log.Printf("Read error: %v. Waiting for port recovery...", err)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for req := range globals.RequestChannel {
			if req.IsRemote() {
				m.mu.RLock()
				if m.activePort != nil {
					pkt := req.ToPacket()
					bytes := pkt.ToBytes()

					_, err := m.activePort.Write(bytes)
					if err != nil {
						log.Printf("Write error: %v", err)
					}
				}
				m.mu.RUnlock()

			} else {
				if portCmd, ok := req.(*command.UpdateSerialPortCommand); ok {
					newPort := portCmd.SerialPort

					err := m.OpenPort(newPort, 115200)
					if err != nil {
						log.Printf("Failed to switch to %s: %v", newPort, err)
					}
				}
			}
		}
	}()
}
