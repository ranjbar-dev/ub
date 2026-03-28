/*
 *
 * FundsPage
 *
 */

import React, { memo, useState, useEffect } from 'react';
import { Helmet } from 'react-helmet';
import { FormattedMessage } from 'react-intl';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import {
  makeSelectBalancePageData,
  makeSelectLocation,
  makeSelectUserData,
  //  makeSelectUserData,
} from './selectors';
import reducer from './reducer';
import saga from './saga';
import translate from './messages';
import { MaxWidthWrapper1600 } from 'components/wrappers/maxWidthWrapper1600';
import SegmentSelector from './components/segmentSelector';
import { FundsPages } from './constants';
import TopInfo from './components/topInfo';
import BalancePage from './pages/balancePage';
import DepositePage from './pages/depositePage/Loadable';
import WithdrawalsPage from './pages/withdrawalsPage/Loadable';
import TransactionsHistoryPage from './pages/transactionsHistoryPage';
import { LocalStorageKeys } from 'services/constants';
import { getBalancePageDataAction, getUserDataAction } from './actions';
import { RegisteredUserSubscriber } from 'services/message_service';
import { toast } from 'components/Customized/react-toastify';
import { AppPages } from 'containers/App/constants';
import { storage } from 'utils/storage';
import { GridLoading } from 'components/grid_loading/gridLoading';
import styled from 'styles/styled-components';
import EmailNotVerifiedIcon from 'images/themedIcons/emailNotVerifiedIcon';
import { useNewPaymentEvent } from './hooks/useNewPaymentEvent';

const stateSelector = createStructuredSelector({
  balancePageData: makeSelectBalancePageData(),
  userData: makeSelectUserData(),
  location: makeSelectLocation(),
});

interface Props {}

const pageSelector = (ActivePage: FundsPages) => {
  switch (ActivePage) {
    case FundsPages.BALANCE:
      return <BalancePage />;
    case FundsPages.DEPOSIT:
      return <DepositePage />;
    case FundsPages.WITHDRAWALS:
      return <WithdrawalsPage />;
    case FundsPages.TRANSACTION_HISTORY:
      return <TransactionsHistoryPage />;
    default:
      return <BalancePage />;
  }
};

function FundsPage (props: Props) {
  useInjectReducer({ key: 'fundsPage', reducer: reducer });
  useInjectSaga({ key: 'fundsPage', saga: saga });
  const dispatch = useDispatch();
  const { balancePageData, location, userData } = useSelector(stateSelector);
  useEffect(() => {
    //@ts-ignore
    if (userData?.isAccountVerified === undefined) {
      dispatch(getUserDataAction());
    }
    dispatch(getBalancePageDataAction({ isSilent: false }));

    return () => {};
  }, []);

  const [ActivePage, setActivePage] = useState(
    localStorage[LocalStorageKeys.FUND_PAGE] === FundsPages.DEPOSIT
      ? FundsPages.DEPOSIT
      : localStorage[LocalStorageKeys.FUND_PAGE] === FundsPages.WITHDRAWALS
      ? FundsPages.WITHDRAWALS
      : FundsPages.BALANCE,
  );
  //@ts-ignore
  const subPage = location.pathname.split(AppPages.Funds)[1]?.replace('/', '');
  useEffect(() => {
    if (subPage === FundsPages.DEPOSIT) {
      document.getElementById('depositeTab')!.click();
    } else if (subPage === FundsPages.WITHDRAWALS) {
      document.getElementById('withdrawTab')!.click();
    }

    return () => {};
  }, [subPage]);

  useNewPaymentEvent({
    toRunAfterNewEvent: () => {
      dispatch(getBalancePageDataAction({ isSilent: true }));
    },
  });

  return (
    <>
      <Helmet>
        <title>Funds </title>
        <meta name='description' content='Description of FundsPage' />
      </Helmet>

      <MaxWidthWrapper1600>
        <SegmentSelector
          defaultIndex={
            localStorage[LocalStorageKeys.FUND_PAGE] === FundsPages.DEPOSIT
              ? 1
              : localStorage[LocalStorageKeys.FUND_PAGE] ===
                FundsPages.WITHDRAWALS
              ? 2
              : 0
          }
          onChange={(page: FundsPages) => {
            setActivePage(page);
            //if (page === FundsPages.BALANCE) {
            //  dispatch(getBalancePageDataAction({ isSilent: true }));
            //}
          }}
          options={[
            {
              title: <FormattedMessage {...translate.balance} />,
              page: FundsPages.BALANCE,
            },
            {
              title: <FormattedMessage {...translate.deposite} />,
              page: FundsPages.DEPOSIT,
            },
            {
              title: <FormattedMessage {...translate.withdrawals} />,
              page: FundsPages.WITHDRAWALS,
            },
            {
              title: <FormattedMessage {...translate.transactionHistory} />,
              page: FundsPages.TRANSACTION_HISTORY,
            },
          ]}
        />
        <TopInfo
          data={{
            etimatedBalance: balancePageData.totalSum
              ? balancePageData.totalSum
              : '',
            availableBalance: balancePageData.availableSum
              ? balancePageData.availableSum
              : '',
            inOrder: balancePageData.inOrderSum
              ? balancePageData.inOrderSum
              : '',
            btcInOrder: balancePageData.btcInOrderSum
              ? balancePageData.btcInOrderSum
              : '',
            btcAvailableBalance: balancePageData.btcAvailableSum
              ? balancePageData.btcAvailableSum
              : '',
            btcEtimatedBalance: balancePageData.btcTotalSum
              ? balancePageData.btcTotalSum
              : '',
          }}
        />
        {userData?.isAccountVerified === true ? (
          pageSelector(ActivePage)
        ) : userData?.isAccountVerified === false ? (
          <MainWrapper>
            <EmailNotVerifiedIcon />
            <UnderText>
              Your email is not verified yet. Please check your inbox and
              confirm your Email
            </UnderText>
          </MainWrapper>
        ) : (
          <GridLoading />
        )}
      </MaxWidthWrapper1600>
    </>
  );
}
const UnderText = styled.p`
  color: #f49806;
  font-size: 14px;
  margin-top: -20px;
`;
const MainWrapper = styled.div`
  min-width: 1000px;
  background: var(--white);
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(95.5vh - 225px);
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  place-content: center;
  /* min-width: 1180px; */
`;

export default memo(FundsPage);
