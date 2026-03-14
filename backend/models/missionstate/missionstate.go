package missionstate

type MissionState byte

const (
	Idle MissionState = iota
	Filling
	SafeIdle
	FillingN2
	PrePressure
	FillingN2O
	PostPressure
	Ready
	Armed
	Ignition
	Launch
	Flight
	Recovery
	Abort
)

func (ms MissionState) IsAMissionState() bool {
	return ms >= Idle && ms <= Abort
}
