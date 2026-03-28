import { ColDef, ValueFormatterParams } from 'ag-grid-community';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { MessageNames } from 'services/messageService';
import { CurrencyFormater } from 'utils/formatters';

import { UserDetailsActions } from '../slice';

interface Props {
	data: InitialUserDetails;
}
function WaletSegment(props: Props) {
	const { data } = props;

	const { t } = useTranslation();

	const staticRows: ColDef[] = [
		{
			headerName: t(translations.Grid.Name()),
			field: 'name',
			maxWidth: 100,
			minWidth: 100,
		},
		{
			headerName: t(translations.Grid.Currency()),
			field: 'code',
			maxWidth: 85,
			minWidth: 85,
		},
		{
			headerName: t(translations.Grid.Address()),
			field: 'address',

		},
		//{
		//	headerName: t(translations.Grid.CreationDate()),
		//	field: 'lastName',
		//},
		{
			headerName: t(translations.Grid.TotalBalance()),
			field: 'totalAmount',
			maxWidth: 140,
			minWidth: 140,
			valueFormatter: (params: ValueFormatterParams) => {
return CurrencyFormater(params.data.totalAmount);
			},
		},
		{
			headerName: t(translations.Grid.Available()),
			field: 'availableAmount',
			maxWidth: 140,
			minWidth: 140,
			valueFormatter: (params: ValueFormatterParams) => {
return CurrencyFormater(params.data.availableAmount);
			},
		},
		{
			headerName: t(translations.Grid.InOrder()),
			field: 'inOrderAmount',
			maxWidth: 140,
			minWidth: 140,
			valueFormatter: (params: ValueFormatterParams) => {
return CurrencyFormater(params.data.inOrderAmount);
			},
		},

		//{
		//	headerName: t(translations.Grid.TotalDeposit()),
		//	field: 'registrationDate',
		//},
		//{
		//	headerName: t(translations.Grid.TotalWithdraw()),
		//	field: 'registeredIP',
		//},
	];

	return (
		<div style={{ width: '100%' }}>
			<SimpleGrid
				containerId="UserDetailsWindow"
				additionalInitialParams={{ id: data.id }}
				arrayFieldName="balances"
				immutableId="name"
				//filters={{}}
				//onRowClick={handleRowClick}
				initialAction={UserDetailsActions.GetWalletsAction}
				messageName={MessageNames.SET_WALLETS_DATA}
				staticRows={staticRows}
			/>
		</div>
	);
}

export default memo(WaletSegment);
