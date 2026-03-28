import {ColDef} from 'ag-grid-community';
import {SimpleGrid} from 'app/components/SimpleGrid/SimpleGrid';
import {InitialUserDetails} from 'app/containers/UserAccounts/types';
import {translations} from 'locales/i18n';
import React,{memo,useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {MessageNames} from 'services/messageService';

import {UserDetailsActions} from '../slice';

//import { Balance, UserBalances } from '../types';

interface Props {
	data: InitialUserDetails;
}
function WhiteAddressesSegemnt(props: Props) {
	const {data}=props;
	const {t}=useTranslation();
	const staticRows: ColDef[]=useMemo(
		() => [
			{
				headerName: t(translations.Grid.Name()),
				field: 'name',
				suppressMenu: true,
				sortable: true,
			},
			{
				headerName: t(translations.Grid.Currency()),
				field: 'code',
				suppressMenu: true,
				sortable: true,
			},
			{
				headerName: t(translations.Grid.Address()),
				field: 'address',
				suppressMenu: true,
				sortable: true,
			},
			{
				headerName: t(translations.Grid.Label()),
				field: 'label',
				suppressMenu: true,
				sortable: true,
			},
		],
		[],
	);

	return (
		<div style={{width: '100%'}}>
			<SimpleGrid
				containerId="UserDetailsWindow"
				additionalInitialParams={{id: data.id}}
				arrayFieldName="addresses"
				immutableId="name"
				//filters={{}}
				//onRowClick={handleRowClick}
				initialAction={UserDetailsActions.GetWhiteAddressesAction}
				messageName={MessageNames.SET_WHITEADDRESSES_DATA}
				staticRows={staticRows}
			/>
		</div>
	);
}

export default memo(WhiteAddressesSegemnt);
