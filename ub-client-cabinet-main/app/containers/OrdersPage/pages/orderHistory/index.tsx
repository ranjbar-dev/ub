import React, { memo, useEffect } from 'react';
import { Order } from 'containers/OrdersPage/types';
import { useDispatch } from 'react-redux';
import { GridLoading } from 'components/grid_loading/gridLoading';
import {
  getOrderHistoryAction,
  getFilteredOrderHistoryAction,
} from 'containers/OrdersPage/actions';
import AnimateChildren from 'components/AnimateChildren';
import DataGrid from './dataGrid';
import GridFilters from 'components/gridFilters';
import { LocalStorageKeys } from 'services/constants';
import { OrderPage } from 'containers/TradePage/constants';
import { MessageService, MessageNames } from 'services/message_service';
let searchCalled = false;
const OrderHistory = (props: {
  orderHistory: Order[];
  isLoadingOrderHistory: boolean;
  isMini?: boolean;
}) => {
  const isLoadingData = props.isLoadingOrderHistory;
  const dispatch = useDispatch();
  useEffect(() => {
    localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] =
      OrderPage.OrderHistory;
    if (props.orderHistory.length === 0) {
      dispatch(getOrderHistoryAction());
    } else {
      requestAnimationFrame(() => {
        dispatch(getOrderHistoryAction({ silent: true }));
      });
    }
    return () => {};
  }, []);
  const onSearchClick = (data: any) => {
    searchCalled = true;
    MessageService.send({
      name: MessageNames.SET_PAGE_FILTERS_WITH_ID,
      id: 'orderHistory',
      payload: data,
    });
    dispatch(getFilteredOrderHistoryAction(data));
  };
  const onCancelClick = () => {
    if (searchCalled === true) {
      dispatch(getOrderHistoryAction());
      searchCalled = false;
      MessageService.send({
        name: MessageNames.SET_PAGE_FILTERS_WITH_ID,
        id: 'orderHistory',
        payload: {},
      });
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
        />
      )}
      <DataGrid isMini={props.isMini} data={props.orderHistory} />
    </AnimateChildren>
  );
};
export default memo(OrderHistory);
