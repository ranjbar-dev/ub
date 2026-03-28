import React, { useEffect, useState, useRef } from 'react';
import { Subscriber, MessageNames } from 'services/message_service';
import styled from 'styles/styled-components';
import { vw, Translator, CurrencyFormater } from 'utils/formatters';
import { CellRenderer } from 'components/renderer';
import externalLink from 'images/themedIcons/externalLink.svg';
import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';

import { AgGridReact } from 'ag-grid-react';
import greenTick from 'images/greenTickIcon.svg';

import { Transaction } from 'containers/FundsPage/types';
import { GridHeaderNames } from 'containers/App/constants';
import withdrawIcon from 'images/withdrawIcon.svg';
import depositeIcon from 'images/depositeIcon.svg';
import { withStyles } from '@material-ui/core';
import { CircularProgress } from '@material-ui/core';
import { DWStatus } from 'containers/FundsPage/constants';
import { injectIntl } from 'react-intl';
import NoOpenOrdersIcon from 'images/themedIcons/noOpenOrdersIcon';
import { FilterApplier } from 'components/gridFilters/filterApply';
import {  ColDef } from 'ag-grid-community';
import { GridLoading } from 'components/grid_loading/gridLoading';
import {
  DetailsStyle,
  RowWithShaddow,
  headerHider,
  whiteShaddowHider,
  resizeToTarget,
} from 'utils/gridUtilities';
import CrossIcon from 'images/themedIcons/crossIcon';
//let gridApi: GridApi, gridColumnApi: any;
interface GridConfigTypes {
  columnDefs: any[];
  rowData: Transaction[];
}
interface Params {
  data: Transaction;
  rowIndex: number;
  node: { id: string };
}
const ColorCircularProgress = withStyles({
  root: {
    color: '#707070',
  },
})(CircularProgress);
const DataGrid = (props: { data: Transaction[]; intl: any }) => {
  //   let gridData: OpenOrder[] = props.data;
  const intl = props.intl;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };

  const gridApi = useRef<any>();

  const [GridData, setGridData] = useState(props.data);
  const [IsLoading, setIsLoading] = useState(false);

  const viewportWidth = 1552;

  const staticRows = useRef<ColDef[]>([
    {
      headerName: Translate({ message: 'Status' }),
      field: 'status',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      minWidth: 40,
      cellRenderer: (params: Params) =>
        CellRenderer(
          <>
            <div className='statusWrapper'>
              <div className={`${'icon' + params.data.status}`}>
                {params.data.status === DWStatus.Complete && (
                  <img
                    src={
                      params.data.type === 'deposit'
                        ? depositeIcon
                        : withdrawIcon
                    }
                  />
                )}
                {params.data.status === DWStatus.CONFIRMED ? (
                  <img src={greenTick} />
                ) : params.data.status === DWStatus.IN_PROGRESS ? (
                  <ColorCircularProgress disableShrink size={15} />
                ) : params.data.status === DWStatus.Cancel ? (
                  <div className='dash'>-</div>
                ) : params.data.status === DWStatus.Complete ? (
                  ''
                ) : (
                  <CrossIcon color='var(--redText)' />
                )}
              </div>
              <span className={`upperFirst ${'text' + params.data.status}`}>
                {params.data.status}
              </span>
            </div>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'Coin' }),
      field: 'code',
      width: vw(5, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Type' }),
      field: 'type',
      width: vw(7, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Amount' }),
      field: 'amount',
      width: vw(8, viewportWidth),
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
      headerName: Translate({ message: 'Time' }),
      field: 'createdAt',
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: true,
      minWidth: 155,
      filter: 'agNumberColumnFilter',
      filterValueGetter: params => {
        return params.data.createdAtToFilter;
      },
    },
    {
      headerName: Translate({ message: 'Info' }),
      field: 'address',
      width: vw(28, viewportWidth),
      suppressMenu: true,

      cellRenderer: (params: Params) =>
        CellRenderer(
          <>
            <div className='info'>
              <span className='address'>
                {Translate({ message: 'Address' })}:{' '}
              </span>
              <span
                onClick={() => {
                  params.data.txIdExplorerUrl &&
                    window.open(params.data.addressExplorerUrl);
                }}
                className={`${params.data.txIdExplorerUrl && 'clickable'}`}
              >
                {params.data.address}

                {params.data.txIdExplorerUrl &&
                  params.data.status !== DWStatus.Cancel && (
                    <img
                      src={externalLink}
                      style={{ marginTop: '-4px', marginLeft: '2px' }}
                    />
                  )}
              </span>
            </div>
          </>,
        ),
    },
    {
      headerName: Translate({ message: 'TransactionId' }),
      field: 'txId',
      width: vw(32, viewportWidth),
      suppressMenu: true,
      cellRenderer: (params: Params) =>
        CellRenderer(
          <span
            onClick={() => {
              params.data.txIdExplorerUrl &&
                window.open(params.data.txIdExplorerUrl);
            }}
            className={`txId ${params.data.txIdExplorerUrl && 'clickable'}`}
            id={'txId' + params.data.id}
          >
            {params.data.txId}
            {params.data.txIdExplorerUrl && <img src={externalLink} />}
          </span>,
        ),
    },

    // {
    //   headerName: 'Info',
    //   field: 'detail',
    //   width: vw(5, viewportWidth),
    //   suppressMenu: true,
    //   cellRenderer: (params: Params) =>
    //     CellRenderer(
    //       <>
    //         <div>
    //           <TradeHistoryDetailWrapper>
    //             <div className="dataWrapper">
    //               <div className="txid">
    //                 {Translate({ message: 'transactionId' })}
    //                 {' : '}
    //               </div>
    //               <div className="txidValue"> {params.data.txId}</div>
    //             </div>
    //           </TradeHistoryDetailWrapper>
    //         </div>

    //         {(params.data.status === DWStatus.Completed ||
    //           params.data.status === DWStatus.CONFIRMED) && (
    //           <IconButton
    //             onClick={() => {
    //               ToggleDetail({
    //                 gridApi,
    //                 params,
    //                 isOpen: params.data.isDetailsOpen,
    //                 height: 80,
    //               });
    //             }}
    //             size="small"
    //             color="primary"
    //           >
    //             <ExpandMore
    //               className="rotateAnimate200"
    //               id={'expandIcon' + params.data.id}
    //             />
    //           </IconButton>
    //         )}
    //       </>,
    //     ),
    // },
  ]);
  const gridConfig: GridConfigTypes = {
    columnDefs: [...staticRows.current],
    rowData: GridData,
  };
  const gridRendered = e => {
    gridApi.current = e.api;
    const width = window.innerWidth;
    if (width < viewportWidth) {
      //@ts-ignore
      gridApi.current.sizeColumnsToFit();
    }
    resizeToTarget({
      elementId: 'transactionHistoryGrid',
      resizeToElementWithId: 'transactionHistoryWrapper',
      additional: 70,
    });
    headerHider({ gridApi: gridApi.current });
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (gridApi.current) {
        if (message.name === MessageNames.RESIZE) {
          resizeToTarget({
            elementId: 'transactionHistoryGrid',
            resizeToElementWithId: 'transactionHistoryWrapper',
            additional: 70,
          });
          //@ts-ignore
          gridApi.current.sizeColumnsToFit();
        }
        if (message.name === MessageNames.SET_GRID_DATA) {
          //  gridApi.current.setRowData(message.payload);
          setGridData(message.payload);
          //   gridData = message.payload;
        }
        if (message.name === MessageNames.SET_GRID_FILTER) {
          FilterApplier({ gridApi: gridApi.current, message: message });
        }

        if (message.name === MessageNames.SETLOADING) {
          setIsLoading(message.payload);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  if (IsLoading === true) {
    return <GridLoading style={{ left: 'calc(40vw - 30px)' }} />;
  }
  const onScroll = e => {
    whiteShaddowHider({ gridApi: gridApi.current });
  };
  return (
    <GridWrapper
      id='transactionHistoryGrid'
      className='ag-theme-balham withExpandableRows'
    >
      <AgGridReact
        onGridReady={gridRendered}
        animateRows={true}
        headerHeight={32}
        onBodyScroll={onScroll}
        singleClickEdit={true}
        rowHeight={40}
        immutableData
        getRowNodeId={data => {
          return data.id.toString();
        }}
        getRowStyle={RowWithShaddow}
        columnDefs={gridConfig.columnDefs}
        rowData={gridConfig.rowData}
        overlayNoRowsTemplate={ReactDOMServer.renderToString(
          <AnimatedNoRows
            icon={<NoOpenOrdersIcon />}
            texts={[
              <span className='black'>
                {props.intl.formatMessage({
                  id: 'containers.FundsPage.noTransactionHistory',
                  defaultMessage: '',
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
    </GridWrapper>
  );
};
export default injectIntl(DataGrid);
const GridWrapper = styled.div`
  ${DetailsStyle}
  .whiteShaddow {
    position: absolute;
    width: 100%;
    height: 58px;
    left: 0px;
    pointer-events: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
    bottom: 0px;
    z-index: 1;
  }
  .clickable {
    cursor: pointer;
  }
  .actionButton {
    top: 0px;
    max-height: 35px;
    padding: 0 12px !important;
    min-height: 32px;
    margin: 0 5px;
    left: -15px;
    font-weight: 600;
  }

  .coinWrapper {
    img {
      height: 25px;
      width: 25px;
      /* border: 1px solid #c1c1c1;
      border-radius: 50px;
      padding: 1px; */
      min-width: 25px;
      min-height: 25px;
    }
    span {
      margin: 0 0px;
      font-weight: 600;
      font-size: 13px;
    }
    .progress {
      margin-left: 5px;
      margin-right: 5px;
      margin-top: 5px;
      margin-bottom: -3px;
    }
  }
  .info {
    span {
      font-size: 12px !important;
    }
    .address {
      color: var(--textGrey);
    }
  }
  .textin {
    color: var(--orange);
  }
  .statusWrapper {
    display: flex;
    align-items: center;
    font-size: 13px;
    .dash {
      margin-left: 5px;
      margin-right: 5px;
    }
    .iconin {
      margin-right: 12px;
      margin-top: 3px;
    }
    .iconcompleted,
    .iconconfirmed {
      margin-left: -4px;
      margin-right: 5px;
      margin-top: -2px;
    }
    .iconcancel {
      padding: 0 5px;
    }
    .textin {
      margin-top: -3px;
      color: var(--orange);
    }
    .textcompleted {
      color: var(--greenText);
    }

    .textrejected,
    .textfailed,
    .textcanceled {
      margin: 0 10px;
    }
    .textrejected,
    .textfailed {
      color: var(--redText);
    }
    .upperFirst {
      text-transform: capitalize !important;
      font-size: 13px;
    }
  }
  .txId {
    color: var(--blackText);
    font-size: 12px;
    transition: opacity 0.2s;
    img {
      margin-top: -4px;
      margin-left: 2px;
    }
  }
`;
