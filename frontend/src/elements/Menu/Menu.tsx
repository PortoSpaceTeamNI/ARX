import { Menu as BaseMenu } from '@base-ui/react/menu';
import type { ReactNode } from 'react';

import styles from './Menu.module.scss';

type MenuProps = {
  children: ReactNode;
};

function Menu({ children }: MenuProps) {
  return <BaseMenu.Root>{children}</BaseMenu.Root>;
}

Menu.Trigger = function MenuTrigger({ children }: { children: ReactNode }) {
  return (
    <BaseMenu.Trigger tabIndex={-1} className={styles.trigger}>
      {children}
    </BaseMenu.Trigger>
  );
};

Menu.Items = function MenuItems({ children }: { children: ReactNode }) {
  return (
    <BaseMenu.Portal>
      <BaseMenu.Positioner sideOffset={8}>
        <BaseMenu.Popup className={styles.popup}>{children}</BaseMenu.Popup>
      </BaseMenu.Positioner>
    </BaseMenu.Portal>
  );
};

Menu.Item = function MenuItem({ children }: { children: ReactNode }) {
  return <BaseMenu.Item tabIndex={-1}>{children}</BaseMenu.Item>;
};

export default Menu;
