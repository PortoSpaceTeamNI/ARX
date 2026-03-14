package missioncontrol

import (
	"fmt"
	"log"
	"missioncontrol/backend/globals"
	"missioncontrol/backend/models/command"
	"missioncontrol/backend/models/telemetry"
	"missioncontrol/backend/models/valve"
	"missioncontrol/backend/models/valve/valvestate"
	"time"
)

type MissionControl struct {
	CurrentCommand      command.ICommand
	CommandQueue        []command.ICommand
	PendingValves       map[valve.Valve]valvestate.ValveState
	Telemetry           telemetry.Telemetry
	LastCommandSentTime time.Time
}

func NewMissionControl() *MissionControl {
	return &MissionControl{
		PendingValves: make(map[valve.Valve]valvestate.ValveState),
	}
}

func (mc *MissionControl) Run() {
	defer close(globals.TelemetryChannel)
	defer close(globals.RequestChannel)

	heartbeat := time.NewTicker(globals.HeartbeatInterval)
	defer heartbeat.Stop()

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

		case <-heartbeat.C:
			mc.OnTick()
		}
	}
}

func (mc *MissionControl) HandleResponse(res command.Response) {
	if res.Data == nil { // Not a command response but a status for the network conditions
		mc.Telemetry.DataRate = res.DataRate
		return
	}

	mc.Telemetry.Latency = float64(res.SyncByteTime.Sub(mc.LastCommandSentTime).Microseconds()) / 1e3
	mc.Telemetry.CommandLog = fmt.Sprintf("[RCV]: %s", mc.CurrentCommand.ToString())

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

func (mc *MissionControl) HandleCommand(cmd command.ICommand) {
	if valveCmd, ok := cmd.(*command.UpdateValveCommand); ok {
		switch valveCmd.State {
		case valvestate.Closed:
			mc.PendingValves[valveCmd.Valve] = valvestate.ClosingNotAcked

		case valvestate.Opened:
			mc.PendingValves[valveCmd.Valve] = valvestate.OpeningNotAcked
		}
	}

	mc.CurrentCommand = cmd

	pkt := cmd.ToPacket()
	bytes := pkt.ToBytes()

	log.Printf("Sending %s command", cmd.ToString())

	mc.LastCommandSentTime = time.Now()

	globals.RequestChannel <- bytes
}

func (mc *MissionControl) OnTick() {
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

func (mc *MissionControl) HandleQueue() {
	if len(mc.CommandQueue) == 0 {
		return
	}

	nextCmd := mc.CommandQueue[0]
	mc.CommandQueue = mc.CommandQueue[1:]

	mc.HandleCommand(nextCmd)
}

/*func (mc *MissionControl) resolveValvesState(data *command.ParsedStatusResponse) {
	resolveValve := func(valve valve.Valve, currentState *valvestate.ValveState) {
		pendingState, isPending := mc.PendingValves[valve]
		if !isPending {
			return
		}

		if (pendingState == valvestate.Opening && *currentState == valvestate.Opened) || // Valve has reached the desired state so we remove it from the map
			(pendingState == valvestate.Closing && *currentState == valvestate.Closed) {
			delete(mc.PendingValves, valve)

		} else { // Valve hasn't reached the desired state so we overwrite the data with the pending state
			*currentState = pendingState
		}
	}

	resolveValve(valve.Pressurizing, &data.HydraUF.PressurizingValve)
	resolveValve(valve.Vent, &data.HydraUF.VentValve)

	resolveValve(valve.Abort, &data.HydraLF.AbortValve)
	resolveValve(valve.Main, &data.HydraLF.MainValve)

	resolveValve(valve.N2OFill, &data.HydraFS.N2O.FillValve)
	resolveValve(valve.N2OPurge, &data.HydraFS.N2O.PurgeValve)
	resolveValve(valve.N2Fill, &data.HydraFS.N2.FillValve)
	resolveValve(valve.N2Purge, &data.HydraFS.N2.PurgeValve)
	resolveValve(valve.N2OQuickDc, &data.HydraFS.QuickDC.N2OValve)
	resolveValve(valve.N2QuickDc, &data.HydraFS.QuickDC.N2Valve)
}
*/
