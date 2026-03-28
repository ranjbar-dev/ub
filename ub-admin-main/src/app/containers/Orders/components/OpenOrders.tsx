import { ColDef, ValueFormatterParams } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React, { memo } from 'react';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { CurrencyFormater } from 'utils/formatters';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { OrdersActions } from '../slice';
import { selectOpenOrdersData } from '../selectors';


interface Props {
  data: InitialUserDetails;
}
function OpenOrders(props: Props) {
  const { data } = props;
  const { t } = useTranslation();
  const openOrdersData = useSelector(selectOpenOrdersData);

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
      field: 'price',
      valueFormatter: (params: ValueFormatterParams) => {
        return CurrencyFormater(params.data.price);
      },
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
      headerName: t(translations.CommonTitles.Total()),
      field: 'total',
      valueFormatter: (params: ValueFormatterParams) => {
        return CurrencyFormater(params.data.total);
      },
    },
  ];
  const filters: FilterArrayElement = {
    //dateCols: ['createdAt'],
    dropDownCols: [
      {
        id: 'side',
        substituteId: 'side',
        options: [
          {
            name: 'Buy',
            value: 'buy',
          },
          {
            name: 'Sell',
            value: 'sell',
          },
        ],
      },
    ],
  };
  return (
    <div style={{ width: '100%' }}>
      <SimpleGrid
        containerId="UserDetailsWindow"
        additionalInitialParams={{ user_id: data.id }}
        arrayFieldName="orders"
        immutableId="id"
        filters={filters}
        userId={data.id}
        //onRowClick={handleRowClick}
        initialAction={OrdersActions.GetOpenOrdersAction}
        externalData={openOrdersData}
        staticRows={staticRows}
      />
    </div>
  );
}

export default memo(OpenOrders);
