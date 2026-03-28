import React, { useState, useEffect, useRef } from 'react';
import styled from 'styles/styled-components';
import { createStructuredSelector } from 'reselect';
import makeSelectOrdersPage, {
  makeSelectLoggedIn,
} from 'containers/OrdersPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import OrderTabs from './tabs';
import OrderPageSelector from './pageSelector';
import { OrderPage } from 'containers/TradePage/constants';
import { LocalStorageKeys } from 'services/constants';
import {
  RegisteredUserSubscriber,
  MessageNames,
  SideSubscriber,
} from 'services/message_service';
import { getCurrencyPairInfoAction } from 'containers/OrdersPage/actions';
import { savedPairID } from 'utils/sharedData';

const stateSelector = createStructuredSelector({
  orders: makeSelectOrdersPage(),
  loggedIn: makeSelectLoggedIn(),
});

export default function Orders () {
  const { loggedIn, orders } = useSelector(stateSelector);

  const selectedPairId = useRef(savedPairID());
  useEffect(() => {
    const Subscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        selectedPairId.current = message.payload.id;
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const onTabChanged = e => {
    localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] = e;
    setSelectedPage(e);
  };

  useEffect(() => {
    localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] = OrderPage.OpenOrders;
    return () => {
      localStorage.removeItem(LocalStorageKeys.VISIBLE_ORDER_SECTION);
    };
  }, []);
  const [SelectedPage, setSelectedPage] = useState(OrderPage.OpenOrders);

  return (
    <Wrapper>
      <div className='tabs'>
        <OrderTabs onTabChange={onTabChanged} />
        <OrderPageSelector
          loggedIn={loggedIn}
          selectedStep={SelectedPage}
          ordersData={orders}
        />
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  background: var(--white);
  border-radius: var(--cardBorderRadius);
  .pageSelectorWrapper {
    padding: 0 5px;
  }
  .registration {
    width: 100%;
    height: -webkit-fill-available;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: -25px;
    min-height: 200px;
  }
`;
