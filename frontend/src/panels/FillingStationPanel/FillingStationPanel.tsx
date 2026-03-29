import fillingStationDiagram from '@/assets/fillingStationDiagram.png';
import Valve from '@/components/Valve';
import { ValveID } from '@/utils/models/valve';
import { useMissionControl } from '@/utils/store';

import styles from './FillingStationPanel.module.scss';

export default function FillingStationPanel() {
  const telemetry = useMissionControl((state) => state.telemetry);

  return (
    <div className={styles.fillingStationPanel}>
      <img src={fillingStationDiagram} alt="Diagram of the filling station" />

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
    </div>
  );
}
