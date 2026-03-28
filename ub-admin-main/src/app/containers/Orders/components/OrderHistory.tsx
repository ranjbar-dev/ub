import { ColDef, ValueFormatterParams } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { memo } from 'react';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { CurrencyFormater } from 'utils/formatters';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { OrdersActions } from '../slice';
import { selectOrderHistoryData } from '../selectors';


interface Props {
  data: InitialUserDetails;
}
function OrderHistory(props: Props) {
  const { data } = props;
  const { t } = useTranslation();
  const orderHistoryData = useSelector(selectOrderHistoryData);

  const staticRows: ColDef[] = [
    {
      headerName: t(translations.CommonTitles.ID()),
      field: 'id',
      maxWidth: 120,
    },
    {
      headerName: t(translations.CommonTitles.Date()),
      field: 'createdAt',
      minWidth: 160,
    },
    {
      headerName: t(translations.CommonTitles.Type()),
      field: 'side',
      ...cellColorAndNameFormatter('side'),
    },
    {
      headerName: t(translations.CommonTitles.Pair()),
      field: 'pair',
    },

    {
      headerName: t(translations.CommonTitles.Side()),
      field: 'type',
    },
    {
      headerName: t(translations.CommonTitles.Price()),
      field: 'averagePrice',
    },
    {
      headerName: t(translations.CommonTitles.Amount()),
      field: 'amount',
      minWidth: 160,
      valueFormatter: (params: ValueFormatterParams) => {
        return CurrencyFormater(params.data.amount);
      },
    },

    {
      headerName: t(translations.CommonTitles.Status()),
      field: 'status',
    },
  ];

  return (
    <div style={{ width: '100%' }}>
      <SimpleGrid
        containerId="UserDetailsWindow"
        additionalInitialParams={{ user_id: data.id, status: 'filled' }}
        arrayFieldName="orders"
        immutableId="id"
        filters={{}}
        userId={data.id}
        //onRowClick={handleRowClick}
        initialAction={OrdersActions.GetOrderHistoryAction}
        externalData={orderHistoryData}
        staticRows={staticRows}
      />
    </div>
  );
}

export default memo(OrderHistory);
