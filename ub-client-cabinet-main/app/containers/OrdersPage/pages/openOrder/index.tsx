import React, { memo, useEffect } from 'react';
import { Order } from 'containers/OrdersPage/types';
import { useDispatch } from 'react-redux';
import { getOpenOrdersAction } from 'containers/OrdersPage/actions';
import AnimateChildren from 'components/AnimateChildren';
import DataGrid from './dataGrid';
import { GridLoading } from 'components/grid_loading/gridLoading';
import { LocalStorageKeys } from 'services/constants';
import { OrderPage } from 'containers/TradePage/constants';
import useCookiesAuth from 'utils/hooks/useCookiesAuth';

const OpenOrder = (props: {
  openOrders: Order[];
  isMini?: boolean;
  isLoadingOpenOrders: boolean;
}) => {
  const isAuthed = useCookiesAuth();
  const isLoadingData = props.isLoadingOpenOrders;
  const dispatch = useDispatch();
  useEffect(() => {
    localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] = OrderPage.OpenOrders;
    let rafId: number;
    if (props.openOrders.length === 0 && isAuthed === true) {
      rafId = requestAnimationFrame(() => {
        dispatch(getOpenOrdersAction());
      });
    } else if (isAuthed === true) {
      rafId = requestAnimationFrame(() => {
        dispatch(getOpenOrdersAction({ silent: true }));
      });
    }
    return () => {
      if (rafId) cancelAnimationFrame(rafId);
    };
  }, [isAuthed]);
  if (isLoadingData) {
    return <GridLoading className={`${props.isMini ? 'contained' : ''}`} />;
  }
  return (
    <AnimateChildren isLoading={false}>
      <DataGrid isMini={props.isMini} data={props.openOrders} />
    </AnimateChildren>
  );
};
export default memo(OpenOrder);
