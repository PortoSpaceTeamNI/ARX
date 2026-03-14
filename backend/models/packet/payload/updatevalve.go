package payload

type UpdateValveRequestPayload struct {
	NoOpValidate
	ManualUpdateValveStateCommandID byte
	Valve                           byte
	ValveState                      byte
}

func (p *UpdateValveRequestPayload) ToBytes() []byte {
	return []byte{p.ManualUpdateValveStateCommandID, p.Valve, p.ValveState}
}
