import { Circle, Search } from 'lucide-react';

import Button from '@/elements/Button';
import Select from '@/elements/Select';
import { useMissionControl } from '@/utils/store';

import styles from './ConfigPanel.module.scss';

export default function ConfigPanel() {
  const telemetry = useMissionControl((state) => state.telemetry);
  const updateSerialPort = useMissionControl((state) => state.updateSerialPort);

  return (
    <ul className={styles.configPanel}>
      <li className={styles.serialPort}>
        <Select
          label="Serial Port"
          value={telemetry.currentPort.port}
          onValueChange={(val) => val && updateSerialPort(val)}
          className={styles.serialField}
        >
          {telemetry.availablePorts.map((port) => (
            <Select.Item
              key={port.port}
              value={port.port}
              className={styles.item}
            >
              {port.port}
              <Circle
                className={port.state ? styles.connected : styles.disconnected}
              />
            </Select.Item>
          ))}
        </Select>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => {
            updateSerialPort('find');
          }}
        >
          <Search />
        </Button>
      </li>
    </ul>
  );
}
