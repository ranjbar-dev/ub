import { ValueFormatterParams } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React from 'react'
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { LocalStorageKeys } from 'services/constants';
import { CurrencyFormater } from 'utils/formatters';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { selectAllQueueData } from '../selectors';
import { ExternalOrdersActions } from '../slice';

interface Props { }

function AllQueue(props: Props) {

	const { t } = useTranslation();
	const allQueueData = useSelector(selectAllQueueData);

	const staticRows = [

		{
			headerName: t(translations.CommonTitles.OrderID()),
			field: 'orderId',
		},
		{
			headerName: t(translations.CommonTitles.Date()),
			field: 'date',

		},
		{
			headerName: 'Email',
			field: 'username',
		},
		{
			headerName: t(translations.CommonTitles.UserId()),
			field: 'userId',
		},
		{
			headerName: t(translations.CommonTitles.Pair()),
			field: 'pairCurrencyName',
		},
		{
			headerName: t(translations.CommonTitles.Side()),
			field: 'type',
			maxWidth: 120,
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
			valueFormatter: (params: ValueFormatterParams) => {
return CurrencyFormater(params.data.amount);
			},
		},
	]

	const filters: FilterArrayElement = {
		disabledCols: ['orderId', 'userId', 'type', 'price', 'amount'],
		dropDownCols: [
			{
				id: 'pairCurrencyName',
				substituteId: 'pair_currency_id',
				options: JSON.parse(localStorage[LocalStorageKeys.PAIRS])
			}
		],
	};


	return (
		<SimpleGrid
			containerId="externalOrders"
			additionalInitialParams={{}}
			arrayFieldName="externalExchangeOrders"
			immutableId="orderId"
			filters={filters}
			//  onRowClick={handleRowClick}
			additionalHeight={0}
			pullUpPaginationBy={32}
			initialAction={ExternalOrdersActions.GetAllQueueAction}
			externalData={allQueueData}
			staticRows={staticRows}
		/>
	)
}

export default AllQueue
