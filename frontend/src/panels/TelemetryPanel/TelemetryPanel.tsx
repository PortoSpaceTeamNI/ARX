import { Table, TableBody, TableCell, TableRow } from '@/elements/Table';
import { useMissionControl } from '@/utils/store';

import styles from './TelemetryPanel.module.scss';

export default function TelemetryPanel() {
  const telemetry = useMissionControl((state) => state.telemetry);

  return (
    <Table className={styles.telemetryTable}>
      <TableBody>
        <TableRow>
          <TableCell>Packet Loss</TableCell>
          <TableCell className={styles.telemetryValue}>
            {telemetry.packetLoss}
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
  );
}
