package hub

import (
	"encoding/json"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/telemetry"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Mu            sync.Mutex
	GroundStation *websocket.Conn

	GrafanaQueue chan telemetry.Telemetry
}

func NewHub() *Hub {
	return &Hub{
		GrafanaQueue: make(chan telemetry.Telemetry, 100),
	}
}

func (h *Hub) Run() {
	defer close(globals.CommandChannel)

	go h.RunGrafanaPusher()

	for data := range globals.TelemetryChannel {
		msg, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling telemetry: %v", err)
			continue
		}

		h.Broadcast(msg)

		if data.Status != nil {
			select {
			case h.GrafanaQueue <- data:
			default:
			}
		}
	}
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

func (h *Hub) RunGrafanaPusher() {
	url := "http://grafana:3000/api/live/push/telemetry_stream"

	for t := range h.GrafanaQueue {
		req, err := http.NewRequest("POST", url, strings.NewReader(t.ToGrafanaString()))
		if err != nil {
			log.Printf("Failed to create Grafana request: %v", err)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to push to Grafana: %v", err)
			continue
		}
		resp.Body.Close()
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
		log.Println("Upgrade error: %v", err)
		return
	}

	h.Mu.Lock()
	if h.GroundStation != nil {
		h.GroundStation.Close() // Disconnect previous session if it exists
	}
	h.GroundStation = conn
	h.Mu.Unlock()

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
			log.Println("Read error: %v", err)
			break
		}

		if messageType == websocket.TextMessage {
			var wc command.WebCommand

			if err := json.Unmarshal(payload, &wc); err != nil {
				log.Println("Failed to parse command envelope: %v", err)
				continue
			}

			t, exists := command.CommandRegistry[wc.Type]
			if !exists {
				log.Println("Unknown command type: %s", wc.Type)
				continue
			}

			cmd := reflect.New(t).Interface().(command.ICommand)
			if err := json.Unmarshal(wc.Data, cmd); err != nil {
				log.Println("Failed to parse command data for %s: %v", wc.Type, err)
				continue
			}

			if globals.CommandChannel != nil {
				globals.CommandChannel <- cmd
			}
		}
	}
}
