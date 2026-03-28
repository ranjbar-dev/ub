import React, { useEffect, useRef, useCallback, useMemo } from 'react';
import { injectIntl } from 'react-intl';
import styled from 'styles/styled-components';
import { AgGridReact } from 'ag-grid-react';
import { ColDef, GridApi } from 'ag-grid-community';
import { GridHeaderNames, MqttTopicsPrefixes } from 'containers/App/constants';
import { Translator, vw, CurrencyFormater, zeroFixer } from 'utils/formatters';
import {
  SideSubscriber,
  MessageNames,
  Subscriber,
  MarketTradeSubscriber,
  dataInjectSubscriber,
  DataInjectMessageNames,
} from 'services/message_service';
import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';
import { ResizeGridHeigth } from '../utils/tradeUtilities';
import { MqttService } from 'services/MqttService2';
import { savedPairName } from 'utils/sharedData';
import { useDispatch } from 'react-redux';
import { getInitialMarketTradeDataAction } from '../actions';

const viewportWidth = 285;

const dataLength = 50;
let canUpdate = true;

const MarketTradeGrid = (props: {
  intl: any;
  subject: string;
  uniqueId: string;
  enabled: boolean;
}) => {
  const data: any = useRef([]);
  //  let gridApi: GridApi;
  const selectedPair = useRef(savedPairName());
  const gridApi: any = useRef();
  const mqtt2 = useRef(MqttService.getInstance());
  const initialDataLoaded = useRef(false);

  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(
      getInitialMarketTradeDataAction({ pairName: selectedPair.current }),
    );
  }, []);

  //////////////////////resize useEffect
  useEffect(() => {
    if (props.enabled) {
      mqtt2.current.ConnectToSubject({
        subject: `${MqttTopicsPrefixes.MarketTradeAddress}${selectedPair.current}`,
      });
    }
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESIZE) {
        setTimeout(() => {
          ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -13 });
          setTimeout(() => {
            gridApi.current.sizeColumnsToFit();
          }, 0);
        }, 185);
      }
      if (message.name === MessageNames.LAYOUT_RESIZE) {
        if (message.payload.i === 'MARKETTRADE') {
          requestAnimationFrame(() => {
            ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -13 });
            gridApi.current.sizeColumnsToFit();
          });
        }
      }
      if (message.name === MessageNames.LAYOUT_CHANGE) {
        setTimeout(() => {
          ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -13 });
          gridApi.current.sizeColumnsToFit();
        }, 220);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  ///////////////end of resize useEffect
  /////////trade mqtt

  useEffect(() => {
    let counter = 0;
    const SideSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        mqtt2.current.ConnectToNewSubject({
          newSubject: `${MqttTopicsPrefixes.MarketTradeAddress}${message.payload.name}`,
          oldsubject:
            MqttTopicsPrefixes.MarketTradeAddress + selectedPair.current,
        });
        selectedPair.current = message.payload.name;

        data.current = [];
        gridApi.current.setRowData([]);

      }
      initialDataLoaded.current = false;
      dispatch(
        getInitialMarketTradeDataAction({ pairName: selectedPair.current }),
      );
    });

    const dataInjectSubscription = dataInjectSubscriber.subscribe(
      (message: any) => {
        if (
          message.name === DataInjectMessageNames.MARKET_TRADES_INITIAL_DATA
        ) {
          if (gridApi.current && canUpdate === true) {
            counter = message.data.length;
            const gridData = [...message.data];
            gridData.forEach((item, index) => {
              item.id = index + item.price;
            });
            gridApi.current.setRowData(gridData);
          }
          initialDataLoaded.current = true;
        }
      },
    );

    const MarketTradeSubscription = MarketTradeSubscriber.subscribe(
      (message: any) => {
        if (
          gridApi.current &&
          canUpdate === true &&
          initialDataLoaded.current === true
        ) {
          const payload = message.payload;
          payload.id = counter + payload.price;
          counter++;
          data.current.unshift(payload);
          if (data.current.length < dataLength) {
            gridApi.current.applyTransactionAsync({
              add: [payload],
              addIndex: 0,
            });
          } else {
            const last = data.current.pop();
            gridApi.current.applyTransactionAsync({
              remove: [last],
              add: [payload],
              addIndex: 0,
            });
          }
        }
      },
    );
    return () => {
      SideSubscription.unsubscribe();

      MarketTradeSubscription.unsubscribe();

      dataInjectSubscription.unsubscribe();

      mqtt2.current.DisconnectFromSubject({
        subject: MqttTopicsPrefixes.MarketTradeAddress + selectedPair.current,
      });
    };
  }, []);

  const intl = props.intl;

  const Translate = useCallback((data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  }, []);

  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: Translate({ message: 'Price' }),
        field: 'price',
        colId: 'price',
        width: vw(35, viewportWidth),
        minWidth: 80,
        suppressMenu: true,
        cellStyle: ({ data }) => {
          return {
            color:
              data.isMaker === true ? 'var(--redText)' : 'var(--greenText)',
          };
        },
        valueFormatter: ({ data }) => {
          return zeroFixer(data.price + '');
        },
        sortable: false,
      },
      {
        headerName: Translate({ message: 'Amount' }),
        field: 'amount',
        colId: 'amount',
        width: vw(40, viewportWidth),
        suppressMenu: true,
        sortable: false,
      },
      {
        headerName: Translate({ message: 'Time' }),
        field: 'createdAt',
        width: vw(25, viewportWidth),
        minWidth: 70,
        suppressMenu: true,
        sortable: false,
        cellStyle: () => {
          return {
            textAlign: 'end',
          };
        },
        valueFormatter: ({ data }) => {
          return data.createdAt.split(' ')[1];
        },
      },
    ],
    [],
  );
  const gridConfig = {
    columnDefs: [...staticRows],
    rowData: [],
  };
  const gridRendered = (e) => {
    gridApi.current = e.api;
    gridApi.current.sizeColumnsToFit();
    ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: -13 });
  };
  const onScrollEvent = (e) => {
    if (e.top > 10) {
      canUpdate = false;
    } else {
      canUpdate = true;
    }
  };
  return (
    <GridWrapper
      id={'ag-grid-wrapper-' + props.uniqueId}
      className="ag-theme-balham"
    >
      <AgGridReact
        onGridReady={gridRendered}
        animateRows={true}
        headerHeight={24}
        //singleClickEdit={true}
        rowHeight={24}
        immutableData={true}
        getRowNodeId={(data) => {
          return data.id.toString();
        }}
        onBodyScroll={onScrollEvent}
        columnDefs={gridConfig.columnDefs}
        rowData={gridConfig.rowData}
        overlayNoRowsTemplate={ReactDOMServer.renderToString(
          <AnimatedNoRows
            texts={[
              <span className="black" style={{ fontSize: '9px' }}>
                {Translate({ message: 'NoData' })}
              </span>,
            ]}
          />,
        )}
      ></AgGridReact>
    </GridWrapper>
  );
};
export default injectIntl(MarketTradeGrid);
const GridWrapper = styled.div`
  height: calc(15vh);
  div[col-id='createdAt'] {
    .ag-header-cell-label {
      justify-content: flex-end;
    }
  }
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
  }
`;
