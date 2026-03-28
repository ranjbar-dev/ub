import React, { useEffect, useRef, useMemo } from 'react';
import { Subscriber, MessageNames } from 'services/message_service';
import styled from 'styles/styled-components';
import { vw, Translator, CurrencyFormater, toFraction } from 'utils/formatters';
import { CellRenderer } from 'components/renderer';

import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { AgGridReact } from 'ag-grid-react';

import { Button } from '@material-ui/core';
import { Balance } from 'containers/FundsPage/types';
import { GridFilterTypes, GridHeaderNames } from 'containers/App/constants';
import { LocalStorageKeys } from 'services/constants';
import { injectIntl } from 'react-intl';
import NoOpenOrdersIcon from 'images/themedIcons/noOpenOrdersIcon';
import {
  headerHider,
  resizeToTarget,
  whiteShaddowHider,
} from 'utils/gridUtilities';
import { BodyScrollEvent, GridApi } from 'ag-grid-community';

interface GridConfigTypes {
  columnDefs: any[];
  rowData: Balance[];
}
interface Params {
  data: Balance;
  rowIndex: number;
  node: { id: string };
}
const DataGrid = (props: { data: Balance[]; intl: any }) => {
  //@ts-ignore
  const gridApi: { current: GridApi } = useRef();
  const timeout = useRef<any>();
  const intl = useMemo(() => props.intl, []);
  const GridData: Balance[] = props.data;
  //  const [GridData, setGridData] = useState(props.data);
  const viewportWidth = 1552;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };
  const staticRows = useRef([
    {
      headerName: Translate({ message: 'Coin' }),
      field: 'code',
      width: vw(9, viewportWidth),
      minWidth: 95,
      suppressMenu: true,
      sortable: true,
      cellRenderer: (params: Params) =>
        CellRenderer(
          <>
            <div className='coinWrapper'>
              <img src={params.data.image} alt='' />
              <span> {params.data.code}</span>
            </div>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Name' }),
      field: 'name',
      width: vw(9, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'TotalBalance' }),
      field: 'totalAmount',
      width: vw(15, viewportWidth),
      filter: 'agNumberColumnFilter',
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.totalAmount);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'Available' }),
      field: 'availableAmount',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.availableAmount);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'InOrder' }),
      field: 'inOrderAmount',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.inOrderAmount);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'BTCValue' }),
      field: 'btcTotalEquivalentAmount',
      width: vw(25.5, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return `${CurrencyFormater(params.data.btcTotalEquivalentAmount)}
		${
      params.data.btcTotalEquivalentAmount &&
      Number(params.data.btcTotalEquivalentAmount) !== 0
        ? ` ≈ $${toFraction(
            CurrencyFormater(params.data.equivalentTotalAmount),
            2,
          )}`
        : ''
    }`;
      },
    },

    {
      headerName: Translate({ message: 'Actions' }),
      field: 'details',
      width: vw(15, viewportWidth),
      minWidth: 210,
      maxWidth: 210,
      suppressMenu: true,
      cellRenderer: (params: Params) =>
        CellRenderer(
          <>
            <Button
              className='actionButton'
              onClick={() => {
                goToDeposite(params.data.code);
              }}
              color='primary'
            >
              {Translate({ message: 'deposit' })}
            </Button>
            <Button
              onClick={() => {
                goToWithdraw(params.data.code);
              }}
              className='actionButton'
              color='primary'
            >
              {Translate({ message: 'withdraw' })}
            </Button>
          </>,
          'margin-left:-8px;',
        ),
    },
  ]);
  const goToDeposite = (coin: string) => {
    localStorage[LocalStorageKeys.SELECTED_COIN] = coin;
    document.getElementById('depositeTab')!.click();
  };
  const goToWithdraw = (coin: string) => {
    localStorage[LocalStorageKeys.SELECTED_COIN] = coin;
    document.getElementById('withdrawTab')!.click();
  };

  const gridConfig: GridConfigTypes = {
    columnDefs: [...staticRows.current],
    rowData: GridData,
  };
  const gridRendered = e => {
    gridApi.current = e.api;

    const width = window.innerWidth;
    if (width < viewportWidth) {
      gridApi.current.sizeColumnsToFit();
    }
    resizeToTarget({
      elementId: 'balancePageGrid',
      resizeToElementWithId: 'balancePageWrapper',
    });
    headerHider({
      gridApi: gridApi.current,
      hideBottomBorderFromTitledComponent: true,
      wrapperId: 'balancePageWrapper',
    });
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (gridApi.current) {
        if (message.name === MessageNames.RESIZE) {
          resizeToTarget({
            elementId: 'balancePageGrid',
            resizeToElementWithId: 'balancePageWrapper',
          });
          gridApi.current.sizeColumnsToFit();
        }

        if (message.name === MessageNames.SET_BALANCE_PAGE_DATA) {
          // setGridData(message.payload);
          gridApi.current.setRowData(message.payload);
        }
        if (message.name === MessageNames.SET_GRID_FILTER) {
          if (message.payload.showSmallBalances === false) {
            const filterInstance = gridApi.current.getFilterInstance(
              'totalAmount',
            );
            if (filterInstance) {
              filterInstance.setModel({
                type: GridFilterTypes.Greater_Than,
                filter: Number(message.payload.minimum),
                filterTo: null,
              });
            }
            gridApi.current.onFilterChanged();
          } else if (message.payload.showSmallBalances === true) {
            gridApi.current.destroyFilter('totalAmount');
          }
          if (message.payload.searchCoin != undefined) {
            gridApi.current.setQuickFilter(message.payload.searchCoin);

            // let filterInstance = gridApi.current.getFilterInstance('name');
            // if (filterInstance) {
            //   filterInstance.setModel({
            //     type: GridFilterTypes.Contains,
            //     filter: message.payload.searchCoin,
            //   });
            // }
            // gridApi.current.onFilterChanged();
          } else if (
            (message.payload.searchCoin &&
              message.payload.searchCoin.length === 0) ||
            !message.payload.searchCoin
          ) {
            gridApi.current.destroyFilter('name');
          }
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const onScroll = (e: BodyScrollEvent) => {
    if (timeout.current) {
      clearTimeout(timeout.current);
    }
    timeout.current = setTimeout(() => {
      whiteShaddowHider({ gridApi: gridApi.current, isMini: false });
    }, 100);
  };
  return (
    <GridWrapper className='ag-theme-balham' id='balancePageGrid'>
      <AgGridReact
        onGridReady={gridRendered}
        animateRows={true}
        headerHeight={32}
        onBodyScroll={onScroll}
        singleClickEdit={true}
        getRowNodeId={data => data.code}
        immutableData={true}
        rowHeight={40}
        columnDefs={gridConfig.columnDefs}
        rowData={gridConfig.rowData}
        overlayNoRowsTemplate={ReactDOMServer.renderToString(
          <AnimatedNoRows
            icon={<NoOpenOrdersIcon />}
            texts={[
              <span className='black'>
                {Translate({ message: 'noBalance' })}
              </span>,
            ]}
          />,
        )}
      ></AgGridReact>
      {/* {GridData.length === 0 && (
        <Button className={`noRowButton ${Buttons.SimpleRoundButton}`}>
          <FormattedMessage {...translate.Gototrade} />
        </Button>
      )} */}
      <div className='whiteShaddow' id='whiteShaddow'></div>
    </GridWrapper>
  );
};
export default injectIntl(DataGrid);
const GridWrapper = styled.div`
  min-width: 1000px;
  .whiteShaddow {
    position: absolute;
    width: 100%;
    height: 90px;
    left: 0px;
    bottom: 0px;
    pointer-events: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
  }

  .actionButton {
    top: 0px;
    max-height: 35px;
    padding: 0 12px !important;
    min-height: 32px;
    margin: 0 0px;
    left: -5px;
    font-weight: 600;
    span {
      font-size: 13px;
    }
  }

  .coinWrapper {
    img {
      height: 25px;
      width: 25px;
      border: 1px solid #c1c1c1;
      border-radius: 50px;
      padding: 1px;
      min-width: 25px;
      min-height: 25px;
    }
    span {
      margin: 0 8px;
      font-weight: 600;
      font-size: 13px;
    }
  }
`;
