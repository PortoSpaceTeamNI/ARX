import { Button as BaseButton } from '@base-ui/react';
import { useLocation } from 'react-router-dom';

import styles from './Button.module.scss';

type ButtonProps = {
  variant?: 'default' | 'outline' | 'ghost' | 'destructive' | 'link';
  size?: 'sizeDefault' | 'xs' | 'icon' | 'iconLg';
  to?: string; // variant link only prop
  children: React.ReactNode;
} & BaseButton.Props;

export default function Button({
  variant = 'default',
  size = 'sizeDefault',
  to,
  children,
  className,
  ...props
}: ButtonProps) {
  const location = useLocation();

  return (
    <BaseButton
      className={`${styles.button} ${className} ${styles[variant]} ${variant === 'link' && location.pathname === to ? styles.active : ''} ${styles[size]}`}
      {...props}
      render={variant === 'link' ? <a href={to}>{children}</a> : undefined}
      role={variant === 'link' ? 'link' : undefined}
    >
      {variant !== 'link' && children}
    </BaseButton>
  );
}
