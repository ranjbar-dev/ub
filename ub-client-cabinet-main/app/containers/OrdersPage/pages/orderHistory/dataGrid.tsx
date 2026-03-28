import React, {
  useEffect,
  useState,
  useMemo,
  useRef,
  useCallback,
} from 'react';
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
  formatCurrencyWithMaxFraction,
  formatTableCell,
} from 'utils/formatters';
import { CellRenderer } from 'components/renderer';

import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { AgGridReact } from 'ag-grid-react';

import { Button } from '@material-ui/core';
import { Buttons, GridHeaderNames, AppPages } from 'containers/App/constants';
import { injectIntl } from 'react-intl';
import {
  DetailsStyle,
  ToggleDetail,
  RowWithShaddow,
  headerHider,
  resizeToTarget,
  whiteShaddowHider,
} from 'utils/gridUtilities';
import { useDispatch } from 'react-redux';
import {
  getOrderDetailAction,
  getOrderHistoryAction,
  getPaginatedOrderHistoryAction,
} from 'containers/OrdersPage/actions';

import NoOpenOrdersIcon from 'images/themedIcons/noOpenOrdersIcon';
import { DWStatus } from 'containers/FundsPage/constants';
import HistoryDetailRow from 'containers/OrdersPage/components/detailRow/HistoryDetailRow';
import { FilterApplier } from 'components/gridFilters/filterApply';
import ExpandMore from 'images/themedIcons/expandMore';
import { GridReadyEvent } from 'ag-grid-community';
import { push } from 'redux-first-history';
import RxLoader from 'components/RXLoader/RxLoader';
import { OrderPages } from 'containers/OrdersPage/constants';
import { storage } from 'utils/storage';
import { LocalStorageKeys } from 'services/constants';
import { currencyMap } from 'utils/sharedData';
import { miniGridStyles } from 'containers/OrdersPage/styles';

const DataGrid = (props: { data: Order[]; intl: any; isMini?: boolean }) => {
  const gridApi: any = useRef();
  const { data, intl, isMini } = props;

  if (data?.length > 0 && !isMini) {
    data.forEach(item => {
      item.isDetailsOpen = false;
      delete item?.details;
    });
  }

  const [IsGridLoading, setIsGridLoading] = useState(false);
  const dispatch = useDispatch();
  const viewportWidth = 1552;
  const Translate = useCallback((data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  }, []);
  const ordersData = useRef(data);
  const filtersRef: any = useRef({});
  const canGetNewData = useRef(true);
  const lastScrollPosition = useRef(200);
  const lockDataGet = useRef(false);

  const gridColumnsRef = useRef([
    {
      headerName: Translate({ message: 'Date' }),
      field: 'createdAt',
      width: vw(10, viewportWidth),
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
      width: vw(8, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => PairFormat(data.pair),
    },
    {
      headerName: Translate({ message: 'Type' }),
      field: 'type',
      width: vw(6, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: function ({ data, rowIndex }) {
        return {
          textTransform: 'capitalize',
        };
      },
    },
    {
      headerName: Translate({ message: 'Side' }),
      field: 'side',
      width: vw(6, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: function ({ data, rowIndex }) {
        return {
          textTransform: ' capitalize',
          color: data.side === 'buy' ? 'var(--greenText)' : 'var(--redText)',
        };
      },
    },
    {
      headerName: Translate({ message: 'Price' }),
      field: 'price',
      width: vw(10, viewportWidth),
      valueFormatter: (params: any) => {
        return params.data.price
          ? CurrencyFormater(params.data.price)
          : params.data.type === 'market'
          ? CurrencyFormater(params.data.averagePrice)
          : Translate({ message: 'market' });
      },
      comparator: function (a, b) {
        return a - b;
      },
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Amount' }),
      field: 'amount',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => formatTableCell({ value: data.amount }),
      comparator: function (a, b) {
        return +a.split(' ')[0] - +b.split(' ')[0];
      },
    },
    {
      headerName: Translate({ message: 'FilledPercent' }),
      field: 'executed',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,

      valueFormatter: ({ data }) => {
        if (data.executed && data.executed.includes('%')) {
          const tmp = data.executed.split('%')[0];
          if (Number(tmp) > 99.9999) {
            data.executed = '100%';
          }
        }
        return data.executed && data.executed.length > 0 ? data.executed : '-';
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
    {
      headerName: Translate({ message: 'TriggerConditions' }),
      field: 'triggerCondition',
      width: vw(12, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: { paddingLeft: '45px' },
      valueFormatter: ({ data }) => {
        return data.triggerCondition && data.triggerCondition !== ''
          ? data.triggerCondition
          : '-';
      },
    },
    {
      headerName: Translate({ message: 'status' }),
      field: 'status',
      width: vw(11, viewportWidth),
      suppressMenu: true,
      sortable: true,
      maxWidth: isMini ? 80 : 300,
      minWidth: 80,
      cellRenderer: ({ data, rowIndex }) =>
        CellRenderer(
          <>
            <div
              className={`upperFirst statusBadge ${data.status} ${
                isMini ? 'mini' : ''
              }`}
            >
              {data.status}
            </div>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Info' }),
      field: 'detail',
      width: vw(5, viewportWidth),
      minWidth: 70,
      hide: isMini ? true : false,
      suppressMenu: true,
      getQuickFilterText: function (params) {
        const data: Order = params.data;
        return data.status == 'canceled' ? 'showAllOrders' : 'showNotCancelled';
      },
      cellRenderer: params =>
        CellRenderer(
          <>
            {params.data.details && (
              <div>
                <HistoryDetailRow
                  details={params.data.details}
                  mainData={params.data}
                />
              </div>
            )}
            {params.data.status === DWStatus.Completed ||
            params.data.status === DWStatus.CONFIRMED ? (
              <Button
                onClick={() => {
                  ToggleDetail({
                    isMini,
                    gridApi: gridApi.current,
                    params,
                    isOpen: params.data.isDetailsOpen,
                  });
                  if (!params.data.details) {
                    dispatch(
                      getOrderDetailAction({
                        order_id: params.data.id,
                        rowId: params.node.id,
                      }),
                    );
                  }
                }}
                // disableRipple
                className={`detailButton black  ${
                  params.data.isDetailsOpen === true ? 'grey' : ''
                }`}
                id={'dtButton' + params.data.id}
                color='primary'
                endIcon={
                  <ExpandMore
                    className={
                      params.data.isDetailsOpen === true
                        ? 'rotateAnimate200 rotated'
                        : 'rotateAnimate200'
                    }
                    id={'expandIcon' + params.data.id}
                  />
                }
              >
                {Translate({ message: 'Detail' })}
              </Button>
            ) : params.data.status === DWStatus.Expired ? (
              <div></div>
            ) : (
              <div style={{ marginLeft: '-3px' }}>
                {Translate({ message: 'Detail' })}
                <span style={{ marginLeft: '4px' }}>{' -'}</span>
              </div>
            )}
          </>,
        ),
    },
  ]);

  useEffect(() => {
    const EventSubscription = EventSubscriber.subscribe((message: any) => {
      if (
        (message.name === EventMessageNames.REFRESH_ORDER_GRID &&
          message.id !== OrderPages.OPEN_ORDER) ||
        message.name === MessageNames.RECONNECT_EVENT
      ) {
        dispatch(getOrderHistoryAction({ silent: true }));
      }
    });
    return () => {
      EventSubscription.unsubscribe();
    };
  }, []);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (gridApi.current) {
        if (message.name === MessageNames.RESIZE) {
          setTimeout(() => {
            resizeToTarget({
              elementId: 'OrderHistoryGridWrapper',
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
                elementId: 'OrderHistoryGridWrapper',
                resizeToElementWithId: 'ORDERSContainer',
              });
              gridApi.current.sizeColumnsToFit();
            });
          }
        }
        if (message.name === MessageNames.SET_ORDER_HISTORY_DATA) {
          gridApi.current.setRowData(message.payload);
          ordersData.current = message.payload;
          canGetNewData.current = true;
          headerHider({
            gridApi: gridApi.current,
            gridId: 'OrderHistoryGridWrapper',
          });
        }
        if (
          message.name === MessageNames.SET_PAGE_FILTERS_WITH_ID &&
          message.id === 'orderHistory'
        ) {
          filtersRef.current = message.payload;
          lockDataGet.current = false;
          canGetNewData.current = true;
          lastScrollPosition.current = 200;
        }
        if (message.name === MessageNames.SET_PAGINATED_ORDER_HISTORY_DATA) {
          if (message.payload.length > 0) {
            ordersData.current = [...ordersData.current, ...message.payload];
            gridApi.current.setRowData(ordersData.current);
            canGetNewData.current = true;
            lockDataGet.current = false;
          } else {
            lockDataGet.current = true;
          }
        }

        if (message.name === MessageNames.SET_ORDER_DETAIL) {
          gridApi.current.forEachNode((node, index) => {
            if (node.id === message.rowId) {
              node.data.details = message.payload;
              node.data.isDetailsOpen = true;
              gridApi.current.applyTransaction({
                update: [node.data],
              });
              gridApi.current.redrawRows({ rowNodes: [node] });
              return;
            }
          });
        }

        if (message.name === MessageNames.SET_GRID_FILTER) {
          FilterApplier({ gridApi: gridApi.current, message: message });
        }
        if (message.name === MessageNames.SETLOADING) {
          setIsGridLoading(message.payload);
        }
        if (message.name === MessageNames.ADD_ONE_ORDER_TO_HISTORY) {
          gridApi.current.applyTransaction({
            add: [message.payload],
            addIndex: 0,
          });
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [IsGridLoading]);

  const gridConfig: any = {
    columnDefs: [...gridColumnsRef.current],
    rowData: ordersData.current,
  };
  const gridRendered = (e: GridReadyEvent) => {
    gridApi.current = e.api;
    const width = window.innerWidth;
    if (width < viewportWidth || isMini) {
      gridApi.current.sizeColumnsToFit();
      if (isMini) {
        resizeToTarget({
          elementId: 'OrderHistoryGridWrapper',
          resizeToElementWithId: 'ORDERSContainer',
        });
      }
    }
    headerHider({
      gridApi: gridApi.current,
      gridId: 'OrderHistoryGridWrapper',
    });
  };
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
        getPaginatedOrderHistoryAction({
          last_id: ordersData.current[ordersData.current.length - 1].id,
          ...filtersRef.current,
        }),
      );
    }
    whiteShaddowHider({ gridApi: gridApi.current, isMini });
  };
  return useMemo(
    () => (
      <>
        <RxLoader style={{ left: '38vw', top: '4vh' }} id='orderHistory' />
        <GridWrapper
          id='OrderHistoryGridWrapper'
          className={`ag-theme-balham withExpandableRows ${
            isMini ? 'miniGrid' : ''
          }`}
        >
          <AgGridReact
            onGridReady={gridRendered}
            animateRows={true}
            headerHeight={isMini ? 24 : 32}
            rowHeight={isMini ? 24 : 40}
            immutableData={true}
            onBodyScroll={onScroll}
            getRowNodeId={(data: Order) => {
              return data.id.toString();
            }}
            getRowStyle={RowWithShaddow}
            columnDefs={gridConfig.columnDefs}
            rowData={gridConfig.rowData}
            overlayNoRowsTemplate={ReactDOMServer.renderToString(
              <AnimatedNoRows
                isMini={isMini}
                icon={<NoOpenOrdersIcon />}
                texts={[
                  <span className='black'>
                    {Translate({ message: 'noOrder' })}
                  </span>,
                ]}
              />,
            )}
          ></AgGridReact>
          {ordersData.current.length === 0 && !isMini && (
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
    ),
    [IsGridLoading],
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
    z-index: 1;
    pointer-events: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
  }

  .noRowButton {
    bottom: calc(40vh - 160px) !important;
  }
  .upperFirst {
    text-transform: capitalize;
    font-size: 13px;
  }
  ${DetailsStyle}
`;
