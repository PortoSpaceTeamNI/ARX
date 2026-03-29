import { useLayoutEffect, useRef } from 'react';

import { useMissionControl } from '@/utils/store';

import styles from './CommandLogPanel.module.scss';

export default function CommandLogPanel() {
  const log = useMissionControl((state) => state.log);

  const logContainerRef = useRef<HTMLDivElement>(null);

  useLayoutEffect(() => {
    if (logContainerRef.current && log.length > 0) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    }
  }, [log]);

  return (
    <div className={styles.commandLogPanel}>
      <div className={styles.logs} ref={logContainerRef}>
        {log.map((entry) => (
          <div key={entry.timestamp} className={styles.log}>
            <p className={styles.info}>{entry.info}</p>
            {entry.nRepeated > 1 && <div>{entry.nRepeated}</div>}
          </div>
        ))}
      </div>
    </div>
  );
}
