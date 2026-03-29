import React, { useEffect, useRef } from 'react';
import { injectIntl } from 'react-intl';
import styled from 'styles/styled-components';
import { AgGridReact } from 'ag-grid-react';
import { ColDef, CellClickedEvent } from 'ag-grid-community';
import { GridHeaderNames, MqttTopicsPrefixes } from 'containers/App/constants';
import { Translator, vw } from 'utils/formatters';

import {
  SideSubscriber,
  MessageNames,
  Subscriber,
  MessageService,
  OrderBookSubscriber,
} from 'services/message_service';
import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';
import {
  ResizeGridHeigth,
  DrawDepthChart,
} from 'containers/TradePage/utils/tradeUtilities';
import { CellRenderer } from 'components/renderer';
import PrecisionSelect from 'containers/TradePage/components/orderBook/precisionSelect';
import { MqttService } from 'services/MqttService2';
import { LayoutContainers } from '../layout/layout';
import { savedPairName } from 'utils/sharedData';

const viewportWidth = 500;
const drawDepth = data => {
  requestAnimationFrame(() => {
    DrawDepthChart({
      buyData: data.buy,
      sellData: data.sell,
    });
  });
};

const OrderBookGrid = (props: {
  intl: any;
  subject: string;
  uniqueId: string;
  enabled: boolean;
}) => {
  let data: any = {};
  const arr: any = useRef([]);
  let Precision = 8;
  let sellMax = 0;
  let buyMax = 0;

  const mqtt2 = useRef(MqttService.getInstance());

  const selectedPair = useRef(savedPairName());
  const gridApi: any = useRef();

  const intl = props.intl;

  //////////////////////resize useEffect
  useEffect(() => {
    if (props.enabled === true) {
      drawDepth({ sell: [], buy: [] });
      mqtt2.current.ConnectToSubject({
        subject: `${MqttTopicsPrefixes.OrderBookAddress}${selectedPair.current}`,
      });
    }
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESIZE) {
        setTimeout(() => {
          ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -11 });
          if (data.sell) {
            drawDepth(data);
          }
          setTimeout(() => {
            gridApi.current.sizeColumnsToFit();
          }, 0);
        }, 255);
      }
      if (message.name === MessageNames.LAYOUT_RESIZE) {
        if (message.payload.i === LayoutContainers.ORDERBOOK) {
          requestAnimationFrame(() => {
            ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -11 });
            if (data.sell) {
              drawDepth(data);
            }
            gridApi.current.sizeColumnsToFit();
          });
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  ///////////////end of resize useEffect
  ////////mqtt

  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };

  useEffect(() => {
    let currected: { buy: any[]; sell: any[] } = {
      buy: [],
      sell: [],
    };
    const BSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        requestAnimationFrame(() => {
          mqtt2.current.ConnectToNewSubject({
            oldsubject:
              MqttTopicsPrefixes.OrderBookAddress + selectedPair.current,
            newSubject: `${MqttTopicsPrefixes.OrderBookAddress}${message.payload.name}`,
          });
          selectedPair.current = message.payload.name;
          clearChart();
        });
      }
    });
    const OrderBookSubscription = OrderBookSubscriber.subscribe(
      (message: any) => {
        currected = {
          buy: [],
          sell: [],
        };
        if (message.payload.asks && !document.hidden) {
          const { asks, bids } = message.payload;
          for (let i = 0; i < asks.length; i++) {
            if (i === 0) {
              currected.buy.push({ ...bids[i], price: (+bids[i].price).toFixed(Precision) });
              currected.sell.push({ ...asks[i], price: (+asks[i].price).toFixed(Precision) });
            } else {
              if (
                bids[i] &&
                (+bids[i].price).toFixed(Precision) !==
                  (+bids[i - 1].price).toFixed(Precision)
              ) {
                currected.buy.push({ ...bids[i], price: Number(bids[i].price).toFixed(Precision) });
                currected.sell.push({ ...asks[i], price: Number(asks[i].price).toFixed(Precision) });
              }
            }
          }
        }
        if (message.payload.asks && !document.hidden) {
          drawDepth(currected);
          buyMax = sellMax = 0;
          data = {
            ...currected,
            // sell: currected.sell,
            buy: currected.buy.reverse(),
          };

          for (let i = 0; i < data.sell.length; i++) {
            arr.current[i] = {
              buy: data.buy[i],
              sell: data.sell[i],
              id: i,
            };
          }
          if (gridApi.current) {
            requestAnimationFrame(() => {
              buyMax = data.buy[data.buy.length - 1].sum;
              sellMax = data.sell[data.sell.length - 1].sum;
              gridApi.current.setRowData(arr.current);
            });
          }
        }
      },
    );
    return () => {
      Precision = 6;
      BSubscription.unsubscribe();
      OrderBookSubscription.unsubscribe();
      mqtt2.current &&
        mqtt2.current.DisconnectFromSubject({
          subject: MqttTopicsPrefixes.OrderBookAddress + selectedPair.current,
        });
    };
  }, []);

  const staticRows = useRef<ColDef[]>([
    { field: 'id', hide: true },

    {
      headerName: Translate({ message: 'Sum' }),
      field: 'buy.sum',
      colId: 'buyPrice',
      width: vw(12.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
      cellRenderer: ({ data }) =>
        CellRenderer(
          <>
            <div
              className='overlay'
              style={{
                background: 'var(--greenText)',
                width: `${
                  data.buy.sum
                    ? (((+data.buy.sum / buyMax) * 100) / 2).toFixed(2)
                    : 0
                }%`,
              }}
            ></div>
            <div
              className='overlay'
              style={{
                background: 'var(--redText)',
                width: `${(((+data.sell.sum / sellMax) * 100) / 2).toFixed(
                  2,
                )}%`,
                right: 'unset',
                left: 'calc(50% + 0px)',
              }}
            ></div>
            <span className='small'>{data.buy.sum}</span>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Value' }),
      field: 'buy.value',
      colId: 'buyValue',
      width: vw(16.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
    },
    {
      headerName: Translate({ message: 'Amount' }),
      field: 'buy.amount',
      colId: 'buyAmount',
      width: vw(13.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
    },
    {
      headerName: Translate({ message: 'Bid' }),
      field: 'buy.price',
      colId: 'Bid',
      minWidth: 95,
      width: vw(7.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
      cellRenderer: ({ data }) =>
        CellRenderer(
          <>
            <span className='small' style={{ color: 'var(--greenText)' }}>
              {data.buy.price}
            </span>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Ask' }),
      field: 'sell.price',
      colId: 'Ask',
      minWidth: 95,
      //  maxWidth: 90,
      width: vw(7.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
      cellRenderer: ({ data }) =>
        CellRenderer(
          <>
            <span className='small' style={{ color: 'var(--redText)' }}>
              {data.sell.price}
            </span>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Amount' }),
      field: 'sell.amount',
      colId: 'sellAmount',
      width: vw(13.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
    },
    {
      headerName: Translate({ message: 'Value' }),
      field: 'sell.value',
      colId: 'sellValue',
      //  minWidth: 90,
      //  maxWidth: 90,
      width: vw(16.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
    },
    {
      headerName: Translate({ message: 'Sum' }),
      field: 'sell.sum',
      colId: 'sellSum',
      width: vw(12.5, viewportWidth),
      suppressMenu: true,
      sortable: false,
    },
  ]);
  const gridConfig = {
    columnDefs: [...staticRows.current],
    rowData: [],
  };
  const gridRendered = e => {
    gridApi.current = e.api;
    gridApi.current.sizeColumnsToFit();
    ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -11 });
  };

  const cellClicked = (params: CellClickedEvent) => {
    if (params.colDef.field) {
      MessageService.send({
        name: MessageNames.SELECT_ORDERBOOK_ROW,
        payload: {
          type: params.colDef.field.split('.')[0],
          data: params.data[params.colDef.field.split('.')[0]],
        },
      });
    }
  };
  const clearChart = () => {
    arr.current = [];
    drawDepth({ sell: [], buy: [] });
    gridApi.current.setRowData([]);
  };
  const onPrecisionChange = Pre => {
    clearChart();
    Precision = Pre;
  };
  return (
    <>
      <PrecisionSelect onPrecisionChange={onPrecisionChange} />
      <GridWrapper
        id={'ag-grid-wrapper-' + props.uniqueId}
        className='ag-theme-balham'
      >
        <AgGridReact
          onGridReady={gridRendered}
          animateRows={false}
          rowBuffer={1}
          viewportRowModelBufferSize={0}
          // debounceVerticalScrollbar
          immutableData
          getRowNodeId={data => {
            return data.id.toString();
          }}
          headerHeight={24}
          singleClickEdit
          rowHeight={24}
          onCellClicked={cellClicked}
          columnDefs={gridConfig.columnDefs}
          rowData={gridConfig.rowData}
          overlayNoRowsTemplate={ReactDOMServer.renderToString(
            <AnimatedNoRows
              texts={[
                <span className='black' style={{ fontSize: '9px' }}>
                  {Translate({ message: 'NoData' })}
                </span>,
              ]}
            />,
          )}
        ></AgGridReact>
        {/*<canvas id="greenCanvas"></canvas>
        <canvas id="redCanvas"></canvas>*/}
        <div id='orderBookChartWrapper' className='chartWrapper'>
          <div className='orderBookchart1'></div>
          <div className='orderBookchart2'></div>
        </div>
      </GridWrapper>
    </>
  );
};
export default injectIntl(OrderBookGrid);
const GridWrapper = styled.div`
  /*height: calc(15vh);*/
  .ag-header {
    opacity: 1;
    span {
      font-size: 10px;
      font-weight: 600;
    }
  }
  .ag-cell {
    line-height: 22px !important;
    font-size: 10px;
    color: var(--blackText);
    font-weight: 600;
    span {
      font-weight: 600;
      line-height: 22px !important;
      font-size: 10px;
    }
  }
  div[col-id='Bid'],
  div[col-id='sellSum'] {
    display: flex;
    place-content: flex-end;
    .ag-header-cell-label {
      display: flex;
      justify-content: flex-end;
    }
  }
  div[col-id='Ask'] {
    padding: 0;
  }
  div[col-id='buy'] {
  }

  .chartWrapper {
    pointer-events: none;
    position: absolute;
    top: 55px;
    left: 0;
    .orderBookchart1 {
      border-bottom-left-radius: 7px;
      overflow: hidden;
    }
    .orderBookchart2 {
      position: absolute;
      top: 0;
      border-bottom-right-radius: 7px;
      overflow: hidden;
    }
  }
  .overlay {
    height: 22px;
    position: fixed;
    right: calc(50% + 8px);
    opacity: 0.3;
  }
  .ag-row {
    cursor: pointer;
  }
`;
