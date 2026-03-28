import React, { useEffect, useRef, useState } from 'react';
import ReactDOMServer from 'react-dom/server';

import { AgGridReact } from 'ag-grid-react';

import { WithdrawAddress } from '../types';
import { vw, Translator } from 'utils/formatters';
import styled from 'styles/styled-components';
import { Subscriber, MessageNames } from 'services/message_service';
import { CellRenderer } from 'components/renderer';
import deleteIcon from 'images/deleteIcon.svg';
import greyDeleteIcon from 'images/greyDeleteIcon.svg';
import { IconButton } from '@material-ui/core';
import { useDispatch } from 'react-redux';
import { deleteAddressAction, favoriteAddressAction } from '../actions';
import Anime from 'react-anime';
import { injectIntl } from 'react-intl';

import AnimatedNoRows from 'components/noRows/AnimatedNoRows';
import { GridFilterTypes, GridHeaderNames } from 'containers/App/constants';
import FilledStar from 'images/themedIcons/filledStar';
import EmptyAddressIcon from 'images/themedIcons/emptyAddressesIcon';
import EmptyStar from 'images/themedIcons/emptyStar';
import { ColDef, RowDataUpdatedEvent } from 'ag-grid-community';
import { headerHider } from 'utils/gridUtilities';
import { ConfirmPopup } from '../../../components/ConfirmPopup/confirmPopup';

let gridApi: any, gridColumnApi: any;
interface GridConfigTypes {
  columnDefs: any[];
  rowData: WithdrawAddress[];
}

const DataGrid = (props: { data: WithdrawAddress[]; intl: any }) => {
  const intl = props.intl;
  const Translate = (data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  };
  const rowToDeleteData = useRef<any>(null);
  const [isDeletePopupOpen, setIsDeletePopupOpen] = useState(false);
  const [IsEmpty, setIsEmpty] = useState(false);
  const gridData: WithdrawAddress[] = props.data;
  // if (gridData.length === 0 && IsEmpty === false && !gridApi) {
  //   setIsEmpty(true);
  // }
  const dispatch = useDispatch();
  const viewportWidth = 1152;
  ////hooks
  useEffect(() => {
    const subscription = Subscriber.subscribe((message: any) => {
      if (gridApi) {
        if (message.name === MessageNames.RESIZE) {
          gridApi.sizeColumnsToFit();
        }
        if (message.name == MessageNames.SET_GRID_FILTER) {
          //   gridApi.setQuickFilter(message.value);
          const filterInstance = gridApi.getFilterInstance(message.filterField);
          filterInstance.setModel({
            type: GridFilterTypes.Contains,
            filter: message.value,
          });
          // Get grid to run filter operation again
          gridApi.onFilterChanged();
        }

        if (message.name === MessageNames.ADD_DATA_ROW_TO_GRID) {
          gridApi.applyTransaction({
            add: [message.payload],
            addIndex: message.index,
          });
          setIsEmpty(false);
        }
        if (message.name === MessageNames.DELETE_GRID_ROW) {
          gridApi.forEachNode((node, index) => {
            if (node.data.id === message.payload.id) {
              gridApi.applyTransaction({ remove: [node.data] });
            }
          });
          const count = gridApi.getDisplayedRowCount();
          if (count == 0) {
            setTimeout(() => {
              setIsEmpty(true);
            }, 100);
          }
        }
        if (message.name === MessageNames.SET_FAVIORITE_ADDRESS) {
          gridApi.forEachNode((node, index) => {
            if (node.data.id === message.payload.id) {
              node.data.isFavorite = message.isFavorite;
              gridApi.applyTransaction({
                update: [node.data],
              });
              // gridApi.redrawRows();
            }

            // gridApi.onFilterChanged();
          });
        }
      }
    });
    return () => {
      subscription.unsubscribe();
    };
  }, []);

  ///////////

  //////buttons actions
  const deleteRow = (data: WithdrawAddress, rowIndex: number) => {
    rowToDeleteData.current = { data, rowIndex };
    setIsDeletePopupOpen(true);
  };

  const onConfirmDeleteClick = () => {
    dispatch(deleteAddressAction(rowToDeleteData.current));
    setIsDeletePopupOpen(false);
  };
  const onCancelDeleteClick = () => {
    setIsDeletePopupOpen(false);
  };

  const toggleFavorite = (addressData: WithdrawAddress, rowIndex: number) => {
    const data = {
      action: addressData.isFavorite === false ? 'add' : 'remove',
      id: addressData.id,
    };
    dispatch(favoriteAddressAction({ data, rowIndex }));
  };
  //////////////
  ////grid setup

  const staticRows: ColDef[] = [
    {
      headerName: '',
      field: 'isFavorite',
      width: vw(4, viewportWidth),
      suppressMenu: true,
      sortable: true,
      minWidth: 37,
      cellRenderer: ({ data, rowIndex }) =>
        CellRenderer(
          <div className='starWrapper'>
            <IconButton
              disableRipple={true}
              onClick={() => {
                toggleFavorite(data, rowIndex);
              }}
              className='headerButton star'
              size='small'
            >
              {data && data.isFavorite === true ? (
                <Anime duration={600} scale={[0.1, 1]}>
                  <FilledStar />
                </Anime>
              ) : (
                <Anime duration={600} scale={[0.1, 1]}>
                  <EmptyStar />
                </Anime>
              )}
            </IconButton>
          </div>,
        ),
    },
    {
      headerName: Translate({ message: 'Coin' }),
      field: 'code',
      width: vw(20, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: ({ data }) => {
        return `${data.code}${data.network ? ` (${data.network})` : ''}`;
      },
    },
    {
      headerName: Translate({ message: 'Label' }),
      field: 'label',
      width: vw(25, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: Translate({ message: 'Address' }),
      field: 'address',
      width: vw(45.8, viewportWidth),
      suppressMenu: true,
      sortable: true,
    },
    {
      headerName: '',
      field: 'delete',
      width: vw(5, viewportWidth),
      suppressMenu: true,

      cellRenderer: ({ data, rowIndex }) =>
        CellRenderer(
          <>
            <IconButton
              onClick={() => {
                deleteRow(data, rowIndex);
              }}
              className='headerButton delete'
              size='small'
            >
              <img
                src={rowIndex % 2 == 0 ? deleteIcon : greyDeleteIcon}
                alt=''
              />
            </IconButton>
          </>,
        ),
    },
  ];
  const gridConfig: GridConfigTypes = {
    columnDefs: [...staticRows],
    rowData: gridData,
  };
  /////////////////////////

  const gridRendered = e => {
    gridApi = e.api;
    gridColumnApi = e.columnApi;
    const width = window.innerWidth;
    if (width < viewportWidth) {
      gridApi.sizeColumnsToFit();
    }
    headerHider({ gridApi });
  };
  const rowDataUpdated = (e: RowDataUpdatedEvent) => {
    headerHider({ gridApi: e.api });
  };

  return (
    <>
      <ConfirmPopup
        title={<span>Delete this address ? </span>}
        cancelTitle='No'
        isOpen={isDeletePopupOpen}
        onCancelClick={onCancelDeleteClick}
        onSubmitClick={onConfirmDeleteClick}
        onClose={() => setIsDeletePopupOpen(false)}
        submitTitle='Yes'
      />
      <GridWrapper className='ag-theme-balham'>
        <AgGridReact
          onGridReady={gridRendered}
          // getRowStyle={rowStyles}
          // enableRtl={true}
          // pagination={true}
          // paginationPageSize={10}
          animateRows={true}
          onRowDataUpdated={rowDataUpdated}
          // localeText={translate}
          headerHeight={32}
          singleClickEdit={true}
          rowHeight={40}
          columnDefs={gridConfig.columnDefs}
          rowData={gridConfig.rowData}
          // onCellFocused={cellFocused}
          // onCellKeyPress={keyPressed}
          // onFirstDataRendered={gridRendered}
          overlayNoRowsTemplate={ReactDOMServer.renderToString(
            <AnimatedNoRows
              icon={<EmptyAddressIcon />}
              texts={[
                <span className='black'>
                  {Translate({ message: 'Youhavenowithdrawaddress' })}
                </span>,
                <span className='black'>
                  {Translate({
                    message: 'Pleasecreateaddressandwithdrawcoins',
                  })}
                </span>,
              ]}
            />,
          )}
        ></AgGridReact>
        <div className='whiteShaddow'></div>
        {/* {IsEmpty === true && (
        <EmptyWrapper>
          <Anime duration={1000} scale={[0.5, 1]}>
            <div className="itemsWrapper">
              <img src={emptyAddresses} alt="" />
              <FormattedMessage {...translate.Yourhavenowithdrawaddress} />
              <FormattedMessage
                {...translate.Pleasecreateaddressandwithdrawcoins}
              />
            </div>
          </Anime>
        </EmptyWrapper>
      )} */}
      </GridWrapper>
    </>
  );
};
export default injectIntl(DataGrid);
const GridWrapper = styled.div`
  height: calc(89vh - 265px);
  .whiteShaddow {
    position: absolute;
    width: 100%;
    height: 90px;

    left: 0px;
    bottom: 2px;
    pointer-events: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
  }
  .star {
    max-width: 30px;
    max-height: 30px;
    span {
      margin-top: -9px;
    }
  }
  .delete {
    margin-top: -3px;
  }
`;
