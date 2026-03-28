import {RowClickedEvent,CellClickedEvent,ValueFormatterParams,ICellRendererParams} from 'ag-grid-community';
import {CellRenderer} from 'app/components/renderer';
import {SimpleGrid} from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import {FullWidthWrapper} from 'app/components/wrappers/FullWidthWrapper';
import {WindowTypes} from 'app/constants';
import {translations} from 'locales/i18n';
import {FilterArrayElement} from 'locales/types';
import React,{memo,useCallback,useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {useDispatch} from 'react-redux';
import {LocalStorageKeys} from 'services/constants';
import {MessageNames,GridNames} from 'services/messageService';
import {cellColorAndNameFormatter} from 'utils/stylers';

import {UserAccountsActions} from '../slice';


interface Props { }
function VerificationPage(props: Props) {
	const {t}=useTranslation();
	const dispatch=useDispatch();
	const handleRowClick=useCallback(
		(e: RowClickedEvent) => {
			dispatch(
				UserAccountsActions.getInitialSingleUserDataAndOpenWindowAction({
					id: e.data.id,
					windowType: WindowTypes.Verification,
				}),
			);
		},
		[dispatch],
	);
	const staticRows=useMemo(
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
				valueFormatter: (params: ValueFormatterParams) => {
return params.data.countryFullName;
				},
			},
			{
				headerName: t(translations.Grid.IPAddress()),
				field: 'registeredIP',
			},
			{
				headerName: t(translations.Grid.CreationDate()),
				field: 'registrationDate',
			},
			{
				headerName: t(translations.Grid.Identity()),
				field: 'identityConfirmationStatus',
				...cellColorAndNameFormatter('identityConfirmationStatus'),
			},
			{
				headerName: t(translations.Grid.Address()),
				field: 'addressConfirmationStatus',
				...cellColorAndNameFormatter('addressConfirmationStatus'),
			},
			{
				headerName: t(translations.Grid.Phone()),
				field: 'phoneConfirmationStatus',
				...cellColorAndNameFormatter('phoneConfirmationStatus'),
			},
			{
				width: 0,
				field: 'loading',
				cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
					CellRenderer(
						<>
							<div
								className="loadingGridRow"
								id={'loading'+data.id}
								style={{display: 'none'}}
							></div>
						</>,
					),
			},
		],

		[t],
	);
	const options=useMemo(
		() => [
			{
				name: 'None',
				value: 'none',
			},
			{
				name: t(translations.CommonTitles.Incomplete()),
				value: 'incomplete',
			},
			{
				name: t(translations.CommonTitles.Confirmed()),
				value: 'confirmed',
			},
			{
				name: t(translations.CommonTitles.NotConfirmed()),
				value: 'not_confirmed',
			},
			{
				name: t(translations.CommonTitles.Rejected()),
				value: 'rejected',
			},
		],
		[t],
	);
	const dropDownCols=useMemo(
		() => [
			{
				id: 'identityConfirmationStatus',
				options,
			},
			{
				id: 'addressConfirmationStatus',
				options,
			},
			{
				id: 'phoneConfirmationStatus',
				options:options.filter((item,index)=>index!=0),
			},
		],

		[options],
	);
	const onCellClick=(e: CellClickedEvent) => {
		localStorage[LocalStorageKeys.VERIFICATION_WINDOW_TYPE]=
			e.colDef.headerName;
	};

	const filters: FilterArrayElement={
		dropDownCols,
		sortableCols: ['email'],
		countryCols: ['country'],
		dateCols: ['registrationDate'],
	}

	return (
		<FullWidthWrapper>
			<TitledContainer
				id="Verification"
				title={t(translations.CommonTitles.UserVerification())}
			>
				<SimpleGrid
					containerId="Verification"
					additionalInitialParams={{need_verification: 1}}
					arrayFieldName="users"
					immutableId="id"
					flashCellUpdate={true}
					gridName={GridNames.USER_VERIFICATION}
					filters={filters}
					onRowClick={handleRowClick}
					onCellClick={onCellClick}
					initialAction={UserAccountsActions.GetInitialUserAccountsAction}
					messageName={MessageNames.SET_USER_ACCOUNTS}
					staticRows={staticRows}
				/>
			</TitledContainer>
		</FullWidthWrapper>
	);
}

export default memo(VerificationPage);
