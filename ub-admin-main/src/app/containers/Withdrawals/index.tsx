/**
 *
 * Withdrawals
 *
 */


import { ColDef, ValueFormatterParams, ICellRendererParams, RowClickedEvent } from 'ag-grid-community';
import GridTabs from 'app/components/GridTabs/GridTabs';
import { CellRenderer } from 'app/components/renderer';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React, { memo, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { LocalStorageKeys } from 'services/constants';
import {
	GridNames,
	MessageNames,
	MessageService,
} from 'services/messageService';
import styled from 'styled-components/macro';
import { CurrencyFormater } from 'utils/formatters';
import useOpenWithdrawWindow from 'utils/hooks/useOpenWithdrawWindow';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { withdrawalsSaga } from './saga';
import { WithdrawalsActions, WithdrawalsReducer, sliceKey } from './slice';
import { selectWithdrawalsData } from './selectors';

interface Props { }

export const Withdrawals = memo((props: Props) => {
	useInjectReducer({ key: sliceKey, reducer: WithdrawalsReducer });
	useInjectSaga({ key: sliceKey, saga: withdrawalsSaga });
	useOpenWithdrawWindow({});
	const [AdminStatus, setAdminStatus] = useState('pending');
	const { t } = useTranslation();

	const staticRows: ColDef[] = useMemo(
		() => [
			{
				headerName: t(translations.CommonTitles.ID()),
				field: 'id',
				maxWidth: 120,
			},
			{
				headerName: t(translations.CommonTitles.UserId()),
				field: 'userId',
				maxWidth: 120,
			},
			{
				headerName: t(translations.CommonTitles.UserEmail()),
				field: 'userEmail',
				maxWidth: 300,
			},

			{
				headerName: t(translations.CommonTitles.Method()),
				field: 'currencyCode',
				maxWidth: 120,
				valueFormatter: (params: ValueFormatterParams) => {
					return (params.data.currencyCode) + (params.data.blockchainNetwork ? `(${params.data.blockchainNetwork})` : '');
				},
			},
			{
				headerName: t(translations.CommonTitles.Amount()),
				field: 'amount',
				maxWidth: 190,
				valueFormatter: (params: ValueFormatterParams) => {
					return CurrencyFormater(params.data.amount);
				},
			},

			{
				headerName: t(translations.CommonTitles.ToAddress()),
				field: 'toAddress',
			},

			{
				headerName: t(translations.CommonTitles.CreationDate()),
				field: 'createdAt',
				maxWidth: 180,
			},

			{
				width: 0,
				hidden: true,
				field: 'loading',
				cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
					CellRenderer(
						<>
							<div
								className="loadingGridRow"
								id={'loading_main_withdrawals' + data.id}
								style={{ display: 'none' }}
							></div>
						</>,
					),
			},
		],

		[t],
	);

	const allStaticRows: ColDef[] = useMemo(
		() => [
			{
				headerName: t(translations.CommonTitles.ID()),
				field: 'id',
				maxWidth: 100,
			},
			{
				headerName: t(translations.CommonTitles.Email()),
				field: 'userEmail',
				width: 180
			},
			{
				headerName: t(translations.CommonTitles.Method()),
				field: 'currencyCode',
				maxWidth: 120,
				valueFormatter: (params: ValueFormatterParams) => {
					return (params.data.currencyCode) + (params.data.blockchainNetwork ? ` (${params.data.blockchainNetwork})` : '');
				},
			},
			{
				headerName: t(translations.CommonTitles.Amount()),
				field: 'amount',
				maxWidth: 150,
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
			//{
			//	headerName: t(translations.CommonTitles.TransactionId()),
			//	field: 'txId',
			//},
			{
				headerName: t(translations.CommonTitles.CreationDate()),
				field: 'createdAt',
				maxWidth: 130,
			},
			{
				headerName: t(translations.CommonTitles.LastUpdate()),
				field: 'updatedAt',
				maxWidth: 130,
			},
			{
				headerName: t(translations.CommonTitles.Status()),
				field: 'status',
				maxWidth: 140,
				...cellColorAndNameFormatter('status'),
			},
			{
				width: 0,
				field: 'loading',
				cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
					CellRenderer(
						<>
							<div
								className="loadingGridRow"
								id={'loading_main_withdrawals' + data.id}
								style={{ display: 'none' }}
							></div>
						</>,
					),
			},
		],

		[t],
	);
	const filters: FilterArrayElement = {
		dateCols: ['createdAt', 'updatedAt'],
		hiddenCols: ['loading'],

		dropDownCols: [
			{
				id: 'status',
				options: [
					{
						name: 'Created',
						value: 'created',
					},
					{
						name: 'Pending',
						value: 'pending',
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
						name: 'User Canceled',
						value: 'user_canceled',
					},
					{
						name: 'Canceled',
						value: 'canceled',
					},
					{
						name: 'Rejected',
						value: 'rejected',
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
	};

	const filterTabs = [
		{
			name: t(translations.CommonTitles.Pending()),
			callObject: {
				///when tab is changed

				admin_status: 'pending',
			},
		},
		{
			name: t(translations.CommonTitles.RecheckRequired()),
			callObject: {
				admin_status: 'recheck',
			},
		},
		{
			name: t(translations.CommonTitles.Accepted()),
			callObject: {
				admin_status: 'approved',
			},
		},
		{
			name: t(translations.CommonTitles.All()),
			callObject: {
				admin_status: null,
			},
		},
	];

	const dispatch = useDispatch();
	const withdrawalsData = useSelector(selectWithdrawalsData);
	const handleRowClick = useCallback(
		(e: RowClickedEvent) => {
			//setInitialModalData(e);
			dispatch(WithdrawalsActions.GetWithdrawalDetailAction(e.data));
			//setIsModalOpen(true);
		},
		[dispatch],
	);

	const handleTopTabChange = (e: Record<string, unknown>) => {
		setAdminStatus(e.admin_status as string);
		MessageService.send({
			name: MessageNames.APPLY_PARAMS_TO_GRID,
			gridName: GridNames.MAIN_WITHDRAWALS,
			payload: e,
		});
	};
	return (
		<FullWidthWrapper>
			<TitledContainer
				id="withdrawals"
				title={t(translations.CommonTitles.Withdrawals())}
			>
				<GridTabs onChange={handleTopTabChange} tabs={filterTabs} />
				{!(AdminStatus == null) && (
					<SimpleGrid
						//  topTabs={filterTabs}
						containerId="withdrawals"
						gridName={GridNames.MAIN_WITHDRAWALS}
						additionalInitialParams={{
							type: 'withdraw',
							//initial get
							status: 'created',
							admin_status: 'pending',
						}}
						arrayFieldName="payments"
						immutableId="id"
						filters={filters}
						additionalHeight={0}
						pullUpPaginationBy={32}
						onRowClick={handleRowClick}
						initialAction={WithdrawalsActions.GetWithdrawals}
						messageName={MessageNames.SET_WITHDRAWALS_DATA}
						externalData={withdrawalsData}
						staticRows={staticRows}
					/>
				)}
				{AdminStatus == null && (
					<SimpleGrid
						//  topTabs={filterTabs}
						containerId="withdrawals"
						gridName={GridNames.MAIN_WITHDRAWALS}
						additionalInitialParams={{
							type: 'withdraw',
						}}
						arrayFieldName="payments"
						immutableId="id"
						filters={filters}
						additionalHeight={0}
						pullUpPaginationBy={32}
						onRowClick={handleRowClick}
						initialAction={WithdrawalsActions.GetWithdrawals}
						messageName={MessageNames.SET_WITHDRAWALS_DATA}
						externalData={withdrawalsData}
						staticRows={allStaticRows}
					/>
				)}
			</TitledContainer>
		</FullWidthWrapper>
	);
});

const Wrapper = styled.div``;
