import { Circle } from 'lucide-react';

import Table from '@/elements/Table';
import { useMissionControl } from '@/utils/store';

import styles from './TelemetryPanel.module.scss';

export default function TelemetryPanel() {
  const telemetry = useMissionControl((state) => state.telemetry);

  return (
    <Table className={styles.telemetryTable}>
      <Table.Body>
        <Table.Row>
          <Table.Cell>Connection Status</Table.Cell>
          <Table.Cell className={styles.telemetryValue}>
            {telemetry.currentPort.state ? 'On' : 'Off'}
          </Table.Cell>
          <Table.Cell className={styles.telemetryUnit}>
            <Circle
              className={
                telemetry.currentPort.state
                  ? styles.connected
                  : styles.disconnected
              }
            />
          </Table.Cell>
        </Table.Row>
        <Table.Row>
          <Table.Cell>Packet Loss</Table.Cell>
          <Table.Cell className={styles.telemetryValue}>
            {telemetry.packetLoss}
          </Table.Cell>
          <Table.Cell className={styles.telemetryUnit}>pkts</Table.Cell>
        </Table.Row>
        <Table.Row>
          <Table.Cell>Data Rate</Table.Cell>
          <Table.Cell className={styles.telemetryValue}>
            {telemetry.dataRate}
          </Table.Cell>
          <Table.Cell className={styles.telemetryUnit}>B/s</Table.Cell>
        </Table.Row>
        <Table.Row>
          <Table.Cell>Latency</Table.Cell>
          <Table.Cell className={styles.telemetryValue}>
            {telemetry.latency}
          </Table.Cell>
          <Table.Cell className={styles.telemetryUnit}>ms</Table.Cell>
        </Table.Row>
      </Table.Body>
    </Table>
  );
}
