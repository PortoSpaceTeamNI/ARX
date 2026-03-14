package command

import (
	"missioncontrol/backend/models/missionstate"
	"missioncontrol/backend/models/packet"
	"missioncontrol/backend/models/packet/commandid"
	"missioncontrol/backend/models/packet/communicatorid"
	"missioncontrol/backend/models/packet/payload"
	"missioncontrol/backend/models/valve"
	"missioncontrol/backend/models/valve/valvestate"
)

type ParsedStatusResponse struct {
	OBC       OBC       `json:"obc"`
	HydraUF   HydraUF   `json:"hydraUf"`
	HydraLF   HydraLF   `json:"hydraLf"`
	HydraFS   HydraFS   `json:"hydraFs"`
	Navigator Navigator `json:"navigator"`
	LiftR     LiftR     `json:"liftR"`
	LiftFS    LiftFS    `json:"liftFs"`
	Elytra    Elytra    `json:"elytra"`
}

type OBC struct {
	State     missionstate.MissionState `json:"state"`
	SdLogging bool                      `json:"sdLogging"`
}

type HydraUF struct {
	PressurizingValve valvestate.ValveState `json:"pressurizingValve"`
	VentValve         valvestate.ValveState `json:"ventValve"`
	ProbeThermo1      float64               `json:"probeThermo1"`
	ProbeThermo2      float64               `json:"probeThermo2"`
	ProbeThermo3      float64               `json:"probeThermo3"`
}

type HydraLF struct {
	AbortValve         valvestate.ValveState `json:"abortValve"`
	MainValve          valvestate.ValveState `json:"mainValve"`
	ProbeThermo1       float64               `json:"probeThermo1"`
	ProbeThermo2       float64               `json:"probeThermo2"`
	ChamberTemperature float64               `json:"chamberTemperature"`
	TankPressure       float64               `json:"tankPressure"`
	ChamberPressure    float64               `json:"chamberPressure"`
}

type HydraFS struct {
	N2      N2      `json:"n2"`
	N2O     N2O     `json:"n2o"`
	QuickDC QuickDC `json:"quickDc"`
}

type N2 struct {
	FillValve  valvestate.ValveState `json:"fillValve"`
	PurgeValve valvestate.ValveState `json:"purgeValve"`
	Pressure   float64               `json:"pressure"`
}

type N2O struct {
	FillValve   valvestate.ValveState `json:"fillValve"`
	PurgeValve  valvestate.ValveState `json:"purgeValve"`
	Pressure    float64               `json:"pressure"`
	Temperature float64               `json:"temperature"`
}

type QuickDC struct {
	N2Valve  valvestate.ValveState `json:"n2Valve"`
	N2OValve valvestate.ValveState `json:"n2oValve"`
	Pressure float64               `json:"pressure"`
}

type Navigator struct {
	GPS           GPS     `json:"gps"`
	Kalman        Kalman  `json:"kalman"`
	IMU           IMU     `json:"imu"`
	BarometerAlt1 float64 `json:"barometerAlt1"`
	BarometerAlt2 float64 `json:"barometerAlt2"`
}

type GPS struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Altitude           float64 `json:"altitude"`
	HorizontalVelocity float64 `json:"horizontalVelocity"`
	SatCount           uint8   `json:"satCount"`
}

type Kalman struct {
	VelocityZ     float64 `json:"velocityZ"`
	AccelerationZ float64 `json:"accelerationZ"`
	Altitude      float64 `json:"altitude"`
	MaxAltitude   float64 `json:"maxAltitude"`
	Q1            float64 `json:"q1"`
	Q2            float64 `json:"q2"`
	Q3            float64 `json:"q3"`
	Q4            float64 `json:"q4"`
}

type IMU struct {
	AccelX float64 `json:"accelX"`
	AccelY float64 `json:"accelY"`
	AccelZ float64 `json:"accelZ"`
	GyroX  float64 `json:"gyroX"`
	GyroY  float64 `json:"gyroY"`
	GyroZ  float64 `json:"gyroZ"`
	MagX   float64 `json:"magX"`
	MagY   float64 `json:"magY"`
	MagZ   float64 `json:"magZ"`
}

type LiftR struct {
	MainEmatch bool    `json:"mainEmatch"`
	Loadcell1  float64 `json:"loadcell1"`
	Loadcell2  float64 `json:"loadcell2"`
	Loadcell3  float64 `json:"loadcell3"`
}

type LiftFS struct {
	N2OLoadcell float64 `json:"n2oLoadcell"`
}

type Elytra struct {
	DrogueChuteEmatch bool `json:"drogueChuteEmatch"`
	MainChuteEmatch   bool `json:"mainChuteEmatch"`
}

func (sr *ParsedStatusResponse) GetValveState(v valve.Valve) valvestate.ValveState {
	switch v {
	case valve.Pressurizing:
		return sr.HydraUF.PressurizingValve

	case valve.Vent:
		return sr.HydraUF.VentValve

	case valve.Abort:
		return sr.HydraLF.AbortValve

	case valve.Main:
		return sr.HydraLF.MainValve

	case valve.N2OFill:
		return sr.HydraFS.N2O.FillValve

	case valve.N2OPurge:
		return sr.HydraFS.N2O.PurgeValve

	case valve.N2Fill:
		return sr.HydraFS.N2.FillValve

	case valve.N2Purge:
		return sr.HydraFS.N2.PurgeValve

	case valve.N2OQuickDc:
		return sr.HydraFS.QuickDC.N2OValve

	case valve.N2QuickDc:
		return sr.HydraFS.QuickDC.N2Valve
	}

	panic("Tried getting state of undefined valve")
}

type StatusCommand struct{}

func (c *StatusCommand) ToString() string {
	return "Status"
}

func (c *StatusCommand) ToPacket() packet.Packet {
	return packet.NewEmptyPacket(communicatorid.OBC, commandid.Status)
}

func (c *StatusCommand) ParseResponse(raw packet.Packet) (IResponse, error) {
	p := raw.Payload.(*payload.StatusResponsePayload)

	isSet := func(flag uint16) bool {
		return (p.BitFlags & flag) != 0
	}

	checkValveState := func(flag uint16) valvestate.ValveState {
		if isSet(flag) {
			return valvestate.Opened
		}

		return valvestate.Closed
	}

	scale := func(val int16) float64 {
		return float64(val) / 100.0
	}

	return ParsedStatusResponse{
		OBC: OBC{
			State:     missionstate.MissionState(p.MissionState),
			SdLogging: false,
		},
		HydraUF: HydraUF{
			PressurizingValve: checkValveState(payload.HydraUFPressurizingValveState),
			VentValve:         checkValveState(payload.HydraUFVentValveState),
			ProbeThermo1:      scale(p.HydraUFProbeThermo1),
			ProbeThermo2:      scale(p.HydraUFProbeThermo2),
			ProbeThermo3:      scale(p.HydraUFProbeThermo3),
		},
		HydraLF: HydraLF{
			AbortValve:         checkValveState(payload.HydraLFAbortValveState),
			MainValve:          checkValveState(payload.HydraLFMainValveState),
			ProbeThermo1:       scale(p.HydraLFProbeThermo1),
			ProbeThermo2:       scale(p.HydraLFProbeThermo2),
			ChamberTemperature: scale(p.HydraLFChamberTemperature),
			TankPressure:       scale(p.HydraLFTankPressure),
			ChamberPressure:    scale(p.HydraLFChamberPressure),
		},
		HydraFS: HydraFS{
			N2: N2{
				FillValve:  checkValveState(payload.HydraFSN2FillValveState),
				PurgeValve: checkValveState(payload.HydraFSN2PurgeValveState),
				Pressure:   scale(p.HydraFSN2Pressure),
			},
			N2O: N2O{
				FillValve:   checkValveState(payload.HydraFSN2OFillValveState),
				PurgeValve:  checkValveState(payload.HydraFSN2OPurgeValveState),
				Pressure:    scale(p.HydraFSN2OPressure),
				Temperature: scale(p.HydraFSN2OTemperature),
			},
			QuickDC: QuickDC{
				N2Valve:  checkValveState(payload.HydraFSN2QuickDCValveState),
				N2OValve: checkValveState(payload.HydraFSN2OQuickDCValveState),
				Pressure: scale(p.HydraFSQuickDCPressure),
			},
		},
		Navigator: Navigator{
			GPS:           GPS{Latitude: 0, Longitude: 0, Altitude: 0, HorizontalVelocity: 0, SatCount: 0},
			Kalman:        Kalman{VelocityZ: 0, AccelerationZ: 0, Altitude: 0, MaxAltitude: 0, Q1: 0, Q2: 0, Q3: 0, Q4: 0},
			IMU:           IMU{AccelX: 0, AccelY: 0, AccelZ: 0, GyroX: 0, GyroY: 0, GyroZ: 0, MagX: 0, MagY: 0, MagZ: 0},
			BarometerAlt1: 0, BarometerAlt2: 0,
		},
		LiftR: LiftR{
			MainEmatch: isSet(payload.LiftRMainEmatch),
			Loadcell1:  scale(p.LiftRLoadCell1),
			Loadcell2:  scale(p.LiftRLoadCell2),
			Loadcell3:  scale(p.LiftRLoadCell3),
		},
		LiftFS: LiftFS{
			N2OLoadcell: scale(p.LiftFSN2OLoadCell),
		},
		Elytra: Elytra{
			DrogueChuteEmatch: isSet(payload.ElytraDrogueParachuteEmatch),
			MainChuteEmatch:   isSet(payload.ElytraMainParachuteEmatch),
		},
	}, nil
}
