import type { ValveID, ValveState } from './valve';

export type UpdateValveCommand = {
  type: 'update_valve';
  data: {
    valve: ValveID;
    state: ValveState;
    duration?: number;
  };
};

export type UpdateSerialPortCommand = {
  type: 'update_serial_port';
  data: {
    serial_port: string;
  };
};

export const LogCommand = {
  Start: 0,
  Stop: 1,
} as const;

export type LogCommand = (typeof LogCommand)[keyof typeof LogCommand];

export type LocalLogCommand = {
  type: 'local_log';
  data: {
    command: LogCommand;
  };
};

export type WebCommand =
  | UpdateValveCommand
  | UpdateSerialPortCommand
  | LocalLogCommand;
