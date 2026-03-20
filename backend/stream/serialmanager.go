package stream

import (
	"io"
	"log"
	"missioncontrol/backend/globals"
	"sort"
	"sync"

	"go.bug.st/serial"
)

type SerialStatus struct {
	Ports      []string `json:"ports"`
	Connected  bool     `json:"connected"`
	CurrentPort string   `json:"currentPort,omitempty"`
	BaudRate   int      `json:"baudRate,omitempty"`
	Error      string   `json:"error,omitempty"`
}

type SerialManager struct {
	mu          sync.RWMutex
	stream      *SerialStream
	currentPort string
	baudRate    int
	defaultBaud int
	onStatus    func(SerialStatus)
}

func NewSerialManager(defaultBaud int) *SerialManager {
	return &SerialManager{defaultBaud: defaultBaud}
}

func (m *SerialManager) SetStatusCallback(cb func(SerialStatus)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onStatus = cb
}

func (m *SerialManager) ListPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}

	sort.Strings(ports)
	return ports, nil
}

func (m *SerialManager) GetStatus() SerialStatus {
	ports, err := m.ListPorts()
	status := SerialStatus{Ports: ports}
	if err != nil {
		status.Error = err.Error()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.stream != nil {
		status.Connected = true
		status.CurrentPort = m.currentPort
		status.BaudRate = m.baudRate
	}

	return status
}

func (m *SerialManager) Connect(portName string, baudRate int) error {
	if baudRate <= 0 {
		baudRate = m.defaultBaud
	}

	newStream, err := NewSerialStream(portName, baudRate)
	if err != nil {
		return err
	}

	m.mu.Lock()
	oldStream := m.stream
	m.stream = newStream
	m.currentPort = portName
	m.baudRate = baudRate
	m.mu.Unlock()

	if oldStream != nil {
		_ = oldStream.Close()
	}

	go m.readLoop(newStream, portName)
	m.emitStatus("")

	return nil
}

func (m *SerialManager) Disconnect() error {
	m.mu.Lock()
	active := m.stream
	m.stream = nil
	m.currentPort = ""
	m.baudRate = 0
	m.mu.Unlock()

	if active != nil {
		if err := active.Close(); err != nil {
			m.emitStatus(err.Error())
			return err
		}
	}

	m.emitStatus("")
	return nil
}

func (m *SerialManager) Run() {
	for packetBytes := range globals.RequestChannel {
		m.mu.RLock()
		active := m.stream
		m.mu.RUnlock()

		if active == nil {
			log.Println("Serial write skipped: no active serial connection")
			continue
		}

		if _, err := active.Write(packetBytes); err != nil {
			log.Printf("Serial write error: %v", err)
			m.handleStreamFailure(active, err)
		}
	}
}

func (m *SerialManager) readLoop(active *SerialStream, portName string) {
	buf := make([]byte, 256)

	for {
		n, err := active.Read(buf)

		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			globals.ByteChannel <- data
		}

		if err != nil {
			if err != io.EOF {
				log.Printf("Serial read error (%s): %v", portName, err)
			}

			m.handleStreamFailure(active, err)
			return
		}
	}
}

func (m *SerialManager) handleStreamFailure(active *SerialStream, err error) {
	m.mu.Lock()
	if m.stream != active {
		m.mu.Unlock()
		return
	}

	m.stream = nil
	m.currentPort = ""
	m.baudRate = 0
	m.mu.Unlock()

	_ = active.Close()

	if err == io.EOF {
		m.emitStatus("")
		return
	}

	m.emitStatus(err.Error())
}

func (m *SerialManager) emitStatus(errorMessage string) {
	m.mu.RLock()
	cb := m.onStatus
	m.mu.RUnlock()

	if cb == nil {
		return
	}

	status := m.GetStatus()
	if errorMessage != "" {
		status.Error = errorMessage
	}

	cb(status)
}
