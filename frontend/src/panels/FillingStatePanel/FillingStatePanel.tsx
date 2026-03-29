import MissionStateGroup from '@/components/MissionState';
import { MissionState, MissionStateNames } from '@/utils/models/missionstate';
import { useMissionControl } from '@/utils/store';

import styles from './FillingStatePanel.module.scss';

export default function FillingStatePanel() {
  const telemetry = useMissionControl((state) => state.telemetry);

  return (
    <div className={styles.fillingStatePanel}>
      <MissionStateGroup
        stages={[
          MissionStateNames[MissionState.FillingN2],
          MissionStateNames[MissionState.PrePressure],
          MissionStateNames[MissionState.FillingN2O],
          MissionStateNames[MissionState.PostPressure],
        ]}
        currentStage={MissionStateNames[telemetry.status.obc.state]}
      />
    </div>
  );
}
