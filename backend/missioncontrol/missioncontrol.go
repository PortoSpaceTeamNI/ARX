package missioncontrol

import (
	"fmt"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/telemetry"
	"missioncontrol/backend/models/valve"
	"missioncontrol/backend/models/valve/valvestate"
	"missioncontrol/backend/stream"
	"time"

	"go.bug.st/serial"
)

type MissionControl struct {
	CommandQueue  []command.ICommand
	PendingValves map[valve.Valve]valvestate.ValveState
	Telemetry     telemetry.Telemetry

	CurrentCommand      command.ICommand
	LastCommandSentTime time.Time
}

func NewMissionControl() *MissionControl {
	return &MissionControl{
		PendingValves: make(map[valve.Valve]valvestate.ValveState),
	}
}

func NewLiveMissionControl() *MissionControl {
	mc := NewMissionControl()

	port, err := stream.FindAvailablePort()
	if err != nil {
		log.Printf("Mission Control initialized in IDLE mode: %v", err)

	} else {
		log.Printf("Valid Port Found: %s", port)

		updateCmd := &command.UpdateSerialPortCommand{
			SerialPort: port,
		}

		globals.RequestChannel <- updateCmd
		mc.Telemetry.CurrentPort = port
	}

	return mc
}

func (mc *MissionControl) Run() {
	defer close(globals.TelemetryChannel)
	defer close(globals.RequestChannel)

	heartbeatTicker := time.NewTicker(globals.HeartbeatInterval)
	defer heartbeatTicker.Stop()

	timeoutTicker := time.NewTicker(globals.CommandTimeout)
	defer timeoutTicker.Stop()

	for {
		select {
		case res, ok := <-globals.ResponseChannel:
			if !ok { // ResponseChannel was closed
				return
			}

			mc.HandleResponse(res)

		case cmd, ok := <-globals.CommandChannel:
			if !ok { // CommandChannel was closed
				return
			}

			if mc.CurrentCommand == nil {
				mc.HandleCommand(cmd)

			} else {
				mc.CommandQueue = append(mc.CommandQueue, cmd)
			}

		case <-heartbeatTicker.C:
			mc.Telemetry.AvailablePorts, _ = serial.GetPortsList()
			mc.Heartbeat()

		case <-timeoutTicker.C:
			mc.CheckTimeout()
		}
	}
}

func (mc *MissionControl) HandleResponse(res command.Response) {
	if res.Data == nil { // Not a command response but a status for the network conditions
		mc.Telemetry.DataRate = res.DataRate
		return
	}

	mc.Telemetry.Latency = float64(res.SyncByteTime.Sub(mc.LastCommandSentTime).Microseconds()) / 1e3
	//if mc.CurrentCommand != nil {
	mc.Telemetry.CommandLog = fmt.Sprintf("[RCV]: %s Ack", mc.CurrentCommand.ToString())
	//}

	switch r := res.Data.(type) {
	case command.ParsedManualExecResponse: // TODO: This implementation assumes that, when currentCommand has a response type of ManualExec then the first received response will be relative to that command, which migh not be the case
		if valveCmd, ok := mc.CurrentCommand.(*command.UpdateValveCommand); ok {
			pendingState := mc.PendingValves[valveCmd.Valve]

			if valveCmd.State == valvestate.Opened && pendingState == valvestate.OpeningNotAcked {
				mc.PendingValves[valveCmd.Valve] = valvestate.Opening

			} else if valveCmd.State == valvestate.Closed && pendingState == valvestate.ClosingNotAcked {
				mc.PendingValves[valveCmd.Valve] = valvestate.Closing
			}
		}

		log.Println("Received manual exec response")

		mc.CurrentCommand = nil

	case command.ParsedStatusResponse:
		if _, ok := mc.CurrentCommand.(*command.StatusCommand); ok {
			mc.CurrentCommand = nil
		}

		log.Println("Received status response")

		statusCopy := r
		mc.Telemetry.Status = &statusCopy

		//mc.resolveValvesState(mc.Telemetry.Status)
		for valve := range mc.PendingValves {
			delete(mc.PendingValves, valve)
		}
	}

	globals.TelemetryChannel <- mc.Telemetry
	mc.HandleQueue()
}

func (mc *MissionControl) sendCommand(cmd command.ICommand) {
	log.Printf("Sending %s command", cmd.ToString())

	mc.LastCommandSentTime = time.Now()

	globals.RequestChannel <- cmd
}

func (mc *MissionControl) HandleCommand(cmd command.ICommand) {
	if cmd.IsRemote() {
		if valveCmd, ok := cmd.(*command.UpdateValveCommand); ok {
			switch valveCmd.State {
			case valvestate.Closed:
				mc.PendingValves[valveCmd.Valve] = valvestate.ClosingNotAcked

			case valvestate.Opened:
				mc.PendingValves[valveCmd.Valve] = valvestate.OpeningNotAcked
			}
		}

		mc.CurrentCommand = cmd
	}

	mc.sendCommand(cmd)
}

func (mc *MissionControl) HandleQueue() {
	if len(mc.CommandQueue) == 0 {
		return
	}

	nextCmd := mc.CommandQueue[0]
	mc.CommandQueue = mc.CommandQueue[1:]

	mc.HandleCommand(nextCmd)
}

func (mc *MissionControl) Heartbeat() {
	if mc.CurrentCommand != nil {
		return
	}

	if len(mc.CommandQueue) > 0 {
		mc.HandleQueue()
		return
	}

	//log.Println("Sending status command")
	mc.HandleCommand(&command.StatusCommand{})
}

func (mc *MissionControl) CheckTimeout() {
	if mc.CurrentCommand == nil || time.Since(mc.LastCommandSentTime) <= globals.CommandTimeout {
		return
	}

	portCmdIndex := -1
	for i, cmd := range mc.CommandQueue {
		if _, ok := cmd.(*command.UpdateSerialPortCommand); ok {
			portCmdIndex = i
			break
		}
	}

	if portCmdIndex != -1 {
		log.Println("Timeout: Jumping queue with Port Update")

		portCmd := mc.CommandQueue[portCmdIndex]
		mc.CommandQueue = append(mc.CommandQueue[:portCmdIndex], mc.CommandQueue[portCmdIndex+1:]...)

		globals.RequestChannel <- portCmd

		mc.Telemetry.CommandLog = "[SYS]: Attempting rescue port swap"
		globals.TelemetryChannel <- mc.Telemetry

		return
	}

	log.Println("Retrying command")
	mc.Telemetry.PacketLoss++
	mc.Telemetry.Latency = 0
	mc.Telemetry.CommandLog = fmt.Sprintf("[ERROR]: Retrying %s", mc.CurrentCommand.ToString())
	globals.TelemetryChannel <- mc.Telemetry

	mc.sendCommand(mc.CurrentCommand)
}
