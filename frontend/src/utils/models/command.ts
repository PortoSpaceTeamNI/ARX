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

export type WebCommand = UpdateValveCommand | UpdateSerialPortCommand;
