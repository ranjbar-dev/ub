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
import { selectTradeHistoryData } from '../selectors';


interface Props {
  data: InitialUserDetails;
}
function TradeHistory(props: Props) {
  const { data } = props;
  const { t } = useTranslation();
  const tradeHistoryData = useSelector(selectTradeHistoryData);

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
      ...cellColorAndNameFormatter('type'),
    },
    {
      headerName: t(translations.CommonTitles.Price()),
      field: 'price',
      valueFormatter: (params: ValueFormatterParams) => {
        return CurrencyFormater(params.data.price);
      },
    },
    {
      headerName: t(translations.CommonTitles.Amount()),
      field: 'amount',
      minWidth: 180,
      valueFormatter: (params: ValueFormatterParams) => {
        return CurrencyFormater(params.data.amount);
      },
    },
    {
      headerName: t(translations.CommonTitles.Exchange()),
      field: 'exchange',
    },
    //{
    //  headerName: t(translations.CommonTitles.Status()),
    //  field: 'status',
    //  ...cellColorAndNameFormatter('status'),
    //},
  ];

  return (
    <div style={{ width: '100%' }}>
      <SimpleGrid
        containerId="UserDetailsWindow"
        additionalInitialParams={{ user_id: data.id }}
        arrayFieldName="orders"
        immutableId="id"
        filters={{}}
        userId={data.id}
        //onRowClick={handleRowClick}
        initialAction={OrdersActions.GetTradeHistoryAction}
        externalData={tradeHistoryData}
        staticRows={staticRows}
      />
    </div>
  );
}

export default memo(TradeHistory);
