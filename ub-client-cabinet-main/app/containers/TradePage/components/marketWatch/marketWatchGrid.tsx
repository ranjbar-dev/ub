import React, {
  useState,
  useEffect,
  useCallback,
  useMemo,
  useRef,
} from 'react';
import { injectIntl } from 'react-intl';
import styled from 'styles/styled-components';
import { AgGridReact } from 'ag-grid-react';
import { ColDef, CellClickedEvent } from 'ag-grid-community';
import {
  GridHeaderNames,
  GridFilterTypes,
  CentrifugoChannels,
} from 'containers/App/constants';
import { Translator, vw, PairFormat, zeroFixer } from 'utils/formatters';
import {
  MessageNames,
  SideMessageService,
  Subscriber,
  MarketWatchSubscriber,
  EventSubscriber,
  EventMessageNames,
} from 'services/message_service';
import ReactDOMServer from 'react-dom/server';
import AnimatedNoRows from 'components/noRows/AnimatedNoRows';
import { ResizeGridHeigth } from 'containers/TradePage/utils/tradeUtilities';
import SearchIcon from 'images/themedIcons/searchIcon';
import CoinTabs from './coinTabs';
import { CellRenderer } from 'components/renderer';
import FilledStar from 'images/themedIcons/filledStar';
import EmptyStar from 'images/themedIcons/emptyStar';
import { LocalStorageKeys } from 'services/constants';
import Anime from 'react-anime';
import { CentrifugoPublicService } from 'services/CentrifugoPublicService';
import { useInjectReducer } from 'utils/injectReducer';
import { useInjectSaga } from 'utils/injectSaga';
import reducer from 'containers/TradePage/reducer';
import saga from 'containers/TradePage/saga';
import { useDispatch } from 'react-redux';
import { AddRemoveFavoritePair } from 'containers/TradePage/actions';
import { CookieKeys, cookies } from 'services/cookie';
import { storage } from 'utils/storage';

const viewportWidth = 285;
let favoritePairs = localStorage[LocalStorageKeys.FAV_PAIRS]
  ? JSON.parse(localStorage[LocalStorageKeys.FAV_PAIRS])
  : [];

const MarketWatchGrid = (props: {
  intl: any;
  subject: string;
  uniqueId: string;
  enabled: boolean;
}) => {
  const { intl } = props;
  useInjectReducer({ key: 'tradePage', reducer: reducer });
  useInjectSaga({ key: 'tradePage', saga: saga });
  const dispatch = useDispatch();

  const mqtt2 = useRef(CentrifugoPublicService.getInstance());

  const gridApi: any = useRef();
  const columnApi: any = useRef();
  //////////////////////resize useEffect
  useEffect(() => {
    if (props.enabled === true) {
      mqtt2.current.ConnectToSubject({
        subject: CentrifugoChannels.TickerChannel,
      });
    }
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESIZE) {
        setTimeout(() => {
          ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: 15 });
          setTimeout(() => {
            gridApi.current.sizeColumnsToFit();
          }, 0);
        }, 190);
      }
      if (message.name === MessageNames.LAYOUT_RESIZE) {
        if (message.payload.i === 'MARKETWATCH') {
          requestAnimationFrame(() => {
            ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: 15 });
            gridApi.current.sizeColumnsToFit();
          });
        }
      }
      //clear the favorite pairs when logged out
      if (message.name === MessageNames.LOGGED_OUT) {
        gridApi.current.forEachNode((node, index) => {
          node.data.isFavorite = false;
          gridApi.current.applyTransactionAsync({ update: [node.data] });
        });
        localStorage.removeItem(LocalStorageKeys.FAV_PAIRS);
      }
    });
    return () => {
      mqtt2.current.DisconnectFromSubject({
        subject: CentrifugoChannels.TickerChannel,
      });
      Subscription.unsubscribe();
    };
  }, []);

  ////price mqtt
  useEffect(() => {
    const EventSubscription = EventSubscriber.subscribe((message: any) => {
      if (message.name === EventMessageNames.GOT_FAV_PAIRS) {
        favoritePairs = message.payload;
        if (gridApi.current) {
          gridApi.current.forEachNode((node, index) => {
            if (favoritePairs.indexOf(node.data.name) !== -1) {
              node.data.isFavorite = true;
              gridApi.current.applyTransactionAsync({ update: [node.data] });
            } else {
              node.data.isFavorite = false;
              gridApi.current.applyTransactionAsync({ update: [node.data] });
            }
          });
        }
      }
    });
    return () => {
      EventSubscription.unsubscribe();
    };
  }, []);
  ///////////////end of resize useEffect
  const [GridData, setGridData] = useState([]);

  const Translate = useCallback((data: { message: string }) => {
    return Translator({
      containerPrefix: GridHeaderNames,
      intl,
      message: data.message,
    });
  }, []);

  useEffect(() => {
    const allPairs = {};
    const MarketWatchSubscription = MarketWatchSubscriber.subscribe(
      (message: any) => {
        if (!document.hidden && gridApi.current) {
          if (!allPairs[message.payload.id]) {
            allPairs[message.payload.id] = message.payload;
            gridApi.current.applyTransactionAsync({
              add: [
                {
                  ...message.payload,
                  isFavorite: favoritePairs.includes(message.payload.name),
                },
              ],
              addIndex: 0,
            });
            return;
          }

          const rowNode = gridApi.current.getRowNode(
            message.payload.id.toString(),
          );
          if (rowNode) {
            message.payload.isFavorite = rowNode.data.isFavorite
              ? rowNode.data.isFavorite
              : false;
            rowNode.data = { ...message.payload };
            gridApi.current.applyTransactionAsync({ update: [rowNode.data] });
          }
        }
      },
    );
    return () => {
      MarketWatchSubscription.unsubscribe();
    };
  }, []);

  const staticRows = useRef<ColDef[]>([
    { headerName: 'id', field: 'id', hide: true },
    {
      headerName: '',
      field: 'isFavorite',
      colId: 'fav',
      minWidth: 30,
      maxWidth: 30,
      width: vw(10, viewportWidth),
      suppressMenu: true,
      sortable: false,
      cellRenderer: (params) =>
        CellRenderer(
          <>
            <Anime
              duration={200}
              easing="easeOutCirc"
              scale={[0.5, 1]}
              opacity={[0, 1]}
            >
              <div style={{ marginTop: '-2px' }}>
                {params.data.isFavorite === true ? (
                  <FilledStar size="18" />
                ) : (
                  <EmptyStar size="18" />
                )}
              </div>
            </Anime>
          </>,
        ),
      getQuickFilterText: function ({ data }) {
        return data.isFavorite === true ? 'showFavs' : 'showNotCancelled';
      },
    },
    {
      headerName: Translate({ message: 'Coin' }),
      field: 'name',
      width: vw(35, viewportWidth),
      suppressMenu: true,
      sortable: true,
      // maxWidth: 90,
      getQuickFilterText: ({ data }) => {
        return PairFormat(data.name);
      },
      valueFormatter: ({ data }) => PairFormat(data.name),
    },
    {
      headerName: Translate({ message: 'LastPrice' }),
      field: 'price',
      width: vw(30, viewportWidth),
      suppressMenu: true,
      sortable: true,
      valueFormatter: (params: any) => {
        return zeroFixer(params.data.price);
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
    {
      headerName: Translate({ message: 'Change' }),
      field: 'percentage',
      width: vw(25, viewportWidth),
      suppressMenu: true,
      sortable: true,
      cellStyle: function ({ data }) {
        return {
          color: +data.percentage > 0 ? 'var(--greenText)' : 'var(--redText)',
          textAlign: 'end',
        };
      },
      valueFormatter: ({ data }) => {
        return Number(data.percentage).toFixed(2) + '%';
      },
      comparator: function (a, b) {
        return a - b;
      },
    },
  ]);
  const gridConfig = useMemo(
    () => ({
      columnDefs: [...staticRows.current],
      rowData: GridData,
    }),
    [],
  );
  const gridRendered = useCallback((e) => {
    gridApi.current = e.api;
    gridApi.current.sizeColumnsToFit();
    columnApi.current = e.columnApi;

    ResizeGridHeigth({ uniqueId: props.uniqueId, additinal: 15 });

    const initialSort = [{ colId: 'price', sort: 'desc' }];
    columnApi.current.applyColumnState({ state: initialSort });

    if (localStorage[LocalStorageKeys.FAV_COIN] == 'Favs') {
      gridApi.current.setQuickFilter('showFavs');
    } else if (localStorage[LocalStorageKeys.FAV_COIN] == 'All') {
      gridApi.current.setQuickFilter('');
      gridApi.current.destroyFilter('name');
    } else if (localStorage[LocalStorageKeys.FAV_COIN]) {
      const filterInstance = gridApi.current.getFilterInstance('name');
      filterInstance.setModel({
        type: GridFilterTypes.Ends_With,
        filter: localStorage[LocalStorageKeys.FAV_COIN],
      });
      gridApi.current.onFilterChanged();
    }
  }, []);
  const onCellClicked = useCallback((e: CellClickedEvent) => {
    const field = e.colDef.field;
    if (field === 'isFavorite') {
      let isfav = false;
      gridApi.current.forEachNode((node, index) => {
        if (node.data.id === e.data.id) {
          e.data.isFavorite = isfav = e.data.isFavorite
            ? !e.data.isFavorite
            : true;
          node.data = { ...e.data };
          gridApi.current.applyTransactionAsync({ update: [node.data] });
          if (isfav === true) {
            favoritePairs.push(node.data.name);
            localStorage[LocalStorageKeys.FAV_PAIRS] = JSON.stringify(
              favoritePairs,
            );
            if (cookies.get(CookieKeys.Token)) {
              dispatch(
                AddRemoveFavoritePair({
                  pair_currency_id: node.data.id,
                  action: 'add',
                }),
              );
            }
          } else {
            favoritePairs = favoritePairs.filter((e) => e !== node.data.name);
            localStorage[LocalStorageKeys.FAV_PAIRS] = JSON.stringify(
              favoritePairs,
            );
            if (cookies.get(CookieKeys.Token)) {
              dispatch(
                AddRemoveFavoritePair({
                  pair_currency_id: node.data.id,
                  action: 'remove',
                }),
              );
            }
          }
          return;
        }
      });
      return;
    }
    storage.write(LocalStorageKeys.SAVED_TRADE_PAIR, {
      name: e.data.name,
      id: e.data.id,
    });
    SideMessageService.send({
      name: MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR,
      payload: e.data,
    });
  }, []);
  const onSearch = (e) => {
    gridApi.current.setQuickFilter(e.target.value);
  };
  const handleTabChange = useCallback((code: string) => {
    const filterInstance = gridApi.current.getFilterInstance('name');
    if (code === 'All') {
      gridApi.current.setQuickFilter('');
      gridApi.current.destroyFilter('name');
    } else if (code === 'Favs') {
      gridApi.current.setQuickFilter('showFavs');
      gridApi.current.destroyFilter('name');
      return;
    } else {
      filterInstance.setModel({
        type: GridFilterTypes.Ends_With,
        filter: code,
      });
      gridApi.current.setQuickFilter('');
    }
    gridApi.current.onFilterChanged();
  }, []);

  return (
    <GridWrapper
      id={'ag-grid-wrapper-' + props.uniqueId}
      className="ag-theme-balham"
    >
      <div className="search">
        <input
          placeholder={Translator({
            containerPrefix: 'containers.FundsPage',
            intl,
            message: 'SearchCoin',
          })}
          type="text"
          onChange={onSearch}
          className="searchInput"
        />
        <SearchIcon />
      </div>
      <CoinTabs
        subject={props.subject}
        //onFavChange={handleFavChange}
        onTabChange={handleTabChange}
      />
      <AgGridReact
        onGridReady={gridRendered}
        animateRows={true}
        headerHeight={24}
        singleClickEdit={true}
        rowHeight={24}
        immutableData={true}
        getRowNodeId={(data) => {
          return data.id.toString();
        }}
        //enableCellChangeFlash={true}
        // onRowClicked={onRowClick}
        onCellClicked={onCellClicked}
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
export default injectIntl(MarketWatchGrid);
const GridWrapper = styled.div`
  height: calc(15vh);
  .ag-header {
    opacity: 1;
    span {
      font-size: 10px;
      font-weight: 600;
    }
    div[col-id='percentage'] {
      .ag-header-cell-label {
        justify-content: flex-end;
      }
    }
  }
  .ag-cell {
    line-height: 22px !important;
    font-size: 10px;
    color: var(--blackText);
    font-weight: 600;
  }
  .ag-row {
    cursor: pointer;
  }
  .search {
    position: absolute;
    z-index: 1;
    top: 1px;
    right: 10px;
    .searchInput {
      border: none;
      box-shadow: none;
      outline: none !important;
      border-radius: 5px;
      margin: 0 -27px;
      height: 16px;
      font-size: 10px !important;
      padding: 0 10px;
      width: 100px;
      font-weight: 500;
      background: var(--white);
      color: var(--blackText);
      &::placeholder {
        font-weight: 500;
        font-size: 10px !important;
        font-style: italic;
        color: var(--placeHolderColor);
      }
    }
    svg {
      pointer-events: none;
      transform: scale(0.8);
      margin-top: 0px;
      margin-left: 3px;
    }
  }
`;
