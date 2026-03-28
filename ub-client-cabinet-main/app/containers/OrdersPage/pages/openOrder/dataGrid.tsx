import React, { useEffect, useState, useRef } from 'react';
import { Order } from 'containers/OrdersPage/types';
import {
  Subscriber,
  MessageNames,
  EventSubscriber,
  EventMessageNames,
} from 'services/message_service';
import styled from 'styles/styled-components';
import { vw, Translator, CurrencyFormater, PairFormat } from 'utils/formatters';
import { CellRenderer } from 'components/renderer';

import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { Button, makeStyles, CircularProgress } from '@material-ui/core';
import { Buttons, GridHeaderNames, AppPages } from 'containers/App/constants';
import { injectIntl } from 'react-intl';
import NoOpenOrdersIcon from 'images/themedIcons/noOpenOrdersIcon';

import { AgGridReact } from 'ag-grid-react';
import {
  headerHider,
  resizeToTarget,
  whiteShaddowHider,
} from 'utils/gridUtilities';
import { ColDef, GridReadyEvent } from 'ag-grid-community';
import { useDispatch } from 'react-redux';
import {
  cancelOrderAction,
  getOpenOrdersAction,
} from 'containers/OrdersPage/actions';
import { push } from 'redux-first-history';
import PopupModal from 'components/materialModal/modal';
import { OrderPages } from 'containers/OrdersPage/constants';
import { miniGridStyles } from 'containers/OrdersPage/styles';

interface GridConfigTypes {
  columnDefs: ColDef[];
  rowData: Order[];
}
const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'var(--textBlue)',
  },
});
const DataGrid = (props: { data: Order[]; intl: any; isMini?: boolean }) => {
  const gridApi: any = useRef();
  const classes = materialClasses();
  const [IsSubmitCancelOpen, setIsSubmitCancelOpen] = useState(false);
  //   let gridData: OpenOrder[] = props.data;

  const cancelData = useRef<any>();

  const intl = props.intl;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };
  const [GridData, setGridData] = useState(props.data);
  const dispatch = useDispatch();
  const viewportWidth = 1552;
  const handleCancelOrder = (data: Order) => {
    cancelData.current = data;
    setIsSubmitCancelOpen(true);
  };

  const rowsConfig = useRef<ColDef[]>([
    {
      headerName: Translate({ message: 'Date' }),
      field: 'createdAt',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      minWidth: 145,
      maxWidth: 155,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Pair' }),
      field: 'pair',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => PairFormat(data.pair),
    },
    {
      headerName: Translate({ message: 'Type' }),
      field: 'type',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Side' }),
      field: 'side',
      width: vw(10, viewportWidth),
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
        return CurrencyFormater(params.data.price);
      },
      suppressMenu: true,
      comparator: function (a, b) {
        return a - b;
      },
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Amount' }),
      field: 'amount',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.amount);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'FilledPercent' }),
      field: 'executed',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: ({ data }) => {
        return data.executed && data.executed.length > 0
          ? { paddingLeft: '20px' }
          : { paddingLeft: '30px' };
      },
      valueFormatter: ({ data }) => {
        return data.executed && data.executed.length > 0 ? data.executed : '-';
      },
    },
    {
      headerName: Translate({ message: 'Total' }),
      field: 'total',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return CurrencyFormater(params.data.total);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'TriggerConditions' }),
      field: 'triggerCondition',
      width: vw(15, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: ({ data }) => {
        return data.triggerCondition && data.triggerCondition.length > 0
          ? { paddingLeft: '0px' }
          : { paddingLeft: '50px' };
      },
      valueFormatter: ({ data }) => {
        return data.triggerCondition && data.triggerCondition.length > 0
          ? data.triggerCondition
          : '-';
      },
    },
    {
      headerName: Translate({ message: 'Action' }),
      field: 'delete',
      width: vw(5, viewportWidth),
      minWidth: 80,
      maxWidth: 80,
      suppressMenu: true,

      cellRenderer: ({ data, rowIndex }) =>
        CellRenderer(
          <>
            <Button
              onClick={() => {
                !data.isCanceling ? handleCancelOrder(data) : () => {};
              }}
              className={`cancelButtons detailButton ${
                props.isMini ? 'miniCancel' : ''
              }`}
              color='primary'
            >
              {!data.isCanceling ? (
                Translate({ message: 'Cancel' })
              ) : (
                <CircularProgress
                  size={14}
                  className={classes.loadingIndicator}
                />
              )}
            </Button>
          </>,
        ),
    },
  ]);

  const gridConfig: GridConfigTypes = {
    columnDefs: [...rowsConfig.current],
    rowData: GridData,
  };
  const gridRendered = (e: GridReadyEvent) => {
    gridApi.current = e.api;

    const width = window.innerWidth;
    if (width < viewportWidth || props.isMini) {
      gridApi.current.sizeColumnsToFit();
      if (props.isMini) {
        resizeToTarget({
          elementId: 'OpenOrdersGridWrapper',
          resizeToElementWithId: 'ORDERSContainer',
        });
      }
    }
    headerHider({
      gridApi: gridApi.current,
      gridId: 'OpenOrdersGridWrapper',
      hideBottomBorderFromTitledComponent: true,
    });
  };

  useEffect(() => {
    const EventSubscription = EventSubscriber.subscribe((message: any) => {
      if (
        message.name === EventMessageNames.REFRESH_ORDER_GRID ||
        message.name === MessageNames.RECONNECT_EVENT
      ) {
        dispatch(getOpenOrdersAction({ silent: true }));
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
              elementId: 'OpenOrdersGridWrapper',
              resizeToElementWithId: 'ORDERSContainer',
            });
            setTimeout(() => {
              gridApi.current.sizeColumnsToFit();
            }, 0);
          }, 200);
        }
        if (message.name === MessageNames.LAYOUT_RESIZE) {
          if (message.payload && message.payload.i === 'ORDERS') {
            //setTimeout(() => {

            //}, 0);
            requestAnimationFrame(() => {
              resizeToTarget({
                elementId: 'OpenOrdersGridWrapper',
                resizeToElementWithId: 'ORDERSContainer',
              });
              gridApi.current.sizeColumnsToFit();
            });
          }
        }
        if (message.name === MessageNames.SET_OPEN_ORDERS_DATA) {
          gridApi.current.setRowData(message.payload);
          headerHider({
            gridApi: gridApi.current,
            gridId: 'OpenOrdersGridWrapper',
            hideBottomBorderFromTitledComponent: true,
          });
          // setGridData(message.payload);
          //   gridData = message.payload;
        }
        if (message.name === MessageNames.IS_CANCELING_ORDER) {
          gridApi.current.forEachNode((node, index) => {
            if (node.data.id === message.payload.id) {
              node.data.isCanceling = message.payload.state;
              gridApi.current.applyTransaction({
                update: [node.data],
              });
              gridApi.current.redrawRows({ rowNodes: [node] });
              return;
            }
          });
          // setGridData(message.payload);
          //   gridData = message.payload;
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const onScroll = () => {
    whiteShaddowHider({ gridApi: gridApi.current, isMini: props.isMini });
  };

  const handleSubmitCancel = () => {
    dispatch(cancelOrderAction(cancelData.current));
    setIsSubmitCancelOpen(false);
    cancelData.current = null;
  };

  return (
    <>
      <PopupModal
        isOpen={IsSubmitCancelOpen}
        onClose={() => {
          cancelData.current = null;
          setIsSubmitCancelOpen(false);
        }}
      >
        <div className='alertWrapper alertConfirmWrapper'>
          <span>
            {intl.formatMessage({
              id: 'containers.GoogleAuthenticationPage.Areyousure',
              defaultMessage: 'ET.Label',
            })}
            <span className='red'>
              {' '}
              {intl.formatMessage({
                id: 'containers.OrdersPage.youWantToCancelTheOrder',
                defaultMessage: 'ET.Label',
              })}
            </span>
            {intl.formatMessage({
              id: 'containers.GoogleAuthenticationPage.question',
              defaultMessage: 'ET.Label',
            })}
          </span>
        </div>
        <div className='alertButtonsWrapper'>
          <Button
            onClick={() => {
              cancelData.current = null;
              setIsSubmitCancelOpen(false);
            }}
          >
            {intl.formatMessage({
              id: 'containers.globalTitles.no',
              defaultMessage: 'ET.no',
            })}
          </Button>
          <div className='separator'></div>
          <Button onClick={handleSubmitCancel}>
            {intl.formatMessage({
              id: 'containers.globalTitles.yes',
              defaultMessage: 'ET.yes',
            })}
          </Button>
        </div>
      </PopupModal>

      <GridWrapper
        id='OpenOrdersGridWrapper'
        className={`ag-theme-balham ${props.isMini ? 'miniGrid' : ''}`}
      >
        <AgGridReact
          onGridReady={gridRendered}
          animateRows={true}
          onBodyScroll={onScroll}
          headerHeight={props.isMini ? 25 : 32}
          rowHeight={props.isMini ? 25 : 40}
          columnDefs={gridConfig.columnDefs}
          rowData={gridConfig.rowData}
          immutableData={true}
          getRowNodeId={(data: Order) => {
            return data.id.toString();
          }}
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
        {GridData && GridData.length === 0 && props.isMini === undefined && (
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
  height: calc(90vh - 130px);
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

  .cancelButtons {
    position: absolute;
    top: 2px;
    max-height: 35px;
    padding: 0 !important;
    min-height: 32px;
    left: 0px;
  }
`;
