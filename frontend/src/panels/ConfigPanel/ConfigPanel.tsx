import { Search } from 'lucide-react';

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
          items={telemetry.availablePorts}
          value={telemetry.currentPort}
          onValueChange={(val) => val && updateSerialPort(val)}
          className={styles.select}
        />
        <Button
          variant="ghost"
          size="icon"
          onClick={() => {
            /*updateSerialPort("search") TODO: This*/
          }}
        >
          <Search />
        </Button>
      </li>
    </ul>
  );
}
