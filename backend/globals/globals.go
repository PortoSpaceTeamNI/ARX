package globals

import (
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/telemetry"
	"time"
)

const (
	RsByteTimeout     time.Duration = 1000 * time.Millisecond
	HeartbeatInterval time.Duration = 1000 * time.Millisecond
	CommandTimeout    time.Duration = 1000 * time.Millisecond
	DoubtInterval     time.Duration = 10 * time.Second // TODO: Doubt state
)

// Communication Channels

// Stream -> Parser
var ByteChannel chan []byte = make(chan []byte, 100)

// Parser -> MissionControl
var ResponseChannel chan command.Response = make(chan command.Response, 100)

// Hub -> MissionControl
var CommandChannel chan command.ICommand = make(chan command.ICommand, 100)

// MissionControl -> Hub
var TelemetryChannel chan telemetry.Telemetry = make(chan telemetry.Telemetry, 100)

// MissionControl -> Stream
var RequestChannel chan command.ICommand = make(chan command.ICommand, 100)
