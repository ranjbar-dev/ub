/**
 *
 * Billing
 *
 */

import AccountBalanceIcon from '@material-ui/icons/AccountBalance';
import AssignmentOutlinedIcon from '@material-ui/icons/AssignmentOutlined';
import CallMadeOutlinedIcon from '@material-ui/icons/CallMadeOutlined';
import LocalAtmIcon from '@material-ui/icons/LocalAtm';
import { ColDef, ValueFormatterParams, ICellRendererParams } from 'ag-grid-community';
import { CellRenderer } from 'app/components/renderer';
import UserDetailsTabs from 'app/containers/UserDetails/components';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React, { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { LocalStorageKeys } from 'services/constants';
import styled from 'styled-components/macro';
import { CurrencyFormater } from 'utils/formatters';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { billingSaga } from './saga';
import { BillingReducer, sliceKey } from './slice';
import { InitialUserDetails } from '../UserAccounts/types';
import BillingAllTransactionsDataGrid from './components/BillingAllTransactionsDataGrid';
import BillingDepositsDataGrid from './components/BillingDepositsDataGrid';
import BillingWithdrawalssDataGrid from './components/BillingWithdrawalssDataGrid';

interface Props {
  initialData: InitialUserDetails;
}

export const Billing = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: BillingReducer });
  useInjectSaga({ key: sliceKey, saga: billingSaga });
  const { initialData } = props;
  const { t } = useTranslation();
  const staticRows: ColDef[] = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.IDNO()),
        field: 'id',
      },
      {
        headerName: t(translations.CommonTitles.Method()),
        field: 'currencyCode',
      },
      {
        headerName: t(translations.CommonTitles.Amount()),
        field: 'amount',
        valueFormatter: (params: ValueFormatterParams) => {
          return CurrencyFormater(params.data.amount);
        },
      },
      {
        headerName: t(translations.CommonTitles.FromAddress()),
        field: 'fromAddress',
      },
      {
        headerName: t(translations.CommonTitles.ToAddress()),
        field: 'toAddress',
      },
      {
        headerName: t(translations.CommonTitles.TransactionId()),
        field: 'txId',
      },
      {
        headerName: t(translations.CommonTitles.CreationDate()),
        field: 'createdAt',
      },
      {
        headerName: t(translations.CommonTitles.LastUpdate()),
        field: 'updatedAt',
      },
      {
        headerName: t(translations.CommonTitles.Status()),
        field: 'status',
        ...cellColorAndNameFormatter('status'),
      },
      {
        width: 0,
        field: 'loading',
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(
            <>
              {data.IsLoading === true ? (
                <div
                  className="loadingGridRow"
                  id={'BillingLoading' + data.id}
                ></div>
              ) : (
                ''
              )}
            </>,
          ),
      },
    ],
    [],
  );
  const staticAllTransactionsRows: ColDef[] = useMemo(
    () => [
      {
        headerName: t(translations.CommonTitles.IDNO()),
        field: 'id',
      },
      {
        headerName: t(translations.CommonTitles.Method()),
        field: 'currencyCode',
      },
      {
        headerName: t(translations.CommonTitles.Type()),
        field: 'type',
      },
      {
        headerName: t(translations.CommonTitles.Amount()),
        field: 'amount',
        valueFormatter: (params: ValueFormatterParams) => {
          return CurrencyFormater(params.data.amount);
        },
      },
      //  {
      //    headerName: t(translations.CommonTitles.FromAddress()),
      //    field: 'fromAddress',
      //  },
      {
        headerName: t(translations.CommonTitles.ToAddress()),
        field: 'toAddress',
      },
      //  {
      //    headerName: t(translations.CommonTitles.TransactionId()),
      //    field: 'txId',
      //  },
      {
        headerName: t(translations.CommonTitles.CreationDate()),
        field: 'createdAt',
      },
      {
        headerName: t(translations.CommonTitles.LastUpdate()),
        field: 'updatedAt',
      },
      {
        headerName: t(translations.CommonTitles.Status()),
        field: 'status',
        ...cellColorAndNameFormatter('status'),
      },
      {
        width: 0,
        field: 'loading',
        cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
          CellRenderer(
            <>
              {data.IsLoading === true ? (
                <div
                  className="loadingGridRow"
                  id={'BillingLoading' + data.id}
                ></div>
              ) : (
                ''
              )}
            </>,
          ),
      },
    ],
    [],
  );
  const filters: FilterArrayElement = {
    dropDownCols: [
      {
        id: 'status',
        options: [
          {
            name: 'Created',
            value: 'created',
          },
          {
            name: 'In Progress',
            value: 'in_progress',
          },
          {
            name: 'Completed',
            value: 'completed',
          },

          {
            name: 'Failed',
            value: 'failed',
          },

          {
            name: 'Canceled',
            value: 'cancel',
          },
          {
            name: 'Rejected',
            value: 'reject',
          },
        ],
      },
      {
        id: 'currencyCode',
        options: JSON.parse(
          localStorage[LocalStorageKeys.CURRENCIES],
        ).currencies.map((item: { name: string; code: string }, index: number) => {
          return { name: item.name, value: item.code };
        }),
      },
    ],
    //dateCols: ['createdAt'],
  };
  return (
    <>
      <Wrapper>
        <UserDetailsTabs
          options={[
            {
              title: t(translations.CommonTitles.Deposits()),
              component: (
                <BillingDepositsDataGrid
                  data={initialData}
                  staticRows={staticRows}
                  filters={filters}
                />
              ),
              icon: <AccountBalanceIcon />,
            },
            {
              title: t(translations.CommonTitles.Withdrawals()),
              component: (
                <BillingWithdrawalssDataGrid
                  data={initialData}
                  staticRows={staticRows}
                  filters={filters}
                />
              ),
              icon: <LocalAtmIcon />,
            },
            {
              title: t(translations.CommonTitles.AllTransactions()),
              component: (
                <BillingAllTransactionsDataGrid
                  data={initialData}
                  staticRows={staticAllTransactionsRows}
                  filters={filters}
                />
              ),
              icon: <AssignmentOutlinedIcon />,
            },
          ]}
        />
      </Wrapper>
    </>
  );
});

const Wrapper = styled.div``;
