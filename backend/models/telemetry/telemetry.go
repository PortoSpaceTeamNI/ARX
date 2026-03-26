package telemetry

import "missioncontrol/backend/models/command"

type Telemetry struct {
	PacketLoss     int                           `json:"packetLoss"`
	Latency        float64                       `json:"latency"`
	DataRate       int                           `json:"dataRate"`
	Status         *command.ParsedStatusResponse `json:"status"`
	CommandLog     string                        `json:"commandLog"`
	AvailablePorts []string                      `json:"availablePorts"`
	CurrentPort    string                        `json:"currentPort"`
}
