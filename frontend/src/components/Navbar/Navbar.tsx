import { Plus, Settings } from 'lucide-react';

import Button from '@/elements/Button';

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
          <Button variant="ghost" size="iconLg">
            <Plus />
          </Button>
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
