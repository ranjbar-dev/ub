import React, { memo, useEffect } from 'react';
import { Order, FilterModel } from 'containers/OrdersPage/types';
import { useDispatch } from 'react-redux';
import { GridLoading } from 'components/grid_loading/gridLoading';
import { getTradeHistoryAction } from 'containers/OrdersPage/actions';
import AnimateChildren from 'components/AnimateChildren';
import DataGrid from './dataGrid';
import GridFilters from 'components/gridFilters';
import { LocalStorageKeys } from 'services/constants';
import { OrderPage } from 'containers/TradePage/constants';
import { MessageService, MessageNames } from 'services/message_service';
let searchCalled = false;
const TradeHistory = (props: {
  tradeHistory: Order[];
  isLoadingTradeHistory: boolean;
  isMini?: boolean;
}) => {
  const isLoadingData = props.isLoadingTradeHistory;
  const dispatch = useDispatch();
  useEffect(() => {
    localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] =
      OrderPage.TradeHistory;
    if (props.tradeHistory.length === 0) {
      dispatch(getTradeHistoryAction());
    } else {
      requestAnimationFrame(() => {
        //@ts-ignore
        dispatch(getTradeHistoryAction({ silent: true }));
      });
    }
    return () => {};
  }, []);
  const onSearchClick = (filters: FilterModel) => {
    searchCalled = true;
    MessageService.send({
      name: MessageNames.SET_PAGE_FILTERS_WITH_ID,
      id: 'tradeHistory',
      payload: filters,
    });
    dispatch(getTradeHistoryAction(filters));
  };
  const onCancelClick = () => {
    if (searchCalled === true) {
      MessageService.send({
        name: MessageNames.SET_PAGE_FILTERS_WITH_ID,
        id: 'tradeHistory',
        payload: {},
      });
      dispatch(getTradeHistoryAction());
      searchCalled = false;
    }
  };

  if (isLoadingData) {
    return <GridLoading className={`${props.isMini ? 'contained' : ''}`} />;
  }
  return (
    <AnimateChildren isLoading={false}>
      {!props.isMini && (
        <GridFilters
          onSearchClick={onSearchClick}
          onCancelClick={onCancelClick}
          hideTick={true}
        />
      )}
      <DataGrid isMini={props.isMini} data={props.tradeHistory} />
    </AnimateChildren>
  );
};
export default memo(TradeHistory);
