package payload

import (
	"fmt"
	"missioncontrol/backend/models/packet/commandid"
)

type ManualExecResponsePayload struct {
	NoOpToBytes
	ManualExecCommandID byte
	//SequenceNumber      byte
}

func (p *ManualExecResponsePayload) Validate() error {
	cmdID := commandid.CommandID(p.ManualExecCommandID)
	if cmdID != commandid.ManualExec {
		return fmt.Errorf("Command ID is not ManualExec: 0x%X", p.ManualExecCommandID)
	}

	return nil
}
