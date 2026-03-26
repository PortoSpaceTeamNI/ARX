import { Select as BaseSelect } from '@base-ui/react/select';
import { ChevronDown } from 'lucide-react';
import { useId } from 'react';

import styles from './Select.module.scss';

type SelectProps = {
  label: string;
  items: string[];
  value: string;
  onValueChange: BaseSelect.Root.Props<string>['onValueChange'];
  className?: string;
};

export default function Select({
  label,
  items,
  value,
  onValueChange,
  className,
}: SelectProps) {
  const id = useId();

  return (
    <div className={`${styles.field} ${className}`}>
      <BaseSelect.Root id={id} value={value} onValueChange={onValueChange}>
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
                {items.map((item) => (
                  <BaseSelect.Item
                    key={item}
                    value={item}
                    className={styles.item}
                  >
                    <BaseSelect.ItemText className={styles.text}>
                      {item}
                    </BaseSelect.ItemText>
                  </BaseSelect.Item>
                ))}
              </BaseSelect.List>

              <BaseSelect.ScrollDownArrow />
            </BaseSelect.Popup>
          </BaseSelect.Positioner>
        </BaseSelect.Portal>
      </BaseSelect.Root>
    </div>
  );
}
