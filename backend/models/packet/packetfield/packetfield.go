package packetfield

type PacketField byte

// It's important that these respect the order in which they appear in the packet
const (
	Sync PacketField = iota
	SenderID
	TargetID
	CommandID
	PayloadSize
	Payload
	CRC
)
