import type { DockviewApi } from 'dockview-react';
import { LayoutDashboard, Plus, RefreshCw, Save, Settings } from 'lucide-react';
import { useLocation } from 'react-router-dom';

import Button from '@/elements/Button';
import Menu from '@/elements/Menu';
import { loadFillingDefaultLayout } from '@/pages/FillingPage/FillingPage';

import styles from './Navbar.module.scss';

type NavbarProps = {
  dockviewApi: DockviewApi | null;
};

export default function Navbar({ dockviewApi }: NavbarProps) {
  const location = useLocation();

  const handleLoadDefaultLayout = () => {
    if (!dockviewApi) return;

    if (location.pathname === '/filling') {
      loadFillingDefaultLayout(dockviewApi);
    }

    if (location.pathname === '/launch') {
      //loadLaunchDefaultLayout(dockviewApi);
    }

    if (location.pathname === '/recovery') {
      //loadRecoveryDefaultLayout(dockviewApi);
    }
  };

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
                <Button
                  variant="outline"
                  className={styles.stretched}
                  onClick={handleLoadDefaultLayout}
                >
                  <RefreshCw /> Load Default Layout
                </Button>
              </Menu.Item>
              <Menu.Item>
                <Button variant="outline" className={styles.stretched}>
                  <Save /> Save Current Layout
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
