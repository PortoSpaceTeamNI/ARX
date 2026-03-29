import Button from '@/elements/Button';

import styles from './AbortPanel.module.scss';

export default function AbortPanel() {
  return (
    <div className={styles.abortContainer}>
      <Button variant="destructive" className={styles.abortButton}>
        ABORT
      </Button>
    </div>
  );
}
