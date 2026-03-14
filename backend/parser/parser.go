package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"maps"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/packet"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/packet/packetfield"
	"missioncontrol/backend/models/packet/payload"
	"reflect"
	"slices"
	"time"
)

type Field struct {
	Field   packetfield.PacketField
	Size    int
	ParseFn func(buf []byte, p *packet.Packet, fields []Field) error
}

type Parser struct {
	FieldIndex       int
	FieldBuf         []byte
	Fields           []Field
	Packet           packet.Packet
	RollingCRC       uint16
	SyncByteTime     time.Time
	LastByteTime     time.Time
	AccumulatedBytes int
}

func NewParser() *Parser {
	p := &Parser{}
	p.Fields = []Field{
		{
			Field: packetfield.Sync,
			Size:  1,
			ParseFn: func(buf []byte, _ *packet.Packet, _ []Field) error {
				if buf[0] != packet.SyncByte {
					return fmt.Errorf("Non sync byte ignored: 0x%X", buf[0])
				}

				p.SyncByteTime = time.Now()
				return nil
			},
		},
		{
			Field: packetfield.SenderID,
			Size:  1,
			ParseFn: func(buf []byte, packet *packet.Packet, _ []Field) error {
				packet.SenderID = communicatorid.CommunicatorID(buf[0])
				if packet.SenderID != communicatorid.OBC {
					return fmt.Errorf("Packet sender is not OBC")
				}

				return nil
			},
		},
		{
			Field: packetfield.TargetID,
			Size:  1,
			ParseFn: func(buf []byte, packet *packet.Packet, _ []Field) error {
				packet.TargetID = communicatorid.CommunicatorID(buf[0])
				if packet.TargetID != communicatorid.MissionControl {
					return fmt.Errorf("Packet target is not Mission Control")
				}

				return nil
			},
		},
		{
			Field: packetfield.CommandID,
			Size:  1,
			ParseFn: func(buf []byte, packet *packet.Packet, _ []Field) error {
				packet.CommandID = commandid.CommandID(buf[0])
				if packet.CommandID != commandid.Ack {
					return fmt.Errorf("Packet command ID is not Ack: 0x%X", packet.CommandID)
				}

				return nil
			},
		},
		{
			Field: packetfield.PayloadSize,
			Size:  1,
			ParseFn: func(buf []byte, packet *packet.Packet, fields []Field) error {
				packet.PayloadSize = buf[0]

				validSizes := slices.Collect(maps.Values(payload.PayloadSizeRegistry))
				if !slices.Contains(validSizes, uintptr(packet.PayloadSize)) {
					return fmt.Errorf("Invalid payload size: %X", packet.PayloadSize)
				}

				fields[packetfield.Payload].Size = int(packet.PayloadSize)
				return nil
			},
		},
		{
			Field: packetfield.Payload,
			Size:  0, // Size 0 means the value is dynamic. In this case it gets overwritten when payloadSize is parsed
			ParseFn: func(buf []byte, packet *packet.Packet, _ []Field) error {
				newPayload, err := p.determinePayloadType(packet, buf)
				if err != nil {
					return err
				}

				reader := bytes.NewReader(buf)
				if err := binary.Read(reader, binary.BigEndian, newPayload); err != nil {
					return fmt.Errorf("Failed to map buffer: %w", err)
				}

				if err := newPayload.Validate(); err != nil {
					return fmt.Errorf("Payload validation failed: %w", err)
				}

				packet.Payload = newPayload
				return nil
			},
		},
		{
			Field: packetfield.CRC,
			Size:  2,
			ParseFn: func(buf []byte, packet *packet.Packet, _ []Field) error {
				receivedCRC := binary.BigEndian.Uint16(buf)

				if p.RollingCRC != receivedCRC {
					//return fmt.Errorf("CRC Mismatch: calculated %04X, received %04X", p.RollingCRC, receivedCRC) // TODO: Uncomment to do the crc check
				}

				packet.CRC = receivedCRC
				return nil
			},
		},
	}
	p.Reset()
	return p
}

func (p *Parser) determinePayloadType(packet *packet.Packet, buf []byte) (payload.IPayload, error) {
	if commandid.CommandID(buf[0]) == commandid.Status && packet.PayloadSize == byte(payload.PayloadSizeRegistry[reflect.TypeFor[payload.StatusResponsePayload]()]) {
		return &payload.StatusResponsePayload{}, nil

	} else if commandid.CommandID(buf[0]) == commandid.ManualExec && packet.PayloadSize == byte(payload.PayloadSizeRegistry[reflect.TypeFor[payload.ManualExecResponsePayload]()]) {
		return &payload.ManualExecResponsePayload{}, nil

	} else {
		return nil, fmt.Errorf("Invalid payload type")
	}
}

func (p *Parser) Reset() {
	p.FieldIndex = 0
	p.FieldBuf = nil
	p.RollingCRC = packet.CRCSeed
	p.Packet = packet.Packet{}
}

func (p *Parser) CheckByteTimeout() {
	if p.FieldIndex != int(packetfield.Sync) && !p.LastByteTime.IsZero() {
		duration := time.Since(p.LastByteTime)

		if duration > globals.RsByteTimeout {
			log.Printf("Parser timeout: %v since last byte", duration)
			p.Reset()
		}
	}
}

func (p *Parser) ParseByte(rxByte byte) (*packet.Packet, error) {
	p.CheckByteTimeout()

	field := &p.Fields[p.FieldIndex]

	p.LastByteTime = time.Now()
	p.FieldBuf = append(p.FieldBuf, rxByte)

	switch field.Field {
	case packetfield.SenderID, packetfield.TargetID, packetfield.CommandID, packetfield.PayloadSize, packetfield.Payload:
		packet.UpdateCRC(&p.RollingCRC, rxByte)
	}

	// Field isnt fully parsed yet
	if len(p.FieldBuf) < field.Size && field.Size != 0 {
		return nil, nil
	}

	if err := field.ParseFn(p.FieldBuf, &p.Packet, p.Fields); err != nil {
		p.Reset()
		return nil, err
	}

	// Field was successfully parsed
	p.FieldBuf = nil
	p.FieldIndex++

	if p.FieldIndex == len(p.Fields) {
		packet := p.Packet
		p.Reset()
		return &packet, nil
	}

	return nil, nil
}

func (p *Parser) ParsePacket(pkt *packet.Packet) (command.Response, error) {
	var data command.IResponse
	var err error

	switch pkt.Payload.(type) {
	case *payload.StatusResponsePayload:
		cmd := &command.StatusCommand{}
		data, err = cmd.ParseResponse(*pkt)

	case *payload.ManualExecResponsePayload:
		cmd := &command.ManualExecCommand{}
		data, err = cmd.ParseResponse(*pkt)

	default:
		return command.Response{}, fmt.Errorf("unsupported payload: %T", pkt.Payload)
	}

	if err != nil {
		return command.Response{}, fmt.Errorf("parse error: %w", err)
	}

	return command.Response{
		SyncByteTime: p.SyncByteTime,
		DataRate:     -1,
		Data:         data,
	}, nil
}

func (p *Parser) Run() {
	defer close(globals.ResponseChannel)

	timeoutTicker := time.NewTicker(globals.RsByteTimeout / 2)
	defer timeoutTicker.Stop()

	networkPerformanceTicker := time.NewTicker(1 * time.Second)
	defer networkPerformanceTicker.Stop()

	for {
		select {
		case chunk, ok := <-globals.ByteChannel:
			if !ok {
				return
			}

			p.AccumulatedBytes += len(chunk)

			for _, b := range chunk {
				rawPacket, err := p.ParseByte(b)
				if err != nil {
					log.Printf("Parser error: %v", err)
					continue
				}

				if rawPacket != nil {
					res, err := p.ParsePacket(rawPacket)
					if err != nil {
						log.Printf("Failed to parse packet: %v", err)
						continue
					}

					globals.ResponseChannel <- res
				}
			}

		case <-timeoutTicker.C:
			p.CheckByteTimeout()

		case <-networkPerformanceTicker.C:
			bps := p.AccumulatedBytes
			p.AccumulatedBytes = 0

			globals.ResponseChannel <- command.Response{
				SyncByteTime: time.Now(),
				DataRate:     bps,
				Data:         nil,
			}
		}
	}
}
