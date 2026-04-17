import { LayoutDashboard, Plus, RefreshCw, Save, Settings } from 'lucide-react';

import Button from '@/elements/Button';
import Menu from '@/elements/Menu';

import styles from './Navbar.module.scss';

export default function Navbar() {
  return (
    <nav className={styles.navbar}>
      <ul>
        <li>
          <Button variant="link" to="/filling">
            FILLING
          </Button>
        </li>
        <li>
          <Button variant="link" to="/launch">
            LAUNCH
          </Button>
        </li>
        <li>
          <Button variant="link" to="/recovery">
            RECOVERY
          </Button>
        </li>
      </ul>

      <ul>
        <li>
          <Menu>
            <Menu.Trigger>
              <Button variant="ghost" size="iconLg">
                <LayoutDashboard />
              </Button>
            </Menu.Trigger>
            <Menu.Items>
              <Menu.Item>
                <Button variant="outline" className={styles.stretched}>
                  <RefreshCw /> Refresh Default Layout
                </Button>
              </Menu.Item>
              <Menu.Item>
                <Button variant="outline" className={styles.stretched}>
                  <Save /> Save Current Layout
                </Button>
              </Menu.Item>
              <Menu.Item>
                <Button variant="outline" className={styles.stretched}>
                  <Plus /> Add Panel
                </Button>
              </Menu.Item>
            </Menu.Items>
          </Menu>
        </li>
        <li>
          <Button variant="ghost" size="iconLg">
            <Settings />
          </Button>
        </li>
      </ul>
    </nav>
  );
}
