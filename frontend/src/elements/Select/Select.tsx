import { Select as BaseSelect } from '@base-ui/react/select';
import { ChevronDown } from 'lucide-react';
import { type ReactNode, useEffect, useId, useState } from 'react';

import styles from './Select.module.scss';

type SelectProps = {
  label: string;
  children: ReactNode;
  value: string;
  onValueChange: BaseSelect.Root.Props<string>['onValueChange'];
  className?: string;
};

export default function Select({
  label,
  children,
  value,
  onValueChange,
  className,
}: SelectProps) {
  const id = useId();
  // Fallback to null to guarantee the component is initialized as controlled!
  const [localValue, setLocalValue] = useState<string | null>(value ?? null);

  useEffect(() => {
    setLocalValue(value ?? null);
  }, [value]);

  const handleValueChange: BaseSelect.Root.Props<string>['onValueChange'] = (
    newValue,
    event
  ) => {
    if (newValue !== null) {
      setLocalValue(newValue);
    }
    if (onValueChange) {
      onValueChange(newValue, event);
    }
  };

  return (
    <div className={`${styles.field} ${className}`}>
      <BaseSelect.Root
        id={id}
        value={localValue}
        onValueChange={handleValueChange}
      >
        <label htmlFor={id} className={styles.label}>
          {label}
        </label>

        <BaseSelect.Trigger className={styles.select}>
          <BaseSelect.Value
            placeholder={`Select ${label}`}
            className={styles.value}
          />

          <BaseSelect.Icon className={styles.icon}>
            <ChevronDown />
          </BaseSelect.Icon>
        </BaseSelect.Trigger>

        <BaseSelect.Portal>
          <BaseSelect.Positioner
            alignItemWithTrigger={false}
            className={styles.positioner}
            sideOffset={4}
          >
            <BaseSelect.Popup className={styles.popup}>
              <BaseSelect.ScrollUpArrow />

              <BaseSelect.List className={styles.list}>
                {children}
              </BaseSelect.List>

              <BaseSelect.ScrollDownArrow />
            </BaseSelect.Popup>
          </BaseSelect.Positioner>
        </BaseSelect.Portal>
      </BaseSelect.Root>
    </div>
  );
}

type SelectItemProps = {
  value: string;
  children: ReactNode;
  className?: string;
};

Select.Item = function SelectItem({
  value,
  children,
  className,
}: SelectItemProps) {
  return (
    <BaseSelect.Item
      value={value}
      className={`${styles.item} ${className || ''}`.trim()}
    >
      {children}
    </BaseSelect.Item>
  );
};
