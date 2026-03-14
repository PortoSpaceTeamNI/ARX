import React from 'react';

import styles from './MissionStage.module.scss';

type MissionStageGroupProps = {
  currentStage: string;
  stages: string[];
};

export default function MissionStageGroup({
  currentStage,
  stages,
}: MissionStageGroupProps) {
  return (
    <div className={styles.missionStageGroup}>
      {stages.map((stage, index) => (
        <React.Fragment key={stage}>
          <div className={styles.stageSlot}>
            <div
              className={`${styles.stage} ${index < stages.indexOf(currentStage) ? styles.complete : ''} ${index === stages.indexOf(currentStage) ? styles.current : ''}`}
            >
              {stage}
            </div>
          </div>
          {index < stages.length - 1 && (
            <div
              className={`${styles.stageSeparator} ${index < stages.indexOf(currentStage) ? styles.complete : ''}`}
            />
          )}
        </React.Fragment>
      ))}
    </div>
  );
}
