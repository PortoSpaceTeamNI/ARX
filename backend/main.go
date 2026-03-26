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
	s := stream.NewStreamManager()
	p := parser.NewParser()
	mc := missioncontrol.NewLiveMissionControl()
	h := hub.NewHub()

	go s.Run()
	go p.Run()
	go mc.Run()
	go h.Run()

	http.HandleFunc("/ws", h.WSHandler)
	log.Println("Mission Control Hub starting on :8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
