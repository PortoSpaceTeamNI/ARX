package telemetry

import (
	"fmt"
	"missioncontrol/backend/models/command"
	"strings"
)

type Telemetry struct {
	PacketLoss     int                           `json:"packetLoss"`
	Latency        float64                       `json:"latency"`
	DataRate       int                           `json:"dataRate"`
	Status         *command.ParsedStatusResponse `json:"status"`
	CommandLog     string                        `json:"commandLog"`
	AvailablePorts []string                      `json:"availablePorts"`
	CurrentPort    string                        `json:"currentPort"`
}

func (t *Telemetry) ToGrafanaString() string {
	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	header := "telemetry,source=backend"

	fields := []string{
		// Filling
		fmt.Sprintf("N2O\\ Tank\\ Pressure=%f", t.Status.HydraFS.N2O.Pressure),
		// TODO: Loadcell tanque
		fmt.Sprintf("N2O\\ Bottle\\ Weight=%f", t.Status.LiftFS.N2OLoadcell),
		fmt.Sprintf("Vent\\ Valve=%d", t.Status.HydraUF.VentValve.Binarize()),
		fmt.Sprintf("N2O\\ Fill\\ Valve=%d", t.Status.HydraFS.N2O.FillValve.Binarize()),
		fmt.Sprintf("N2\\ Fill\\ Valve=%d", t.Status.HydraFS.N2.FillValve.Binarize()),

		// Launch
		fmt.Sprintf("Lift\\ Thrust\\ 1\\ Weight=%f", t.Status.LiftR.Loadcell1),
		fmt.Sprintf("Lift\\ Thrust\\ 2\\ Weight=%f", t.Status.LiftR.Loadcell2),
		fmt.Sprintf("Lift\\ Thrust\\ 3\\ Weight=%f", t.Status.LiftR.Loadcell3),
		fmt.Sprintf("Main\\ Valve=%d", t.Status.HydraLF.MainValve.Binarize()),
		fmt.Sprintf("Main\\ Parachute\\ Ematch=%d", boolToInt(t.Status.Elytra.MainChuteEmatch)),
		fmt.Sprintf("Drogue\\ Parachute\\ Ematch=%d", boolToInt(t.Status.Elytra.DrogueChuteEmatch)),
		fmt.Sprintf("Chamber\\ Pressure=%f", t.Status.HydraLF.ChamberPressure),
		fmt.Sprintf("Altitude=%f", t.Status.Navigator.GPS.Altitude),
		fmt.Sprintf("Velocity\\ Z=%f", t.Status.Navigator.Kalman.VelocityZ),
		fmt.Sprintf("Acceleration\\ Z=%f", t.Status.Navigator.Kalman.AccelerationZ),

		// Recovery
	}

	return fmt.Sprintf("%s %s", header, strings.Join(fields, ","))
}
