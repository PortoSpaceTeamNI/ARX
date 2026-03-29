import { Popover } from '@base-ui/react';

import Button from '@/elements/Button';
import { type ValveID, ValveNames, ValveState } from '@/utils/models/valve';
import { useMissionControl } from '@/utils/store';

import styles from './Valve.module.scss';

type ValveProps = {
  valve: ValveID;
  state: ValveState;
  nameSide?: 'left' | 'right' | 'bottom' | 'top';
  className?: string;
};

export default function Valve({
  valve,
  state,
  nameSide = 'top',
  className,
}: ValveProps) {
  const openValve = useMissionControl((state) => state.openValve);
  const closeValve = useMissionControl((state) => state.closeValve);

  return (
    <Popover.Root>
      <Popover.Trigger
        nativeButton={false}
        openOnHover
        delay={0}
        tabIndex={-1}
        className={`${styles.valve} ${className}`}
        render={<div />}
      >
        <Button
          size="xs"
          variant="outline"
          tabIndex={-1}
          onClick={() => {
            closeValve(valve);
          }}
          className={`${styles.closeButton} ${state === ValveState.Closed ? styles.closed : ''} ${state === ValveState.Closing ? styles.stateTransition : ''}`}
        >
          Close
        </Button>

        <Button
          size="xs"
          variant="outline"
          tabIndex={-1}
          onClick={() => {
            openValve(valve);
          }}
          className={`${styles.openButton} ${state === ValveState.Opened ? styles.opened : ''} ${state === ValveState.Opening ? styles.stateTransition : ''}`}
        >
          Open
        </Button>
      </Popover.Trigger>

      <Popover.Portal>
        <Popover.Positioner side={nameSide} sideOffset={8}>
          <Popover.Popup>
            <div className={styles.valveName}>{ValveNames[valve]}</div>
          </Popover.Popup>
        </Popover.Positioner>
      </Popover.Portal>
    </Popover.Root>
  );
}
