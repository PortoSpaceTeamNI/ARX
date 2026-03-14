package telemetry

import "missioncontrol/backend/models/command"

type Telemetry struct {
	Latency    float64                       `json:"latency"`
	DataRate   int                           `json:"dataRate"`
	Status     *command.ParsedStatusResponse `json:"status"`
	CommandLog string                        `json:"commandLog"`
}
