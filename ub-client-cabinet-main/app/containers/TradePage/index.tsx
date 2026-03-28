/*
 *
 * TradePage
 *
 */

import React, { memo, useEffect } from 'react';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import reducer from 'containers/OrdersPage/reducer';
import saga from 'containers/OrdersPage/saga';
import StreamComponentsWrapper from './components/StreamComponentsWrapper';
import TradeHelmet from './helmet';
import { LocalStorageKeys } from 'services/constants';
import { useDispatch } from 'react-redux';
import { getCurrenciesAction } from './actions';
import styled from 'styles/styled-components';

interface Props { }

function TradePage(props: Props) {
  useInjectReducer({ key: 'ordersPage', reducer: reducer });
  useInjectSaga({ key: 'ordersPage', saga: saga });
  const dispatch = useDispatch();
  useEffect(() => {

    dispatch(getCurrenciesAction());

    return () => { };
  }, []);
  return (
    <Wrapper>
      <TradeHelmet />
      <StreamComponentsWrapper />
    </Wrapper>
  );
}
const Wrapper = styled.div`
  .ag-row {
    background-color: var(--white) !important;
    &:hover {
      background-color: var(--lightBlue) !important;
    }
  }
`;

export default memo(TradePage);
