import type { ValveID, ValveState } from './valve';

export type UpdateValveCommand = {
  type: 'update_valve';
  data: {
    valve: ValveID;
    state: ValveState;
    duration?: number;
  };
};

export type ListSerialPortsCommand = {
  type: 'list_serial_ports';
  data: Record<string, never>;
};

export type ConnectSerialCommand = {
  type: 'connect_serial';
  data: {
    port: string;
    baudRate?: number;
  };
};

export type DisconnectSerialCommand = {
  type: 'disconnect_serial';
  data: Record<string, never>;
};

export type WebCommand =
  | UpdateValveCommand
  | ListSerialPortsCommand
  | ConnectSerialCommand
  | DisconnectSerialCommand;
