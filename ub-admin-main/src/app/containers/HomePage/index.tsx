/**
 *
 * HomePage
 *
 */

import { ICellRendererParams } from 'ag-grid-community';
import { CellRenderer } from 'app/components/renderer';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';
import useOpenWithdrawWindow from 'utils/hooks/useOpenWithdrawWindow';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import DashboardCard from './components/DashboardCard';
import SearchInput from './components/SearchInput';
import { homePageSaga } from './saga';
import { HomePageReducer, sliceKey, HomePageActions } from './slice';


interface Props {}

export const HomePage = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: HomePageReducer });
  useInjectSaga({ key: sliceKey, saga: homePageSaga });
  useOpenWithdrawWindow({});
  const dispatch = useDispatch();
  const { t } = useTranslation();
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {});
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  const staticRows = useMemo(
    () => [
      {
        headerName: t(translations.Grid.IDNO()),
        field: 'id',
      },
      {
        headerName: t(translations.Grid.EmailAddress()),
        field: 'email',
      },
      {
        headerName: t(translations.Grid.FirstName()),
        field: 'firstName',
      },
      {
        headerName: t(translations.Grid.LastName()),
        field: 'lastName',
      },
      {
        headerName: t(translations.Grid.Country()),
        field: 'country',
      },
      {
        headerName: t(translations.Grid.ReferKey()),
        field: 'referKey',
      },
      {
        headerName: t(translations.Grid.ReferalID()),
        field: 'referralId',
      },
      {
        headerName: t(translations.Grid.Manager()),
        field: 'manager',
      },
      {
        headerName: t(translations.Grid.RegisterDate()),
        field: 'registrationDate',
      },
      {
        headerName: t(translations.Grid.RegisterIP()),
        field: 'registeredIP',
      },
      {
        width: 0,
        field: 'loading',
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(
            <>
              <div
                className="loadingGridRow"
                id={'loading' + data.id}
                style={{ display: 'none' }}
              ></div>
            </>,
          ),
      },
    ],

    [t],
  );
  const handleEnter = (data: { field: string; value: string }) => {
    if (data.field === 'email') {
      dispatch(HomePageActions.getUserByIdAction({ email: data.value }));
      return;
    }
    if (data.field === 'systemId') {
      dispatch(HomePageActions.getUserByIdAction({ id: data.value }));
      return;
    }
    if (data.field === 'withdrawalId') {
      dispatch(HomePageActions.getWithdrawalByIdAction({ id: data.value }));
      return;
    }
  };

  return (
    <FullWidthWrapper className="noAlignment">
      <div className="gridPlaceHolder" style={{ display: 'none' }}>
        <SimpleGrid
          containerId="userAccounts"
          additionalInitialParams={{}}
          arrayFieldName="users"
          immutableId="id"
          filters={{
            countryCols: ['country'],
            dateCols: ['registrationDate'],
          }}
          //  onRowClick={handleRowClick}
          initialAction={HomePageActions.someAction}
          messageName={MessageNames.SET_USER_ACCOUNTS}
          staticRows={staticRows}
        />
      </div>
      <DashboardCard title={t(translations.CommonTitles.QuickAccess())}>
        <SearchInput
          title={t(translations.CommonTitles.UserSystemId())}
          onEnter={value => handleEnter({ field: 'systemId', value })}
        />
        <SearchInput
          title={t(translations.CommonTitles.UserEmail())}
          onEnter={value => handleEnter({ field: 'email', value })}
        />
        <SearchInput
          title={t(translations.CommonTitles.WithdrawalID())}
          onEnter={value => handleEnter({ field: 'withdrawalId', value })}
        />
        <SearchInput
          isLast={true}
          disabled={true}
          title={t(translations.CommonTitles.DepositID())}
          onEnter={value => handleEnter({ field: 'depositId', value })}
        />
      </DashboardCard>
    </FullWidthWrapper>
  );
});
