package stream

import (
	"fmt"
	"io"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"os"
	"path/filepath"
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
	logFile       *os.File
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

func (m *StreamManager) StartLogging() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.logFile != nil {
		return nil
	}

	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("could not create log directory: %w", err)
	}

	fileName := time.Now().Format("log_02_01_2006_150405") + ".bin"
	fullPath := filepath.Join(logDir, fileName)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	m.logFile = f
	log.Printf("Started local binary logging to: %s", fileName)
	return nil
}

func (m *StreamManager) StopLogging() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.logFile != nil {
		m.logFile.Close()
		m.logFile = nil
		log.Println("Stopped local binary logging.")
	}
}

func (m *StreamManager) WriteToLog(data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.logFile != nil {
		_, err := m.logFile.Write(data)
		if err != nil {
			log.Printf("Logging error: %v", err)
		}
	}
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
				m.WriteToLog(data)
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

					m.WriteToLog(bytes)
					_, err := m.activePort.Write(bytes)
					if err != nil {
						log.Printf("Write error: %v", err)
					}
				}
				m.mu.RUnlock()

			} else {
				if portCmd, ok := req.(*command.UpdateSerialPortCommand); ok {
					err := m.OpenPort(portCmd.SerialPort, 115200)
					if err != nil {
						log.Printf("Failed to switch to %s: %v", portCmd.SerialPort, err)
					}

				} else if logCmd, ok := req.(*command.LocalLogCommand); ok {
					if logCmd.Command == command.LogStart {
						if err := m.StartLogging(); err != nil {
							log.Printf("Failed to start log: %v", err)
						}
					} else {
						m.StopLogging()
					}
				}
			}
		}
	}()
}
