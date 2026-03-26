import type { MissionState } from './missionstate';
import type { ValveState } from './valve';

export type Telemetry = {
  packetLoss: number;
  latency: number;
  dataRate: number;
  status: {
    obc: OBC;
    hydraUf: HydraUF;
    hydraLf: HydraLF;
    hydraFs: HydraFS;
    navigator: Navigator;
    elytra: Elytra;
    liftR: LiftR;
    liftFs: LiftFS;
  };
  commandLog: string;
  availablePorts: string[];
  currentPort: string;
};

type OBC = {
  state: MissionState;
  sdLogging: boolean;
};

type HydraUF = {
  ventValve: ValveState;
  pressurizingValve: ValveState;
  probeThermo1: number;
  probeThermo2: number;
  probeThermo3: number;
};

type HydraLF = {
  mainValve: ValveState;
  abortValve: ValveState;
  probeThermo1: number;
  probeThermo2: number;
  chamberTemperature: number;
  tankPressure: number;
  chamberPressure: number;
};

type HydraFS = {
  n2: N2;
  n2o: N2O;
  quickDc: QuickDC;
};

type N2 = {
  fillValve: ValveState;
  purgeValve: ValveState;
  pressure: number;
};

type N2O = {
  fillValve: ValveState;
  purgeValve: ValveState;
  pressure: number;
  temperature: number;
};

type QuickDC = {
  n2Valve: ValveState;
  n2oValve: ValveState;
  pressure: number;
};

type Navigator = {
  gps: GPS;
  kalman: Kalman;
  imu: IMU;
  barometerAlt1: number;
  barometerAlt2: number;
};

type GPS = {
  latitude: number;
  longitude: number;
  satCount: number;
  altitude: number;
  horizontalVelocity: number;
};

type Kalman = {
  velocityZ: number;
  accelerationZ: number;
  altitude: number;
  maxAltitude: number;
  q1: number;
  q2: number;
  q3: number;
  q4: number;
};

type IMU = {
  accelX: number;
  accelY: number;
  accelZ: number;
  gyroX: number;
  gyroY: number;
  gyroZ: number;
  magX: number;
  magY: number;
  magZ: number;
};

type Elytra = {
  mainChuteEmatch: boolean;
  drogueChuteEmatch: boolean;
};

type LiftR = {
  loadcell1: number;
  loadcell2: number;
  loadcell3: number;
  mainEmatch: boolean;
};

type LiftFS = {
  n2oLoadcell: number;
};
