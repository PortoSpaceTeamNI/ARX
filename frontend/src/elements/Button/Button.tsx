import { Button as BaseButton } from '@base-ui/react';

import styles from './Button.module.scss';

type ButtonProps = {
  variant?: 'default' | 'outline' | 'ghost' | 'destructive';
  size?: 'sizeDefault' | 'xs' | 'icon' | 'iconLg';
  children: React.ReactNode;
} & BaseButton.Props;

export default function Button({
  variant = 'default',
  size = 'sizeDefault',
  children,
  className,
  ...props
}: ButtonProps) {
  return (
    <BaseButton
      className={`${styles.button} ${className} ${styles[variant]} ${styles[size]}`}
      {...props}
    >
      {children}
    </BaseButton>
  );
}
