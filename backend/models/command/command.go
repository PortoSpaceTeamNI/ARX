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
}

type ICommand interface {
	ToString() string
	ToPacket() packet.Packet
	ParseResponse(p packet.Packet) (IResponse, error)
}

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
