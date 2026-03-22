package command

import (
	"encoding/json"
	"missioncontrol/backend/models/packet"
	"reflect"
	"time"
)

type WebCommand struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

var CommandRegistry = make(map[string]reflect.Type)

func init() {
	register := func(name string, cmd ICommand) {
		t := reflect.TypeOf(cmd).Elem()
		CommandRegistry[name] = t
	}

	register("update_valve", &UpdateValveCommand{})
	register("update_serial_port", &UpdateSerialPortCommand{})
}

type ICommand interface {
	ToString() string
	IsRemote() bool

	ToPacket() packet.Packet
	ParseResponse(p packet.Packet) (IResponse, error)
}

type RemoteCommand struct{}

func (c *RemoteCommand) IsRemote() bool { return true }

type LocalCommand struct{}

func (c *LocalCommand) IsRemote() bool                                   { return false }
func (c *LocalCommand) ToPacket() packet.Packet                          { return packet.Packet{} }
func (c *LocalCommand) ParseResponse(p packet.Packet) (IResponse, error) { return nil, nil }

type IResponse interface {
	IsResponse()
}

func (ParsedManualExecResponse) IsResponse() {}
func (ParsedStatusResponse) IsResponse()     {}

type Response struct {
	SyncByteTime time.Time
	DataRate     int
	Data         IResponse
}
