import { create } from 'zustand';

import { LogCommand } from './models/command';
import { MissionState } from './models/missionstate';
import type { Telemetry } from './models/telemetry';
import { type ValveID, ValveState } from './models/valve';
import { sendMessage } from './webSocketManager';

const MAX_LOG_SIZE = 100;

type MissionControlState = {
  log: { timestamp: number; info: string; nRepeated: number }[];
  telemetry: Telemetry;
  updateTelemetry: (newTelemetry: Telemetry) => void;
  openValve: (valve: ValveID) => void;
  closeValve: (valve: ValveID) => void;
  updateSerialPort: (serialPort: string) => void;
  startLocalLog: () => void;
  stopLocalLog: () => void;
};

export const useMissionControl = create<MissionControlState>()((set) => ({
  log: [],
  telemetry: {
    packetLoss: 0,
    latency: 0,
    dataRate: 0,
    commandLog: '',
    availablePorts: [],
    currentPort: "",
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
    set((state) => {
      const safeStatus = newTelemetry.status ?? state.telemetry.status;

      const newCommand = newTelemetry.commandLog;
      let nextLog = [...state.log];

      if (newCommand) {
        const lastEntry = nextLog[nextLog.length - 1];
        if (lastEntry && lastEntry.info === newCommand) {
          nextLog[nextLog.length - 1] = {
            ...lastEntry,
            nRepeated: lastEntry.nRepeated + 1,
          };
        } else {
          nextLog.push({
            timestamp: Date.now(),
            info: newCommand,
            nRepeated: 1,
          });
        }

        if (nextLog.length > MAX_LOG_SIZE) {
          nextLog = nextLog.slice(-MAX_LOG_SIZE);
        }
      }

      return {
        log: nextLog,
        telemetry: {
          ...newTelemetry,
          status: safeStatus,
        },
      };
    }),
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
  updateSerialPort: (serialPort: string) => {
    sendMessage({
      type: 'update_serial_port',
      data: {
        serial_port: serialPort,
      },
    });
  },
  startLocalLog: () => {
    sendMessage({
      type: 'local_log',
      data: {
        command: LogCommand.Start,
      },
    });
  },
  stopLocalLog: () => {
    sendMessage({
      type: 'local_log',
      data: {
        command: LogCommand.Stop,
      },
    });
  },
}));
