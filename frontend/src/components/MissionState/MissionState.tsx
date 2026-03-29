import React from 'react';

import styles from './MissionState.module.scss';

type MissionStageGroupProps = {
  currentStage: string;
  stages: string[];
};

export default function MissionStageGroup({
  currentStage,
  stages,
}: MissionStageGroupProps) {
  return (
    <div className={styles.missionStateGroup}>
      {stages.map((stage, index) => (
        <React.Fragment key={stage}>
          <div className={styles.stateSlot}>
            <div
              className={`${styles.state} ${index < stages.indexOf(currentStage) ? styles.complete : ''} ${index === stages.indexOf(currentStage) ? styles.current : ''}`}
            >
              {stage}
            </div>
          </div>
          {index < stages.length - 1 && (
            <div
              className={`${styles.stateSeparator} ${index < stages.indexOf(currentStage) ? styles.complete : ''}`}
            />
          )}
        </React.Fragment>
      ))}
    </div>
  );
}
