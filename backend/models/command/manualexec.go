package command

import (
	"fmt"
	"missioncontrol/backend/models/packet"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/packet/payload"
	"missioncontrol/backend/models/valve"
	"missioncontrol/backend/models/valve/valvestate"
)

type ManualExecCommand struct {
	RemoteCommand
}

type ParsedManualExecResponse struct {
	ManualExecCommandID byte
	//SequenceNumber      byte
}

func (c *ManualExecCommand) ParseResponse(raw packet.Packet) (IResponse, error) {
	payload := raw.Payload.(*payload.ManualExecResponsePayload)

	return ParsedManualExecResponse{
		ManualExecCommandID: payload.ManualExecCommandID,
		//SequenceNumber:      payload.SequenceNumber,
	}, nil
}

type UpdateValveCommand struct {
	ManualExecCommand
	Valve    valve.Valve           `json:"valve"`
	State    valvestate.ValveState `json:"state"`
	Duration uint                  `json:"duration,omitempty"`
}

func (c *UpdateValveCommand) ToString() string {
	return fmt.Sprintf("%s %s Valve", c.State.ToString(), c.Valve.ToString())
}

func (c *UpdateValveCommand) ToPacket() packet.Packet {
	return packet.NewPacket(communicatorid.OBC, commandid.ManualExec, &payload.UpdateValveRequestPayload{
		ManualUpdateValveStateCommandID: byte(commandid.ManualUpdateValveState),
		Valve:                           byte(c.Valve),
		ValveState:                      byte(c.State),
	})
}

// TODO: Update Valve ms and SD commands
