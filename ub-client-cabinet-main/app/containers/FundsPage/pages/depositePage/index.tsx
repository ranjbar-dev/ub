import React, { memo, useEffect } from 'react';
import { createStructuredSelector } from 'reselect';

import { useSelector, useDispatch } from 'react-redux';
import AnimatedChildren from 'components/AnimateChildren';
import TitledComponent from 'components/titled';
import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';
import {
  makeSelectdepositAndWithDrawData,
  makeSelectIsLoadingdepositAndWithDrawData,
} from 'containers/FundsPage/selectors';
import { getDepositAndWithDrawDataAction } from 'containers/FundsPage/actions';
import DepositeAddress from './components/depositeAddress';
import { LocalStorageKeys } from 'services/constants';
import DataGrid from 'containers/FundsPage/components/common/dataGrid';

const stateSelector = createStructuredSelector({
  depositAndWithDrawData: makeSelectdepositAndWithDrawData(),
  isLoadingDepositeAndWithdraw: makeSelectIsLoadingdepositAndWithDrawData(),
});

const DepositePage = props => {
  const { depositAndWithDrawData, isLoadingDepositeAndWithdraw } = useSelector(
    stateSelector,
  );
  const dispatch = useDispatch();
  useEffect(() => {
    if (localStorage[LocalStorageKeys.SELECTED_COIN]) {
      dispatch(
        getDepositAndWithDrawDataAction({
          code: localStorage[LocalStorageKeys.SELECTED_COIN],
          type: 'deposit',
        }),
      );
      return;
    }
    dispatch(getDepositAndWithDrawDataAction({ code: 'BTC', type: 'deposit' }));

    return () => {
      localStorage.removeItem(LocalStorageKeys.FUND_PAGE);
    };
  }, []);

  return (
    <Wrapper>
      <AnimatedChildren isLoading={isLoadingDepositeAndWithdraw}>
        <HalfWrapper1>
          <TitledComponent
            title={<FormattedMessage {...translate.DEPOSITSADDRESS} />}
          >
            <DepositeAddress />
          </TitledComponent>
        </HalfWrapper1>
        <HalfWrapper2>
          <TitledComponent
            id='depositAndWithdrawGridWrapper'
            title={<FormattedMessage {...translate.DEPOSITSHISTORY} />}
          >
            <DataGrid
              coinCode={localStorage[LocalStorageKeys.SELECTED_COIN] ?? 'BTC'}
              data={
                depositAndWithDrawData.depositTransactions
                  ? depositAndWithDrawData.depositTransactions
                  : []
              }
              sectionName='deposit'
            />
          </TitledComponent>
        </HalfWrapper2>
      </AnimatedChildren>
    </Wrapper>
  );
};

export default memo(DepositePage);
const HalfWrapper1 = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 200px);
  min-height: 460px;
  margin-top: 0vh;
  /* width: 658px; */
  display: flex;
  flex-direction: column;
  align-items: center;
`;
const HalfWrapper2 = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 200px);
  min-height: 460px;
  margin-top: 0vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 580px;
`;
const Wrapper = styled.div`
  display: flex;
  max-width: 100vw;
  overflow-x: auto;
  .animm0 {
    flex: 70;
  }
  & > .animm0:first-child {
    flex: 45;
    margin-right: 12px;
  }
  /* justify-content: space-between; */
`;
