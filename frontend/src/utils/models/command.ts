import type { ValveID, ValveState } from './valve';

export type UpdateValveCommand = {
  type: 'update_valve';
  data: {
    valve: ValveID;
    state: ValveState;
    duration?: number;
  };
};

export type WebCommand = UpdateValveCommand;
