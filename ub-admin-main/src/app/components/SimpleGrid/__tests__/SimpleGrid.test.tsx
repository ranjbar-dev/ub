import React from 'react';
import { renderWithProviders, screen } from 'utils/testUtils';
import { SimpleGrid } from '../SimpleGrid';
import { ColDef } from 'ag-grid-community';
import { createAction } from '@reduxjs/toolkit';

// Mock AG Grid React
jest.mock('ag-grid-react', () => ({
  AgGridReact: (props: Record<string, unknown>) => (
    <div data-testid="ag-grid" data-row-data={JSON.stringify(props.rowData)} />
  ),
}));

jest.mock('app/components/grid_loading/gridLoading', () => ({
  GridLoading: () => <div data-testid="grid-loading">Loading...</div>,
}));

jest.mock('app/components/PaginationComponent/PaginationComponent', () => {
  return function PaginationComponent() {
    return <div data-testid="pagination">Pagination</div>;
  };
});

jest.mock('services/messageService', () => ({
  Subscriber: {
    subscribe: jest.fn(() => ({ unsubscribe: jest.fn() })),
  },
  MessageService: { send: jest.fn() },
  MessageNames: {
    GRID_RESIZE: 'GRID_RESIZE',
    SET_ROW_LOADING: 'SET_ROW_LOADING',
    UPDATE_GRID_ROW: 'UPDATE_GRID_ROW',
    REFRESH_GRID: 'REFRESH_GRID',
    APPLY_PARAMS_TO_GRID: 'APPLY_PARAMS_TO_GRID',
  },
  GridNames: {},
}));

jest.mock('utils/gridUtilities', () => ({
  getPageSize: jest.fn(() => 20),
  resizeToTarget: jest.fn(),
  filterHeight: 50,
  randomColor: jest.fn(() => 'test-color'),
}));

// Create a real Redux action creator for initialAction
const mockInitialAction = createAction<Record<string, unknown>>('test/fetchData');

describe('SimpleGrid', () => {
  const defaultProps = {
    initialAction: mockInitialAction,
    arrayFieldName: 'data',
    staticRows: [
      { field: 'id', headerName: 'ID' },
      { field: 'name', headerName: 'Name' },
    ] as ColDef[],
    immutableId: 'id',
    additionalInitialParams: {},
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders AG Grid component', () => {
    renderWithProviders(<SimpleGrid {...defaultProps} />);
    expect(screen.getByTestId('ag-grid')).toBeInTheDocument();
  });

  it('renders with externalData prop', () => {
    const externalData = {
      data: [
        { id: 1, name: 'Item 1' },
        { id: 2, name: 'Item 2' },
      ],
      count: 2,
    };

    renderWithProviders(
      <SimpleGrid {...defaultProps} externalData={externalData} />,
    );

    expect(screen.getByTestId('ag-grid')).toBeInTheDocument();
  });

  it('renders empty grid when externalData is null', () => {
    renderWithProviders(<SimpleGrid {...defaultProps} externalData={null} />);
    const agGrid = screen.getByTestId('ag-grid');
    const rowData = JSON.parse(agGrid.getAttribute('data-row-data') || '[]');
    expect(rowData).toEqual([]);
  });
});