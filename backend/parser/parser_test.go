package parser

import (
	"missioncontrol/backend/models/packet"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/packet/payload"
	"testing"
)

func feedBytes(p *Parser, data []byte) (*packet.Packet, error) {
	var lastPacket *packet.Packet
	var lastErr error

	for _, b := range data {
		pkt, err := p.ParseByte(b)
		if err != nil {
			lastErr = err
		}
		if pkt != nil {
			lastPacket = pkt
		}
	}
	return lastPacket, lastErr
}

func TestUpdateCRC(t *testing.T) {
	p := NewParser()

	testData := []byte("123456789")

	for _, b := range testData {
		packet.UpdateCRC(&p.RollingCRC, b)
	}

	expected := uint16(0x29B1)
	if p.RollingCRC != expected {
		t.Errorf("CRC Check Failed! Expected 0x%04X, got 0x%04X", expected, p.RollingCRC)
	}
}

func TestParser_Validators(t *testing.T) {
	p := NewParser()

	t.Run("Invalid Sync Byte", func(t *testing.T) {
		p.Reset()
		_, err := p.ParseByte(0x00)
		if err == nil {
			t.Error("Expected sync error")
		}
	})

	t.Run("Invalid Sender ID", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, 0x02})
		if err == nil {
			t.Error("Expected invalid sender error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})

	t.Run("Invalid Target ID", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, byte(communicatorid.OBC), 0x02})
		if err == nil {
			t.Error("Expected invalid target error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})

	t.Run("Invalid Command ID", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), 0x01})
		if err == nil {
			t.Error("Expected invalid command error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})

	t.Run("Reflection Payload Sizes", func(t *testing.T) {
		expectedSizes := map[string]uintptr{
			"StatusResponsePayload":     36,
			"ManualExecResponsePayload": 2,
		}

		for tType, actualSize := range payload.PayloadSizeRegistry {
			name := tType.Name()
			expectedSize, exists := expectedSizes[name]
			if !exists {
				t.Errorf("Expected sizes map missing entry for registry payload type: %s", name)
				continue
			}

			if actualSize != expectedSize {
				t.Errorf("Size mismatch: Registry has size %d for %s, but expected %d", actualSize, name, expectedSize)
			}
		}

		if len(payload.PayloadSizeRegistry) != len(expectedSizes) {
			t.Errorf("Registry size mismatch: expected %d entries, got %d", len(expectedSizes), len(payload.PayloadSizeRegistry))
		}
	})

	t.Run("Invalid Payload Size", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), byte(commandid.Ack), 255})
		if err == nil {
			t.Error("Expected invalid payload size error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})

	t.Run("Invalid Manual Exec Payload Command ID", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), byte(commandid.Ack), 2, 255, 5})
		if err == nil {
			t.Error("Expected invalid payload command ID error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})

	t.Run("Invalid CRC", func(t *testing.T) {
		p.Reset()

		_, err := feedBytes(p, []byte{packet.SyncByte, byte(communicatorid.OBC), byte(communicatorid.MissionControl), byte(commandid.Ack), 2, byte(commandid.ManualExec), 5, 0xFF, 0xFF})
		if err == nil {
			t.Error("Expected invalid CRC error")
		}

		if p.FieldIndex != 0 {
			t.Error("Parser did not reset after error")
		}
	})
}
