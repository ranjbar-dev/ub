import React, { memo, useEffect } from 'react';
import { createStructuredSelector } from 'reselect';
import {
  makeSelectTransactionHistoryPageData,
  makeSelectIsLoadingTransactionHistoryPageData,
} from 'containers/FundsPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import { Transaction } from 'containers/FundsPage/types';
import AnimatedChildren from 'components/AnimateChildren';
import DataGrid from './components/dataGrid';
import TitledComponent from 'components/titled';
import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';
import { getTransactionHistoryPageDataAction } from 'containers/FundsPage/actions';
import GridFilters from 'components/gridFilters';
import { MessageNames, MessageService } from 'services/message_service';
import { useNewPaymentEvent } from 'containers/FundsPage/hooks/useNewPaymentEvent';

const stateSelector = createStructuredSelector({
  transactionHistoryPageData: makeSelectTransactionHistoryPageData(),
  isLoadingTransactionHistory: makeSelectIsLoadingTransactionHistoryPageData(),
});
let searchCalled = false;
const TransactionHistoryPage = props => {
  const {
    transactionHistoryPageData,
    isLoadingTransactionHistory,
  } = useSelector(stateSelector);
  const dispatch = useDispatch();
  const gridData: Transaction[] = transactionHistoryPageData
    ? transactionHistoryPageData
    : [];

  useEffect(() => {
    if (gridData.length === 0) {
      requestAnimationFrame(() => {
        dispatch(getTransactionHistoryPageDataAction());
      });
    } else {
      requestAnimationFrame(() => {
        //@ts-ignore
        dispatch(getTransactionHistoryPageDataAction({ silent: true }));
      });
    }
    return () => {};
  }, []);

  useNewPaymentEvent({
    toRunAfterNewEvent: () => {
      //@ts-ignore
      dispatch(getTransactionHistoryPageDataAction({ silent: true }));
    },
  });

  const onSearchReset = () => {
    if (searchCalled === true) {
      dispatch(getTransactionHistoryPageDataAction());
      searchCalled = false;
    }
  };
  const onSearchClick = (data: any) => {
    searchCalled = true;
    MessageService.send({
      name: MessageNames.SET_GRID_DATA,
      payload: [],
    });
    dispatch(getTransactionHistoryPageDataAction(data));
  };
  return (
    <div style={{ overflowX: 'auto', overflowY: 'hidden' }}>
      <AnimatedChildren isLoading={isLoadingTransactionHistory}>
        <MainWrapper>
          <TitledComponent
            id='transactionHistoryWrapper'
            title={<FormattedMessage {...translate.Exchangeaccoount} />}
          >
            <GridFilters
              hideTick
              TimePeriod={false}
              CurrencyPair={false}
              BuySell={false}
              Coin
              DWType
              Address
              onSearchClick={onSearchClick}
              onCancelClick={onSearchReset}
            />
            <DataGrid data={gridData} />
          </TitledComponent>
        </MainWrapper>
      </AnimatedChildren>
    </div>
  );
};

export default memo(TransactionHistoryPage);
const MainWrapper = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 200px);
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 100vw;
  min-width: 1180px;
`;
