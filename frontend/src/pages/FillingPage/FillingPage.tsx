import {
  type DockviewApi,
  DockviewReact,
  type DockviewReadyEvent,
  themeAbyssSpaced,
} from 'dockview-react';
import { useState } from 'react';

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

export const loadFillingDefaultLayout = (dockviewApi: DockviewApi) => {
  dockviewApi.clear();

  dockviewApi.addPanel({
    id: 'CONFIG',
    component: 'ConfigPanel',
  });

  dockviewApi.addPanel({
    id: 'panel_2',
    component: 'default',
    position: {
      direction: 'below',
    },
  });

  dockviewApi.addPanel({
    id: 'panel_3',
    component: 'default',
    position: {
      direction: 'below',
    },
  });

  dockviewApi.addPanel({
    id: 'FILLING STATE',
    component: 'FillingStatePanel',
    position: {
      direction: 'right',
    },
  });

  dockviewApi.addPanel({
    id: 'FILLING STATION DIAGRAM',
    component: 'FillingStationPanel',
    position: {
      direction: 'below',
      referencePanel: 'FILLING STATE',
    },
  });

  dockviewApi.addPanel({
    id: 'LOG COMMANDS',
    component: 'LogCommandsPanel',
    position: {
      direction: 'right',
    },
  });

  dockviewApi.addPanel({
    id: 'TELEMETRY',
    component: 'TelemetryPanel',
    position: {
      direction: 'below',
      referencePanel: 'LOG COMMANDS',
    },
  });

  dockviewApi.addPanel({
    id: 'COMMAND LOG',
    component: 'CommandLogPanel',
    position: {
      direction: 'below',
      referencePanel: 'TELEMETRY',
    },
  });

  dockviewApi.addPanel({
    id: 'ABORT',
    component: 'AbortPanel',
    position: {
      direction: 'below',
      referencePanel: 'COMMAND LOG',
    },
  });

  setTimeout(() => {
    const leftGroup = dockviewApi.getPanel('CONFIG')?.group;
    const rightGroup = dockviewApi.getPanel('LOG COMMANDS')?.group;
    const centerTopGroup = dockviewApi.getPanel('FILLING STATE')?.group;

    const containerWidth = dockviewApi.width || window.innerWidth;
    const containerHeight = dockviewApi.height || window.innerHeight;

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

export default function FillingPage() {
  const [api, setApi] = useState<DockviewApi | null>(null);

  const onReady = (event: DockviewReadyEvent) => {
    setApi(event.api);
    loadFillingDefaultLayout(event.api);
  };

  return (
    <main className={styles.fillingPage}>
      <Navbar dockviewApi={api} />

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
