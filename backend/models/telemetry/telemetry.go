package telemetry

import (
	"fmt"
	"missioncontrol/backend/models/command"
	"strings"
)

type Telemetry struct {
	PacketLoss     int                           `json:"packetLoss"`
	Latency        float64                       `json:"latency"`
	DataRate       int                           `json:"dataRate"`
	Status         *command.ParsedStatusResponse `json:"status"`
	CommandLog     string                        `json:"commandLog"`
	AvailablePorts []string                      `json:"availablePorts"`
	CurrentPort    string                        `json:"currentPort"`
}

func (t *Telemetry) ToGrafanaString() string {
	header := "telemetry,source=backend"

	fields := []string{
		fmt.Sprintf("latency=%f", t.Latency),
	}

	return fmt.Sprintf("%s %s", header, strings.Join(fields, ","))
}
