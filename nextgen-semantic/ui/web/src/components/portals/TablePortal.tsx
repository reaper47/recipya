import {
  Button,
  Dimmer,
  Divider,
  Header,
  Portal,
  Segment,
  Table,
} from 'semantic-ui-react';
import './Portal.css';
import { BasePortalProps } from './props';

interface TableRow {
  cells: TableCell[];
}

type TableCell = string | JSX.Element;

interface TablePortalProps extends BasePortalProps {
  rows: TableRow[];
}

export const TablePortal = ({
  title,
  rows,
  isOpen,
  handleClose,
}: TablePortalProps): JSX.Element => {
  const renderCells = (cells: TableCell[], rowIndex: number): JSX.Element[] =>
    cells.map((cell: TableCell, index: number) => (
      <Table.Cell
        key={`portal-${title}-${rowIndex}${index}`}
        textAlign={index === 0 ? 'right' : 'left'}
      >
        {cell}
      </Table.Cell>
    ));

  const renderTable = (): JSX.Element => {
    const children: JSX.Element[] = rows.map(
      (row: TableRow, rowIndex: number) => (
        <Table.Row key={`portal-${title}-${rowIndex}`}>
          {renderCells(row.cells, rowIndex)}
        </Table.Row>
      )
    );

    return (
      <Table inverted selectable style={{ backgroundColor: 'transparent' }}>
        <Table.Body>{children}</Table.Body>
      </Table>
    );
  };

  return (
    <Portal open={isOpen} onClose={handleClose}>
      <Dimmer
        active={isOpen}
        page
        onClickOutside={handleClose}
        style={{ backgroundColor: 'rgba(0, 0, 0, 0.4)' }}
      >
        <Segment
          style={{
            left: '40%',
            position: 'fixed',
            top: '33%',
            zIndex: 1000,
          }}
          inverted
        >
          <Header>{title}</Header>
          <Divider />
          {renderTable()}
          <Button onClick={handleClose} inverted>
            Close
          </Button>
        </Segment>
      </Dimmer>
    </Portal>
  );
};
