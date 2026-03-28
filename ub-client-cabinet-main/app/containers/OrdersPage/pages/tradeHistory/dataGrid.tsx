import React, { useEffect, useState, useRef } from 'react';
import { Order } from 'containers/OrdersPage/types';
import {
  Subscriber,
  MessageNames,
  EventSubscriber,
  EventMessageNames,
} from 'services/message_service';
import styled from 'styles/styled-components';
import {
  vw,
  Translator,
  CurrencyFormater,
  PairFormat,
  formatTableCell,
} from 'utils/formatters';

import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { AgGridReact } from 'ag-grid-react';

import { Button } from '@material-ui/core';
import { Buttons, GridHeaderNames, AppPages } from 'containers/App/constants';
import { injectIntl } from 'react-intl';
import NoOpenOrdersIcon from 'images/themedIcons/noOpenOrdersIcon';
import { FilterApplier } from 'components/gridFilters/filterApply';
import { ColDef } from 'ag-grid-community';
import {
  headerHider,
  resizeToTarget,
  whiteShaddowHider,
} from 'utils/gridUtilities';
import { useDispatch } from 'react-redux';
import { push } from 'redux-first-history';
import {
  getPaginatedTradeHistoryAction,
  getTradeHistoryAction,
} from 'containers/OrdersPage/actions';
import RxLoader from 'components/RXLoader/RxLoader';
import { miniGridStyles } from 'containers/OrdersPage/styles';

//let gridApi: any, gridColumnApi: any;
interface GridConfigTypes {
  columnDefs: any[];
  rowData: Order[];
  bottomData: any;
}
const DataGrid = (props: { data: Order[]; intl: any; isMini?: boolean }) => {
  //   let gridData: OpenOrder[] = props.data;
  const gridApi: any = useRef();
  const intl = props.intl;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };

  const ordersData = useRef(props.data);
  const filtersRef: any = useRef({});
  const canGetNewData = useRef(true);
  const lastScrollPosition = useRef(200);
  const lockDataGet = useRef(false);

  const [IsGridLoading, setIsGridLoading] = useState(false);
  const viewportWidth = 1552;
  const dispatch = useDispatch();

  useEffect(() => {
    const EventSubscription = EventSubscriber.subscribe((message: any) => {
      if (
        message.name === EventMessageNames.REFRESH_ORDER_GRID ||
        message.name === MessageNames.RECONNECT_EVENT
      ) {
        //@ts-ignore
        dispatch(getTradeHistoryAction({ silent: true }));
      }
    });
    return () => {
      EventSubscription.unsubscribe();
    };
  }, []);

  const staticRows: ColDef[] = [
    {
      headerName: Translate({ message: 'Date' }),
      field: 'createdAt',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      minWidth: 145,
      filter: 'agNumberColumnFilter',
      filterValueGetter: params => {
        return params.data.createdAtToFilter;
      },
    },

    {
      headerName: Translate({ message: 'Pair' }),
      field: 'pair',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => PairFormat(data.pair),
    },
    {
      headerName: Translate({ message: 'Side' }),
      field: 'type',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: function ({ data, rowIndex }) {
        return {
          textTransform: ' capitalize',
          color: data.type === 'buy' ? 'var(--greenText)' : 'var(--redText)',
        };
      },
    },
    {
      headerName: Translate({ message: 'Price' }),
      field: 'price',
      width: vw(15, viewportWidth),
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.price);
      },
      comparator: function (a, b) {
        return +a.split(' ')[0] - +b.split(' ')[0];
      },
      suppressMenu: true,
      sortable: true,
    },

    {
      headerName: Translate({ message: 'Filled' }),
      field: 'executed',
      width: vw(14, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => formatTableCell({ value: data.executed }),
      comparator: function (a, b) {
        return +a.split(' ')[0] - +b.split(' ')[0];
      },
    },
    {
      headerName: Translate({ message: 'Fee' }),
      field: 'fee',
      width: vw(14, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => formatTableCell({ value: data.fee }),
      comparator: function (a, b) {
        return +a.split(' ')[0] - +b.split(' ')[0];
      },
    },
    {
      headerName: Translate({ message: 'Total' }),
      field: 'total',
      width: vw(12, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => formatTableCell({ value: data.total }),
      comparator: function (a, b) {
        return +a.split(' ')[0] - +b.split(' ')[0];
      },
    },
  ];
  const gridConfig: GridConfigTypes = {
    columnDefs: [...staticRows],
    rowData: ordersData.current,
    bottomData: { athlete: 'Total' },
  };
  const gridRendered = e => {
    gridApi.current = e.api;
    //gridColumnApi = e.columnApi;
    const width = window.innerWidth;
    if (width < viewportWidth || props.isMini) {
      gridApi.current.sizeColumnsToFit();
      if (props.isMini) {
        resizeToTarget({
          elementId: 'TradeHistoryGridWrapper',
          resizeToElementWithId: 'ORDERSContainer',
        });
      }
    }
    headerHider({
      gridApi: gridApi.current,
      gridId: 'TradeHistoryGridWrapper',
    });
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (gridApi.current) {
        if (message.name === MessageNames.RESIZE) {
          setTimeout(() => {
            resizeToTarget({
              elementId: 'TradeHistoryGridWrapper',
              resizeToElementWithId: 'ORDERSContainer',
            });
            setTimeout(() => {
              gridApi.current.sizeColumnsToFit();
            }, 0);
          }, 170);
        }
        if (message.name === MessageNames.LAYOUT_RESIZE) {
          if (message.payload && message.payload.i === 'ORDERS') {
            //setTimeout(() => {

            //}, 0);
            requestAnimationFrame(() => {
              resizeToTarget({
                elementId: 'TradeHistoryGridWrapper',
                resizeToElementWithId: 'ORDERSContainer',
              });
              gridApi.current.sizeColumnsToFit();
            });
          }
        }
        if (message.name === MessageNames.SET_GRID_FILTER) {
          FilterApplier({ gridApi: gridApi.current, message: message });
        }
        if (message.name === MessageNames.SET_TRADE_HISTORY_DATA) {
          setIsGridLoading(false);
          ordersData.current = message.payload;
          canGetNewData.current = true;
          gridApi.current.setRowData(message.payload);
          headerHider({
            gridApi: gridApi.current,
            gridId: 'TradeHistoryGridWrapper',
          });
        }
        if (
          message.name === MessageNames.SET_PAGE_FILTERS_WITH_ID &&
          message.id === 'tradeHistory'
        ) {
          filtersRef.current = message.payload;
          lockDataGet.current = false;
          canGetNewData.current = true;
          lastScrollPosition.current = 200;
        }
        if (message.name === MessageNames.SET_PAGINATED_TRADE_HISTORY_DATA) {
          if (message.payload.length > 0) {
            ordersData.current = [...ordersData.current, ...message.payload];
            gridApi.current.setRowData(ordersData.current);
            canGetNewData.current = true;
            lockDataGet.current = false;
          } else {
            lockDataGet.current = true;
          }
        }

        if (message.name === MessageNames.SETLOADING) {
          setIsGridLoading(message.payload);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const onScroll = e => {
    const lastRow = gridApi.current.getLastDisplayedRow();
    if (
      lastRow >= ordersData.current.length - 5 &&
      e.top > lastScrollPosition.current + 100 &&
      lockDataGet.current === false &&
      canGetNewData.current === true
    ) {
      lastScrollPosition.current = e.top;
      canGetNewData.current = false;
      dispatch(
        getPaginatedTradeHistoryAction({
          last_id: ordersData.current[ordersData.current.length - 1].id,
          ...filtersRef.current,
        }),
      );
    }
    whiteShaddowHider({ gridApi: gridApi.current, isMini: props.isMini });
  };
  return (
    <>
      <RxLoader style={{ left: '38vw', top: '4vh' }} id='tradeHistory' />
      <GridWrapper
        id='TradeHistoryGridWrapper'
        className={`ag-theme-balham ${props.isMini ? 'miniGrid' : ''}`}
      >
        <AgGridReact
          onGridReady={gridRendered}
          animateRows={true}
          onBodyScroll={onScroll}
          headerHeight={props.isMini ? 24 : 32}
          rowHeight={props.isMini ? 24 : 40}
          columnDefs={gridConfig.columnDefs}
          immutableData={true}
          getRowNodeId={(data: Order) => {
            return data.id.toString();
          }}
          rowData={gridConfig.rowData}
          overlayNoRowsTemplate={ReactDOMServer.renderToString(
            <AnimatedNoRows
              isMini={props.isMini}
              icon={<NoOpenOrdersIcon />}
              texts={[
                <span className='black'>
                  {Translate({ message: 'noOrder' })}
                </span>,
              ]}
            />,
          )}
        ></AgGridReact>
        {ordersData.current.length === 0 && !props.isMini && (
          <div className='noRowsButtonWrapper'>
            <Button
              onClick={() => {
                dispatch(push(AppPages.TradePage));
              }}
              className={`noRowButton ${Buttons.SimpleRoundButton}`}
            >
              {Translate({ message: 'GoToTrade' })}
            </Button>
          </div>
        )}
        <div id='whiteShaddow' className='whiteShaddow'></div>
      </GridWrapper>
    </>
  );
};
export default injectIntl(DataGrid);
const GridWrapper = styled.div`
  height: calc(90vh - 170px);
  ${miniGridStyles}
  .whiteShaddow {
    position: absolute;
    width: 100%;
    height: 90px;

    left: 0px;
    bottom: -10px;
    pointer-events: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
  }

  .cancelButton {
    position: absolute;
    top: 2px;
    max-height: 35px;
    padding: 0 !important;
    min-height: 32px;
    left: 0px;
  }
  .noRowsButtonWrapper .noRowButton {
    bottom: calc(40vh - 157px);
  }
`;
