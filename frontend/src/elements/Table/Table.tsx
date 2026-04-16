import styles from './Table.module.scss';

type TableComponentProps = {
  className?: string;
  children?: React.ReactNode;
};

export default function Table({
  className = '',
  children,
}: TableComponentProps) {
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

Table.Header = function Header({
  className = '',
  children,
}: TableComponentProps) {
  return <thead className={className}>{children}</thead>;
};

Table.Body = function Body({ className = '', children }: TableComponentProps) {
  return <tbody className={className}>{children}</tbody>;
};

Table.Row = function Row({ className = '', children }: TableComponentProps) {
  return <tr className={`${styles.tableRow} ${className}`}>{children}</tr>;
};

Table.Head = function Head({ className = '', children }: TableComponentProps) {
  return <th className={`${styles.tableHead} ${className}`}>{children}</th>;
};

Table.Cell = function Cell({ className = '', children }: TableComponentProps) {
  return <td className={`${styles.tableCell} ${className}`}>{children}</td>;
};
