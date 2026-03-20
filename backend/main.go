package main

import (
	"log"
	"missioncontrol/backend/hub"
	"missioncontrol/backend/missioncontrol"
	"missioncontrol/backend/parser"
	"missioncontrol/backend/stream"
	"net/http"
)

func main() {
	serialManager := stream.NewSerialManager(115200)
	p := parser.NewParser()
	mc := missioncontrol.NewMissionControl()
	h := hub.NewHub(serialManager)
	serialManager.SetStatusCallback(func(status stream.SerialStatus) {
		h.BroadcastJSON("serial_status", status)
	})

	go serialManager.Run()
	go p.Run()
	go mc.Run()
	go h.Run()

	http.HandleFunc("/ws", h.WSHandler)
	log.Println("Mission Control Hub starting on :8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
