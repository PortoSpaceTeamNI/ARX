package stream

import (
	"bytes"
	"fmt"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"

	"go.bug.st/serial"
)

type SerialStream struct {
}

func FindAvailablePort() (string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return "", err
	}

	cmd := &command.StatusCommand{}
	pkt := cmd.ToPacket()
	bytes := pkt.ToBytes()

	expectedHeader := []byte{pkt.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), byte(commandid.Ack)}

	for _, portName := range ports {
		found, err := probePort(portName, bytes, expectedHeader)
		if err != nil {
			log.Printf("Port %s failed probe: %v", portName, err)
			continue
		}

		if found {
			return portName, nil
		}
	}

	return "", fmt.Errorf("No valid device responded to status probe")
}

func probePort(portName string, probe []byte, expectedHeader []byte) (bool, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	p, err := serial.Open(portName, mode)
	if err != nil {
		return false, err
	}
	defer p.Close()

	p.SetReadTimeout(globals.CommandTimeout)

	_, err = p.Write(probe)
	if err != nil {
		return false, err
	}

	buf := make([]byte, 128)
	n, err := p.Read(buf)
	if err != nil || n < len(expectedHeader) {
		return false, nil
	}

	if bytes.HasPrefix(buf[:n], expectedHeader) {
		return true, nil
	}

	return false, nil
}
