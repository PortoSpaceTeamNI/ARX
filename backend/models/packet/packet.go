package packet

import (
	"bytes"
	"encoding/binary"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/packet/payload"
	"reflect"
)

const SyncByte byte = 0x55
const CRCSeed = 0xFFFF
const CRCPoly = 0x1021

type Packet struct {
	SyncByte    byte
	SenderID    communicatorid.CommunicatorID
	TargetID    communicatorid.CommunicatorID
	CommandID   commandid.CommandID
	PayloadSize byte
	Payload     payload.IPayload
	CRC         uint16
}

func UpdateCRC(crc *uint16, b byte) {
	*crc ^= (uint16(b) << 8)

	for range 8 {
		if (*crc & 0x8000) != 0 {
			*crc = (*crc << 1) ^ CRCPoly

		} else {
			*crc <<= 1
		}
	}
}

func (p *Packet) CalculateCRC() {
	p.CRC = CRCSeed

	UpdateCRC(&p.CRC, byte(p.SenderID))
	UpdateCRC(&p.CRC, byte(p.TargetID))
	UpdateCRC(&p.CRC, byte(p.CommandID))
	UpdateCRC(&p.CRC, p.PayloadSize)

	if p.Payload != nil {
		for _, b := range p.Payload.ToBytes() {
			UpdateCRC(&p.CRC, b)
		}
	}
}

func NewPacket(target communicatorid.CommunicatorID, cmdID commandid.CommandID, p payload.IPayload) Packet {
	payloadSize := byte(0)
	if p != nil {
		t := reflect.TypeOf(p)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		payloadSize = byte(payload.GetRealSize(t))
	}

	pkt := Packet{
		SyncByte:    SyncByte,
		SenderID:    communicatorid.MissionControl,
		TargetID:    target,
		CommandID:   cmdID,
		PayloadSize: payloadSize,
		Payload:     p,
		CRC:         CRCSeed,
	}

	pkt.CalculateCRC()

	return pkt
}

func NewEmptyPacket(target communicatorid.CommunicatorID, cmdID commandid.CommandID) Packet {
	return NewPacket(target, cmdID, nil)
}

func (p *Packet) ToBytes() []byte {
	buf := new(bytes.Buffer)

	buf.WriteByte(p.SyncByte)
	buf.WriteByte(byte(p.SenderID))
	buf.WriteByte(byte(p.TargetID))
	buf.WriteByte(byte(p.CommandID))
	buf.WriteByte(p.PayloadSize)
	if p.Payload != nil {
		buf.Write(p.Payload.ToBytes())
	}
	binary.Write(buf, binary.BigEndian, p.CRC)

	return buf.Bytes()
}
