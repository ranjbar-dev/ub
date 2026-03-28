import React, { useEffect, useState, useMemo, useRef } from 'react';
import { Subscriber, MessageNames } from 'services/message_service';
import styled from 'styles/styled-components';
import { vw, Translator, CurrencyFormater } from 'utils/formatters';
import { CellRenderer } from 'components/renderer';

import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { AgGridReact } from 'ag-grid-react';

import Button from '@material-ui/core/Button';

import { DWData } from 'containers/FundsPage/types';
import greenTick from 'images/greenTickIcon.svg';
import { DWStatus } from 'containers/FundsPage/constants';
import { CircularProgress } from '@material-ui/core';
import withStyles from '@material-ui/styles/withStyles';
import {
  getOrderDetailAction,
  getInfiniteDWAction,
} from 'containers/FundsPage/actions';
import { useDispatch } from 'react-redux';
import DataRow from 'components/dataRow';
import TextLoader from 'components/textLoader';
import {
  DetailsStyle,
  ToggleDetail,
  RowWithShaddow,
  headerHider,
  whiteShaddowHider,
} from 'utils/gridUtilities';
import { injectIntl } from 'react-intl';
import { GridHeaderNames } from 'containers/App/constants';
import NoWDIcon from 'images/themedIcons/noWithdrawDeposit';
import ExpandMore from 'images/themedIcons/expandMore';
import { GridApi, ColDef } from 'ag-grid-community';
import { GridLoading } from 'components/grid_loading/gridLoading';
import { infinitePageSize } from 'containers/FundsPage/saga';
import TradeHistoryDetailWrapper from 'containers/FundsPage/pages/transactionsHistoryPage/components/detailWrapper';
import CrossIcon from 'images/themedIcons/crossIcon';
import { useNewPaymentEvent } from 'containers/FundsPage/hooks/useNewPaymentEvent';
interface Params {
  data: DWData;
  rowIndex: number;
  node: { id: string; rowIndex: number };
}
let gridApi: GridApi;
interface GridConfigTypes {
  columnDefs: ColDef[];
  rowData: DWData[];
  defaultColDef?: ColDef;
}
let PageNumber = 0;
let CanGetNewData = true;
let LoadedDataNumber = 0;

const ColorCircularProgress = withStyles({
  root: {
    color: '#707070',
  },
})(CircularProgress);

const DataGrid = (props: {
  sectionName: string;
  coinCode: string;
  intl: any;
}) => {
  //   let gridData: OpenOrder[] = props.data;
  const intl = props.intl;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };
  const [GridData, setGridData] = useState([]);
  const [IsLoadingInitialData, setIsLoadingInitialData] = useState(true);
  const [IsLoadingInfinite, setIsLoadingInfinite] = useState(false);

  const LastScroll = useRef(0);

  const dispatch = useDispatch();

  useNewPaymentEvent({
    dependencies: [props.sectionName],
    toRunAfterNewEvent: () => {
      PageNumber = 0;
      CanGetNewData = true;
      dispatch(
        getInfiniteDWAction({
          page: 0,
          // code: props.coinCode,
          page_size: infinitePageSize,
          type: props.sectionName,
          silent: true,
        }),
      );
    },
  });

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_INFINITE_SCROLL) {
        PageNumber = 0;
        CanGetNewData = true;
        LoadedDataNumber = 0;
      }
      if (message.name === MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA) {
        PageNumber = 0;
        LoadedDataNumber = message.payload.length;
        setGridData(message.payload);
      }
      if (
        message.name === MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING
      ) {
        setIsLoadingInitialData(message.payload);
      }
      if (gridApi) {
        if (message.name === MessageNames.RESIZE) {
          gridApi.sizeColumnsToFit();
        }
        if (message.name === MessageNames.SET_DATA_TO_INFINITE_BOTTOM) {
          setIsLoadingInfinite(false);
          if (
            message.payload.length === 0 ||
            message.payload.length < infinitePageSize
          ) {
            CanGetNewData = false;
          }
          gridApi.applyTransaction({
            add: [...message.payload],
            addIndex: LoadedDataNumber,
          });
          PageNumber++;
          LoadedDataNumber += message.payload.length;
        }
        if (message.name === MessageNames.ADD_DATA_ROW_TO_WITHDRAWS) {
          gridApi.applyTransaction({
            add: [message.payload],
            addIndex: 0,
          });
        }
        if (message.name === MessageNames.SET_ORDER_DETAIL) {
          gridApi.forEachNode((node, index) => {
            if (node.id === message.rowId) {
              node.data.details = message.payload;
              gridApi.applyTransaction({
                update: [node.data],
              });
              const dtButton = document.getElementById(
                'dtButton' + node.data.id,
              );
              if (dtButton) {
                dtButton.classList.add('grey');
              }
              return;
            }
          });
        }
      }
    });
    return () => {
      PageNumber = 0;
      CanGetNewData = true;
      LoadedDataNumber = 0;
      Subscription.unsubscribe();
    };
  }, []);

  const viewportWidth = 884;

  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: Translate({ message: 'Status' }),
        field: 'status',
        width: vw(21.5, viewportWidth),
        suppressMenu: true,
        minWidth: 40,
        sortable: true,
        cellRenderer: (params: { data: DWData }) =>
          CellRenderer(
            <>
              <div className='statusWrapper'>
                <div className={`${'icon' + params.data.status}`}>
                  {params.data.status === DWStatus.Complete ||
                  params.data.status === DWStatus.CONFIRMED ? (
                    <img src={greenTick} />
                  ) : params.data.status === DWStatus.IN_PROGRESS ||
                    params.data.status === DWStatus.PENDING ? (
                    <ColorCircularProgress disableShrink size={15} />
                  ) : params.data.status === DWStatus.Cancel ? (
                    <div className='dash'>-</div>
                  ) : (
                    <CrossIcon
                      style={{ marginRight: '11px' }}
                      color='var(--redText)'
                    />
                  )}
                </div>
                <div className={`upperFirst ${'text' + params.data.status}`}>
                  {params.data.status}
                </div>
              </div>
            </>,
          ),
      },
      {
        headerName: Translate({ message: 'Coin' }),
        field: 'code',
        width: vw(15, viewportWidth),
        suppressMenu: true,
        sortable: true,
      },

      {
        headerName: Translate({ message: 'Amount' }),
        field: 'amount',
        width: vw(30, viewportWidth),
        suppressMenu: true,
        sortable: true,
        valueFormatter: (params: any) => {
          return CurrencyFormater(params.data.amount);
        },
        comparator: (a, b) => {
          return a - b;
        },
      },
      {
        headerName: Translate({ message: 'Time' }),
        field: 'createdAt',
        minWidth: 145,
        width: vw(25, viewportWidth),
        suppressMenu: true,
        sortable: true,
      },

      {
        headerName: Translate({ message: 'Info' }),
        field: 'details',
        width: vw(15, viewportWidth),
        minWidth: 88,
        maxWidth: 88,
        suppressMenu: true,
        cellRenderer: (params: Params) =>
          CellRenderer(
            <>
              <div>
                <TradeHistoryDetailWrapper className='multiLine'>
                  <DataRow
                    title={Translate({ message: 'toAddress' })}
                    small
                    dense
                    titleWidth={'110px'}
                    clickAddress={
                      params.data.details &&
                      params.data.details.addressExplorerUrl
                    }
                    value={
                      params.data.details ? (
                        params.data.details.address
                      ) : (
                        <TextLoader />
                      )
                    }
                  />
                  {params.data.status !== DWStatus.Rejected ? (
                    <DataRow
                      title='Transaction ID: '
                      titleWidth={'120px'}
                      small
                      dense
                      clickAddress={
                        params.data.details &&
                        params.data.details.txIdExplorerUrl
                      }
                      value={
                        params.data.details ? (
                          params.data.details.txId
                        ) : (
                          <TextLoader />
                        )
                      }
                    />
                  ) : (
                    <DataRow
                      title='Reject Reason : '
                      titleWidth={'120px'}
                      small
                      dense
                      value={
                        params.data.details ? (
                          params.data.details.rejectionReason
                        ) : (
                          <TextLoader />
                        )
                      }
                    />
                  )}
                </TradeHistoryDetailWrapper>
              </div>
              {params.data.status === DWStatus.Complete ||
              params.data.status === DWStatus.CONFIRMED ||
              params.data.status === DWStatus.Rejected ||
              params.data.status === DWStatus.IN_PROGRESS ||
              params.data.status === DWStatus.PENDING ? (
                <Button
                  id={'dtButton' + params.data.id}
                  onClick={() => {
                    ToggleDetail({
                      gridApi,
                      params,
                      isOpen: params.data.isDetailsOpen,
                    });
                    if (!params.data.details) {
                      dispatch(
                        getOrderDetailAction({
                          id: params.data.id,
                          rowId: params.node.id,
                        }),
                      );
                    }
                  }}
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
                  // disableRipple
                  className='detailButton'
                  color='primary'
                >
                  {Translate({ message: 'Detail' })}
                </Button>
              ) : (
                <div style={{ paddingLeft: '10px' }}>{'-'}</div>
              )}
            </>,
          ),
      },
    ],
    [],
  );
  const gridConfig: GridConfigTypes = {
    columnDefs: [...staticRows],
    rowData: GridData,
  };
  const gridRendered = e => {
    gridApi = e.api;
    gridApi.sizeColumnsToFit();

    headerHider({
      gridApi,
      hideBottomBorderFromTitledComponent: true,
      wrapperId: 'depositAndWithdrawGridWrapper',
    });
  };

  const onBodyScroll = e => {
    const lastRowInView = gridApi.getLastDisplayedRow();
    if (
      lastRowInView == LoadedDataNumber - 1 &&
      LoadedDataNumber >= infinitePageSize &&
      GridData.length !== 0 &&
      CanGetNewData === true &&
      e.top > LastScroll.current &&
      e.top > 200
    ) {
      if (IsLoadingInfinite === false) {
        setIsLoadingInfinite(true);
        LastScroll.current = e.top;
        dispatch(
          getInfiniteDWAction({
            page: PageNumber + 1,
            // code: props.coinCode,
            page_size: infinitePageSize,
            type: props.sectionName,
          }),
        );
      }
    }
    whiteShaddowHider({ gridApi });
  };

  return useMemo(() => {
    if (IsLoadingInitialData === true) {
      return <GridLoading style={{ left: ' calc(50% - 30px)' }} />;
    }
    return (
      <GridWrapper className='ag-theme-balham withExpandableRows'>
        <AgGridReact
          onGridReady={gridRendered}
          animateRows={true}
          headerHeight={32}
          onBodyScroll={onBodyScroll}
          rowHeight={40}
          getRowStyle={RowWithShaddow}
          immutableData
          getRowNodeId={(data: any) => {
            return data.id.toString();
          }}
          // getRowHeight={RowWithDetailsHeight}
          columnDefs={gridConfig.columnDefs}
          rowData={gridConfig.rowData}
          overlayNoRowsTemplate={ReactDOMServer.renderToString(
            <AnimatedNoRows
              icon={<NoWDIcon />}
              texts={[
                <span className='black'>
                  {props.sectionName === 'withdraw'
                    ? intl.formatMessage({
                        id: 'containers.FundsPage.noWithdrawHistory',
                        defaultMessage: 'ET.noWithdrawHistory',
                      })
                    : intl.formatMessage({
                        id: 'containers.FundsPage.noDepositHistory',
                        defaultMessage: 'ET.noDepositHistory',
                      })}
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
        <div id='whiteShaddow' className='whiteShaddow'></div>
        {IsLoadingInfinite === true && (
          <div className='infiniteLoader'>
            <GridLoading
              style={{
                position: 'fixed',
                bottom: '3%',
                top: 'unset',
                left: 'calc(50% - 30px)',
                zIndex: '99',
              }}
            />
          </div>
        )}
      </GridWrapper>
    );
  }, [GridData, IsLoadingInitialData, IsLoadingInfinite]);
};
export default injectIntl(DataGrid);
const GridWrapper = styled.div`
  height: calc(91vh - 185px);
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

  .statusWrapper {
    display: flex;
    align-items: center;
    font-size: 13px;
    .dash {
      margin-left: 5px;
      margin-right: 18px;
    }

    .iconin,
    .iconpending {
      margin-right: 11px;
      margin-top: 3px;
    }
    .iconcompleted,
    .iconconfirmed {
      margin-left: -4px;
      margin-right: 7px;
      margin-top: -2px;
    }
    .iconcancel {
      padding: 0 5px;
    }
    .textin,
    .textpending {
      margin-top: -3px;
      color: var(--orange);
    }
    .textcancel {
      margin: 0 10px;
    }
    .textreject {
      color: var(--redText);
      margin: 0 10px;
    }
  }

  ${DetailsStyle}
  div[col-id='details'] {
    overflow: visible;
  }
  div[col-id='createdAt'],
  div[col-id='status'],
  div[col-id='code'],
  div[col-id='amount'] {
    pointer-events: none;
  }

  .infiniteLoader {
    position: fixed;
    bottom: 0;
    z-index: 1;
    width: 95%;
    text-align: center;
    height: 100%;
    background: transparent;
  }
  .upperFirst {
    text-transform: capitalize;
    font-size: 13px;
  }
`;
