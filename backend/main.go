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
	s, err := stream.NewSerialStream("/dev/ttyACM0", 115200)
	if err != nil {
		panic("Serial err: " + err.Error())
	}
	p := parser.NewParser()
	mc := missioncontrol.NewMissionControl()
	h := hub.NewHub()

	go stream.Run(s)
	go p.Run()
	go mc.Run()
	go h.Run()

	http.HandleFunc("/ws", h.WSHandler)
	log.Println("Mission Control Hub starting on :8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
