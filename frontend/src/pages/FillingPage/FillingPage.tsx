import {
  DockviewReact,
  type DockviewReadyEvent,
  themeAbyssSpaced,
} from 'dockview-react';

import Navbar from '@/components/Navbar';
import PopoutButton from '@/components/PopoutButton';
import AbortPanel from '@/panels/AbortPanel';
import CommandLogPanel from '@/panels/CommandLogPanel';
import ConfigPanel from '@/panels/ConfigPanel';
import FillingStatePanel from '@/panels/FillingStatePanel';
import FillingStationPanel from '@/panels/FillingStationPanel';
import LogCommandsPanel from '@/panels/LogCommandsPanel';
import TelemetryPanel from '@/panels/TelemetryPanel';

import styles from './FillingPage.module.scss';

export default function FillingPage() {
  const onReady = (event: DockviewReadyEvent) => {
    event.api.addPanel({
      id: 'CONFIG',
      component: 'ConfigPanel',
    });

    event.api.addPanel({
      id: 'panel_2',
      component: 'default',
      position: {
        direction: 'below',
      },
    });

    event.api.addPanel({
      id: 'panel_3',
      component: 'default',
      position: {
        direction: 'below',
      },
    });

    event.api.addPanel({
      id: 'FILLING STATE',
      component: 'FillingStatePanel',
      position: {
        direction: 'right',
      },
    });

    event.api.addPanel({
      id: 'FILLING STATION DIAGRAM',
      component: 'FillingStationPanel',
      position: {
        direction: 'below',
        referencePanel: 'FILLING STATE',
      },
    });

    event.api.addPanel({
      id: 'LOG COMMANDS',
      component: 'LogCommandsPanel',
      position: {
        direction: 'right',
      },
    });

    event.api.addPanel({
      id: 'TELEMETRY',
      component: 'TelemetryPanel',
      position: {
        direction: 'below',
        referencePanel: 'LOG COMMANDS',
      },
    });

    event.api.addPanel({
      id: 'COMMAND LOG',
      component: 'CommandLogPanel',
      position: {
        direction: 'below',
        referencePanel: 'TELEMETRY',
      },
    });

    event.api.addPanel({
      id: 'ABORT',
      component: 'AbortPanel',
      position: {
        direction: 'below',
        referencePanel: 'COMMAND LOG',
      },
    });

    setTimeout(() => {
      const leftGroup = event.api.getPanel('CONFIG')?.group;
      const rightGroup = event.api.getPanel('LOG COMMANDS')?.group;
      const centerTopGroup = event.api.getPanel('FILLING STATE')?.group;

      const containerWidth = event.api.width || window.innerWidth;
      const containerHeight = event.api.height || window.innerHeight;

      if (leftGroup) {
        leftGroup.api.setSize({
          width: containerWidth * 0.23,
          height: containerHeight * 0.248,
        });
      }
      if (rightGroup) {
        rightGroup.api.setSize({ width: containerWidth * 0.25 });
      }
      if (centerTopGroup) {
        centerTopGroup.api.setSize({ height: containerHeight * 0.248 });
      }
    }, 0);
  };

  return (
    <main className={styles.fillingPage}>
      <Navbar />

      <div className={styles.dockviewContainer}>
        <DockviewReact
          theme={themeAbyssSpaced}
          onReady={onReady}
          rightHeaderActionsComponent={PopoutButton}
          components={{
            default: () => <div />,
            ConfigPanel: ConfigPanel,
            FillingStatePanel: FillingStatePanel,
            FillingStationPanel: FillingStationPanel,
            LogCommandsPanel: LogCommandsPanel,
            TelemetryPanel: TelemetryPanel,
            CommandLogPanel: CommandLogPanel,
            AbortPanel: AbortPanel,
          }}
        />
      </div>
    </main>
  );
}
