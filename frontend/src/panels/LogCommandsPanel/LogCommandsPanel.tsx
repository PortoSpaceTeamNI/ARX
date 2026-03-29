import { Play, Square } from 'lucide-react';

import Button from '@/elements/Button';
import { useMissionControl } from '@/utils/store';

import styles from './LogCommandsPanel.module.scss';

export default function LogCommandsPanel() {
  const startLocalLog = useMissionControl((state) => state.startLocalLog);
  const stopLocalLog = useMissionControl((state) => state.stopLocalLog);

  return (
    <div className={styles.logCommandsPanel}>
      <Button variant="ghost" size="iconLg" onClick={startLocalLog}>
        <Play />
      </Button>
      <Button variant="ghost" size="iconLg" onClick={stopLocalLog}>
        <Square />
      </Button>
    </div>
  );
}
