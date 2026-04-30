package missioncontrol

import (
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/telemetry"
	"missioncontrol/backend/stream"
	"time"

	"go.bug.st/serial"
)

type PortMonitor struct {
	currentPort    string
	UpdatePortChan chan string
	RespChan       chan []telemetry.AvailablePort
}

func NewPortMonitor() *PortMonitor {
	return &PortMonitor{
		currentPort:    "",
		UpdatePortChan: make(chan string, 1),
		RespChan:       make(chan []telemetry.AvailablePort, 1),
	}
}

func (pm *PortMonitor) Run() {
	ticker := time.NewTicker(globals.PortScanInterval)
	defer ticker.Stop()

	pm.scan()

	for {
		select {
		case port := <-pm.UpdatePortChan:
			pm.currentPort = port
			pm.scan()

		case <-ticker.C:
			pm.scan()
		}
	}
}

func (pm *PortMonitor) scan() {
	ports, err := serial.GetPortsList()
	if err != nil {
		return
	}

	cmd := &command.StatusCommand{}
	pkt := cmd.ToPacket()
	probeBytes := pkt.ToBytes()
	expectedHeader := []byte{pkt.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), byte(commandid.Ack)}

	res := []telemetry.AvailablePort{}
	for _, portName := range ports {
		if portName == pm.currentPort {
			continue
		}

		valid, _ := stream.ProbePort(portName, probeBytes, expectedHeader)
		res = append(res, telemetry.AvailablePort{Port: portName, State: valid})
	}

	select {
	case pm.RespChan <- res:
	default:
		<-pm.RespChan
		pm.RespChan <- res
	}
}
