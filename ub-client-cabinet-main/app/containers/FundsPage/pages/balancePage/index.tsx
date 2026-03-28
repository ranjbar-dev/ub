import React, { memo, useEffect, useRef } from 'react';
import Filters from './components/filters';
import { createStructuredSelector } from 'reselect';
import {
  makeSelectBalancePageData,
  makeSelectIsLoadingBalancePageData,
} from 'containers/FundsPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import { Balance } from 'containers/FundsPage/types';
import { getBalancePageDataAction } from 'containers/FundsPage/actions';
import AnimatedChildren from 'components/AnimateChildren';
import DataGrid from './components/dataGrid';
import TitledComponent from 'components/titled';
import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';

const stateSelector = createStructuredSelector({
  balancePageData: makeSelectBalancePageData(),
  isLoadingBalances: makeSelectIsLoadingBalancePageData(),
});

const BalancePage = props => {
  const { balancePageData, isLoadingBalances } = useSelector(stateSelector);
  const dispatch = useDispatch();
  const gridData: Balance[] = balancePageData.balances ?? [];
  useEffect(() => {
    requestAnimationFrame(() => {
      dispatch(
        getBalancePageDataAction({ isSilent: !!balancePageData.balances }),
      );
    });

    return () => {};
  }, []);

  return (
    <AnimatedChildren isLoading={isLoadingBalances}>
      <Filters minimumForSwithFilter={balancePageData.minimumOfSmallBalances} />
      <ScrollWrapper>
        <MainWrapper>
          <TitledComponent
            id='balancePageWrapper'
            title={<FormattedMessage {...translate.Exchangeaccoount} />}
          >
            <DataGrid data={[...gridData]} />
          </TitledComponent>
        </MainWrapper>
      </ScrollWrapper>
    </AnimatedChildren>
  );
};

export default memo(BalancePage, () => true);

const ScrollWrapper = styled.div``;

const MainWrapper = styled(Card)`
  min-width: 1000px;
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(95.5vh - 225px);
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  /* min-width: 1180px; */
`;
