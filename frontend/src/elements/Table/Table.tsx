import styles from './Table.module.scss';

type TableComponentProps = {
  className?: string;
  children?: React.ReactNode;
};

export function Table({ className = '', children }: TableComponentProps) {
  return (
    <table
      cellSpacing={0}
      cellPadding={0}
      className={`${styles.table} ${className}`}
    >
      {children}
    </table>
  );
}

export function TableHeader({ className = '', children }: TableComponentProps) {
  return <thead className={className}>{children}</thead>;
}

export function TableBody({ className = '', children }: TableComponentProps) {
  return <tbody className={className}>{children}</tbody>;
}

export function TableRow({ className = '', children }: TableComponentProps) {
  return <tr className={`${styles.tableRow} ${className}`}>{children}</tr>;
}

export function TableHead({ className = '', children }: TableComponentProps) {
  return <th className={`${styles.tableHead} ${className}`}>{children}</th>;
}

export function TableCell({ className = '', children }: TableComponentProps) {
  return <td className={`${styles.tableCell} ${className}`}>{children}</td>;
}
