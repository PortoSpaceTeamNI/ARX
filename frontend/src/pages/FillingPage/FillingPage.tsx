import { useEffect, useMemo, useState } from 'react';
import { Pause, Play, Square } from 'lucide-react';
import { Group, Panel } from 'react-resizable-panels';

import fillingStationDiagram from '@/assets/fillingStationDiagram.png';
import Button from '@/components/Button';
import MissionStageGroup from '@/components/MissionStage';
import { Table, TableBody, TableCell, TableRow } from '@/components/Table';
import Valve from '@/components/Valve';
import { MissionState, MissionStateNames } from '@/utils/models/missionstate';
import { ValveID } from '@/utils/models/valve';
import { useMissionControl } from '@/utils/store';

import styles from './FillingPage.module.scss';

export default function FillingPage() {
  const telemetry = useMissionControl((state) => state.telemetry);
  const serialStatus = useMissionControl((state) => state.serialStatus);
  const requestSerialPorts = useMissionControl((state) => state.requestSerialPorts);
  const connectSerial = useMissionControl((state) => state.connectSerial);
  const disconnectSerial = useMissionControl((state) => state.disconnectSerial);
  const [selectedPort, setSelectedPort] = useState('');

  const canConnect = useMemo(
    () => selectedPort !== '' && (!serialStatus.connected || serialStatus.currentPort !== selectedPort),
    [selectedPort, serialStatus.connected, serialStatus.currentPort]
  );

  useEffect(() => {
    if (serialStatus.currentPort) {
      setSelectedPort(serialStatus.currentPort);
      return;
    }

    if (serialStatus.ports.length > 0 && !serialStatus.ports.includes(selectedPort)) {
      setSelectedPort(serialStatus.ports[0]);
      return;
    }

    if (serialStatus.ports.length === 0 && selectedPort !== '') {
      setSelectedPort('');
    }
  }, [serialStatus.currentPort, serialStatus.ports, selectedPort]);

  return (
    <main className={styles.fillingPage}>
      <Group orientation="vertical" className={styles.group}>
        <Panel defaultSize="20%">
          <Group orientation="horizontal" className={styles.group}>
            <Panel defaultSize="20%" className={styles.panel}>
              <p className={styles.panelName}>SERIAL</p>
              <div className={styles.serialControls}>
                <select
                  className={styles.serialSelect}
                  value={selectedPort}
                  onChange={(event) => setSelectedPort(event.target.value)}
                >
                  <option value="" disabled>
                    {serialStatus.ports.length > 0 ? 'Select port' : 'No ports found'}
                  </option>
                  {serialStatus.ports.map((port) => (
                    <option key={port} value={port}>
                      {port}
                    </option>
                  ))}
                </select>

                <div className={styles.serialButtons}>
                  <Button variant="outline" size="xs" onClick={requestSerialPorts}>
                    Refresh
                  </Button>

                  <Button
                    variant="outline"
                    size="xs"
                    disabled={!canConnect}
                    onClick={() => connectSerial(selectedPort)}
                  >
                    Connect
                  </Button>

                  <Button
                    variant="outline"
                    size="xs"
                    disabled={!serialStatus.connected}
                    onClick={disconnectSerial}
                  >
                    Disconnect
                  </Button>
                </div>

                <p className={styles.serialStatusText}>
                  {serialStatus.connected && serialStatus.currentPort
                    ? `Connected: ${serialStatus.currentPort}`
                    : 'Disconnected'}
                </p>

                {serialStatus.error && (
                  <p className={styles.serialErrorText}>{serialStatus.error}</p>
                )}
              </div>
            </Panel>

            <Panel className={styles.panel}>
              <MissionStageGroup
                stages={[
                  MissionStateNames[MissionState.FillingN2],
                  MissionStateNames[MissionState.PrePressure],
                  MissionStateNames[MissionState.FillingN2O],
                  MissionStateNames[MissionState.PostPressure],
                ]}
                currentStage={MissionStateNames[telemetry.status.obc.state]}
              />
            </Panel>

            <Panel
              defaultSize="20%"
              className={`${styles.panel} ${styles.fillingStageControls}`}
            >
              <Button variant="ghost" size="iconLg">
                <Pause />
              </Button>
              <Button variant="ghost" size="iconLg">
                <Square />
              </Button>
              <Button variant="ghost" size="iconLg">
                <Play />
              </Button>
            </Panel>
          </Group>
        </Panel>

        <Panel>
          <Group orientation="horizontal" className={styles.group}>
            <Panel defaultSize="20%">
              <Group orientation="vertical" className={styles.group}>
                <Panel className={styles.panel}>Vent Camera</Panel>
                <Panel className={styles.panel}>To be Determined</Panel>
              </Group>
            </Panel>

            <Panel
              className={`${styles.panel} ${styles.fillingStageContainer}`}
            >
              <img
                src={fillingStationDiagram}
                alt="Diagram of the filling station"
              />

              <Valve
                valve={ValveID.Pressurizing}
                state={telemetry.status.hydraUf.pressurizingValve}
                className={styles.pressurizingValve}
              />
              <Valve
                valve={ValveID.Vent}
                state={telemetry.status.hydraUf.ventValve}
                className={styles.ventValve}
                nameSide="bottom"
              />
              <Valve
                valve={ValveID.Abort}
                state={telemetry.status.hydraLf.abortValve}
                className={styles.abortValve}
              />
              <Valve
                valve={ValveID.Main}
                state={telemetry.status.hydraLf.mainValve}
                className={styles.mainValve}
              />
              <Valve
                valve={ValveID.N2OFill}
                state={telemetry.status.hydraFs.n2o.fillValve}
                className={styles.n2oFillValve}
              />
              <Valve
                valve={ValveID.N2OPurge}
                state={telemetry.status.hydraFs.n2o.purgeValve}
                className={styles.n2oPurgeValve}
                nameSide="bottom"
              />
              <Valve
                valve={ValveID.N2Fill}
                state={telemetry.status.hydraFs.n2.fillValve}
                className={styles.n2FillValve}
              />
              <Valve
                valve={ValveID.N2Purge}
                state={telemetry.status.hydraFs.n2.purgeValve}
                className={styles.n2PurgeValve}
              />
              <Valve
                valve={ValveID.N2OQuickDc}
                state={telemetry.status.hydraFs.quickDc.n2oValve}
                className={styles.n2oQuickDcValve}
              />
              <Valve
                valve={ValveID.N2QuickDc}
                state={telemetry.status.hydraFs.quickDc.n2Valve}
                className={styles.n2QuickDcValve}
              />
            </Panel>

            <Panel defaultSize="20%">
              <Group orientation="vertical" className={styles.group}>
                <Panel className={styles.panel}>
                  <p className={styles.panelName}>TELEMETRY</p>
                  <Table>
                    <TableBody>
                      <TableRow>
                        <TableCell>Packet Loss</TableCell>
                        <TableCell className={styles.telemetryValue}>
                          {0}
                        </TableCell>
                        <TableCell className={styles.telemetryUnit}>pkts</TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell>Data Rate</TableCell>
                        <TableCell className={styles.telemetryValue}>
                          {telemetry.dataRate}
                        </TableCell>
                        <TableCell className={styles.telemetryUnit}>B/s</TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell>Latency</TableCell>
                        <TableCell className={styles.telemetryValue}>
                          {telemetry.latency}
                        </TableCell>
                        <TableCell className={styles.telemetryUnit}>ms</TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </Panel>

                <Panel className={styles.abortContainer}>
                  <Button variant="destructive" className={styles.abortButton}>
                    ABORT
                  </Button>
                </Panel>
              </Group>
            </Panel>
          </Group>
        </Panel>
      </Group>
    </main>
  );
}
