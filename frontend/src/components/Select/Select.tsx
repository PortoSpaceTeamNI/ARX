import { Select as BaseSelect } from '@base-ui/react/select';
import { ChevronDown } from 'lucide-react';
import { useId } from 'react';

import styles from './Select.module.scss';

type SelectProps = {
  label: string;
  items: { label: string; value: string }[];
  onValueChange: BaseSelect.Root.Props<string>['onValueChange'];
};

export default function Select({ label, items, onValueChange }: SelectProps) {
  const id = useId();

  return (
    <div className={styles.field}>
      <BaseSelect.Root id={id} onValueChange={onValueChange}>
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
                {items.map(({ label, value }) => (
                  <BaseSelect.Item
                    key={label}
                    value={value}
                    className={styles.item}
                  >
                    <BaseSelect.ItemText className={styles.text}>
                      {label}
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
