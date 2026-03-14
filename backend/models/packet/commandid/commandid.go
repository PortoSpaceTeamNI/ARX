package commandid

type CommandID byte

const (
	None CommandID = iota
	Status
	Abort
	Stop
	Arm
	Fire
	FillExec
	ManualExec
	Ack
	Nack
)

type ManualCommandID byte

const (
	ManualStartSDLog ManualCommandID = iota
	ManualStopSDLog
	ManualSDStatus
	ManualUpdateValveState
	ManualValveMs
	ManualAck
)
