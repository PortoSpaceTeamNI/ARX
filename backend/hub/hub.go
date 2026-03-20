package hub

import (
	"encoding/json"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/stream"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Mu            sync.Mutex
	GroundStation *websocket.Conn
	SerialManager *stream.SerialManager
}

func NewHub(serialManager *stream.SerialManager) *Hub {
	return &Hub{SerialManager: serialManager}
}

type OutgoingMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type SerialConnectMessage struct {
	Port     string `json:"port"`
	BaudRate int    `json:"baudRate"`
}

func (h *Hub) Run() {
	defer close(globals.CommandChannel)

	for data := range globals.TelemetryChannel {
		msg, err := json.Marshal(OutgoingMessage{Type: "telemetry", Data: data})
		if err != nil {
			log.Printf("Error marshaling telemetry: %v", err)
			continue
		}

		h.Broadcast(msg)
	}
}

func (h *Hub) BroadcastJSON(messageType string, data any) {
	msg, err := json.Marshal(OutgoingMessage{Type: messageType, Data: data})
	if err != nil {
		log.Printf("Error marshaling %s message: %v", messageType, err)
		return
	}

	h.Broadcast(msg)
}

func (h *Hub) Broadcast(msg []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if h.GroundStation != nil {
		err := h.GroundStation.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Closing Ground Station connection: %v", err)
			h.GroundStation.Close()
			h.GroundStation = nil
		}
	}
}

// TODO: Restrict origins
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	h.Mu.Lock()
	if h.GroundStation != nil {
		h.GroundStation.Close() // Disconnect previous session if it exists
	}
	h.GroundStation = conn
	h.Mu.Unlock()

	h.broadcastSerialStatus("")

	defer func() {
		h.Mu.Lock()
		if h.GroundStation == conn {
			h.GroundStation = nil
		}
		h.Mu.Unlock()
		conn.Close()
	}()

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		if messageType == websocket.TextMessage {
			var wc command.WebCommand

			if err := json.Unmarshal(payload, &wc); err != nil {
				log.Printf("Failed to parse command envelope: %v", err)
				continue
			}

			switch wc.Type {
			case "list_serial_ports":
				h.broadcastSerialStatus("")
				continue

			case "connect_serial":
				if h.SerialManager == nil {
					continue
				}

				var connectMsg SerialConnectMessage
				if err := json.Unmarshal(wc.Data, &connectMsg); err != nil {
					log.Printf("Failed to parse connect_serial data: %v", err)
					h.broadcastSerialStatus("Invalid connect_serial payload")
					continue
				}

				if connectMsg.Port == "" {
					h.broadcastSerialStatus("Serial port is required")
					continue
				}

				if err := h.SerialManager.Connect(connectMsg.Port, connectMsg.BaudRate); err != nil {
					log.Printf("Failed to connect serial port %s: %v", connectMsg.Port, err)
					h.broadcastSerialStatus(err.Error())
					continue
				}

				h.broadcastSerialStatus("")
				continue

			case "disconnect_serial":
				if h.SerialManager == nil {
					continue
				}

				if err := h.SerialManager.Disconnect(); err != nil {
					log.Printf("Failed to disconnect serial port: %v", err)
					h.broadcastSerialStatus(err.Error())
					continue
				}

				h.broadcastSerialStatus("")
				continue
			}

			t, exists := command.CommandRegistry[wc.Type]
			if !exists {
				log.Printf("Unknown command type: %s", wc.Type)
				continue
			}

			cmd := reflect.New(t).Interface().(command.ICommand)
			if err := json.Unmarshal(wc.Data, cmd); err != nil {
				log.Printf("Failed to parse command data for %s: %v", wc.Type, err)
				continue
			}

			if globals.CommandChannel != nil {
				globals.CommandChannel <- cmd
			}
		}
	}
}

func (h *Hub) broadcastSerialStatus(errorMessage string) {
	if h.SerialManager == nil {
		return
	}

	status := h.SerialManager.GetStatus()
	if errorMessage != "" {
		status.Error = errorMessage
	}

	h.BroadcastJSON("serial_status", status)
}
