package payload

import (
	"fmt"
	"missioncontrol/backend/models/missionstate"
	"missioncontrol/backend/models/packet/commandid"
)

type StatusResponsePayload struct {
	NoOpToBytes
	StatusCommandID           byte
	MissionState              byte
	BitFlags                  uint16 // TODO: Maybe check if unused bits are 1
	HydraLFTankPressure       int16
	HydraLFChamberPressure    int16
	HydraFSN2OPressure        int16
	HydraFSN2Pressure         int16
	HydraFSQuickDCPressure    int16
	HydraUFProbeThermo1       int16
	HydraUFProbeThermo2       int16
	HydraUFProbeThermo3       int16
	HydraLFProbeThermo1       int16
	HydraLFProbeThermo2       int16
	HydraLFChamberTemperature int16
	HydraFSN2OTemperature     int16
	LiftFSN2OLoadCell         int32
	LiftRLoadCell1            int32
	LiftRLoadCell2            int32
	LiftRLoadCell3            int32
}

const (
	HydraUFPressurizingValveState uint16 = 1 << iota
	HydraUFVentValveState
	HydraLFAbortValveState
	HydraLFMainValveState
	HydraFSN2OFillValveState
	HydraFSN2OPurgeValveState
	HydraFSN2FillValveState
	HydraFSN2PurgeValveState
	LiftRMainEmatch
	ElytraDrogueParachuteEmatch
	ElytraMainParachuteEmatch
	HydraFSN2OQuickDCValveState
	HydraFSN2QuickDCValveState
)

func (p *StatusResponsePayload) Validate() error {
	cmdID := commandid.CommandID(p.StatusCommandID)
	if cmdID != commandid.Status {
		return fmt.Errorf("Command ID is not Status: 0x%X", p.StatusCommandID)
	}

	missionState := missionstate.MissionState(p.MissionState)
	if !missionState.IsAMissionState() {
		return fmt.Errorf("Invalid MissionState: 0x%X", p.MissionState)
	}

	return nil
}
