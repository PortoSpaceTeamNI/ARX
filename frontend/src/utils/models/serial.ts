export type SerialStatus = {
  ports: string[];
  connected: boolean;
  currentPort?: string;
  baudRate?: number;
  error?: string;
};
