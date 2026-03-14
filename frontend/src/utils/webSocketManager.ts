import type { WebCommand } from './models/command';
import type { Telemetry } from './models/telemetry';
import { useMissionControl } from './store';

const RECONNECT_BASE_DELAY = 1000;
const RECONNECT_MAX_DELAY = 30000;
const MAX_RETRIES = 5;

let ws: WebSocket | null = null;
let reconnectAttempts = 0;

export function initWebSocket() {
  if (ws?.readyState === WebSocket.OPEN) {
    return;
  }

  const host = window.location.host;
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(`${protocol}//${host}/ws`);

  ws.onopen = () => {
    console.log('connected');
    reconnectAttempts = 0;
  };

  ws.onmessage = (event) => {
    const telemetry: Telemetry = JSON.parse(event.data);
    console.log(telemetry);
    useMissionControl.getState().updateTelemetry(telemetry);
  };

  ws.onclose = () => {
    if (reconnectAttempts < MAX_RETRIES) {
      const delay = Math.min(
        RECONNECT_BASE_DELAY * 2 ** reconnectAttempts,
        RECONNECT_MAX_DELAY
      );
      reconnectAttempts++;
      console.log(`Connection lost. Retrying in ${delay}ms...`);

      setTimeout(initWebSocket, delay);
    }
  };

  ws.onerror = () => {
    ws?.close();
  };
}

export function sendMessage(message: WebCommand) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
    return true;
  }

  console.warn('SendMessage failed: WebSocket is not open.');
  return false;
}
