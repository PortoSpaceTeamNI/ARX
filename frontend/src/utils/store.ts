import { create } from 'zustand';

import { MissionState } from './models/missionstate';
import type { Telemetry } from './models/telemetry';
import { type ValveID, ValveState } from './models/valve';
import { sendMessage } from './webSocketManager';

type MissionControlState = {
  telemetry: Telemetry;
  updateTelemetry: (newTelemetry: Telemetry) => void;
  openValve: (valve: ValveID) => void;
  closeValve: (valve: ValveID) => void;
};

export const useMissionControl = create<MissionControlState>()((set) => ({
  telemetry: {
    status: {
      obc: {
        state: MissionState.Idle,
        sdLogging: false,
      },
      hydraUf: {
        ventValve: ValveState.Closed,
        pressurizingValve: ValveState.Closed,
        probeThermo1: 0,
        probeThermo2: 0,
        probeThermo3: 0,
      },
      hydraLf: {
        mainValve: ValveState.Closed,
        abortValve: ValveState.Closed,
        probeThermo1: 0,
        probeThermo2: 0,
        chamberTemperature: 0,
        tankPressure: 0,
        chamberPressure: 0,
      },
      hydraFs: {
        n2: {
          fillValve: ValveState.Closed,
          purgeValve: ValveState.Closed,
          pressure: 0,
        },
        n2o: {
          fillValve: ValveState.Closed,
          purgeValve: ValveState.Closed,
          pressure: 0,
          temperature: 0,
        },
        quickDc: {
          n2Valve: ValveState.Closed,
          n2oValve: ValveState.Closed,
          pressure: 0,
        },
      },
      navigator: {
        gps: {
          latitude: 0,
          longitude: 0,
          satCount: 0,
          altitude: 0,
          horizontalVelocity: 0,
        },
        kalman: {
          velocityZ: 0,
          accelerationZ: 0,
          altitude: 0,
          maxAltitude: 0,
          q1: 0,
          q2: 0,
          q3: 0,
          q4: 0,
        },
        imu: {
          accelX: 0,
          accelY: 0,
          accelZ: 0,
          gyroX: 0,
          gyroY: 0,
          gyroZ: 0,
          magX: 0,
          magY: 0,
          magZ: 0,
        },
        barometerAlt1: 0,
        barometerAlt2: 0,
      },
      elytra: {
        mainChuteEmatch: false,
        drogueChuteEmatch: false,
      },
      liftR: {
        loadcell1: 0,
        loadcell2: 0,
        loadcell3: 0,
        mainEmatch: false,
      },
      liftFs: {
        n2oLoadcell: 0,
      },
    },
  },
  updateTelemetry: (newTelemetry: Telemetry) =>
    set({ telemetry: newTelemetry }),
  openValve: (valve: ValveID, duration?: number) => {
    sendMessage({
      type: 'update_valve',
      data: {
        valve: valve,
        state: ValveState.Opened,
        duration: duration ? duration : undefined,
      },
    });
  },
  closeValve: (valve: ValveID, duration?: number) => {
    sendMessage({
      type: 'update_valve',
      data: {
        valve: valve,
        state: ValveState.Closed,
        duration: duration ? duration : undefined,
      },
    });
  },
}));
