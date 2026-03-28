import { ValueFormatterParams, ICellRendererParams } from 'ag-grid-community';
import { CellRenderer } from 'app/components/renderer';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React from 'react'
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { MessageNames } from 'services/messageService';
import { CurrencyFormater } from 'utils/formatters';

import { BalancesActions } from '../slice';
import { selectBalancesHistoryData } from '../selectors';
import { WalletTypes } from '../types';

interface Props { }

function TransferHistory(props: Props) {
	const { } = props
	const { t } = useTranslation();
	const balancesHistory = useSelector(selectBalancesHistoryData);

	/*
    
	amount: "0.02"
	code: "ETH"
	createdAt: "2020-11-29 08:36:48"
	from: "HOT"
	fromAddress: "0x916849c5ead1E3d40dDf0F63B10b1e522b395bC3"
	id: 4
	metadata: null
	status: "COMPLETED"
	to: "EXTERNAL"
	toAddress: "0x48e4c06ede2c08112cd1b767919c940b5a84c0d4"
	txId: "0xbd3a2ea91808cfcd28c4cf84498d24b1597326e83d197c4de05a74d7e2aee5f7"
	updatedAt: "2020-11-29 08:38:47"
    
	*/


	const staticRows = [
		{
			headerName: t(translations.CommonTitles.Date()),
			field: 'createdAt',
			maxWidth: 150
		},
		{
			headerName: t(translations.CommonTitles.Amount()),
			field: 'amount',
			valueFormatter: (params: ValueFormatterParams) => {
return CurrencyFormater(params.data.amount);
			},
			maxWidth: 150
		},
		{
			headerName: t(translations.CommonTitles.Coin()),
			field: 'code',
			maxWidth: 120
		},

		// {
		//     headerName: t(translations.CommonTitles.FromAddress()),
		//     field: 'fromAddress',
		//     minWidth: 360

		// },
		// {
		//     headerName: t(translations.CommonTitles.ToAddress()),
		//     field: 'toAddress',
		//     minWidth: 360
		// },
		{
			headerName: t(translations.CommonTitles.From()),
			field: 'from',
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<span className="bold">
							{`[${data.from}] - `}
						</span>
						<span>{data.fromAddress}</span>
					</>,
				),
		},
		{
			headerName: t(translations.CommonTitles.To()),
			field: 'to',
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<span className="bold">
							{`[${data.to ?? 'CUSTOM'}] - `}
						</span>
						<span>{data.toAddress}</span>
					</>,
				),
		},
		{
			headerName: t(translations.CommonTitles.Status()),
			field: 'status',
			maxWidth: 120,
			valueFormatter: (params: ValueFormatterParams) => {
return params.data.status.replace('_', ' ');
			},
		},
		{
			headerName: t(translations.CommonTitles.Info()),
			field: 'info',
			maxWidth: 120
		},
	]
	const handleDataRecieved = (data: unknown) => {

	}

	const filters: FilterArrayElement = {
		disabledCols: ['createdAt', 'amount', 'code', 'fromAddress', 'toAddress', 'info'],
		dropDownCols: [
			{
				id: 'from',
				options: [

					{
						name: 'Hot Wallet',
						value: WalletTypes.Hot
					},


					{
						name: 'Cold Wallet',
						value: WalletTypes.Cold
					},


					{
						name: 'Liquidity Wallet',
						value: WalletTypes.External
					},
				]
			},

			{
				id: 'to',
				options: [

					{
						name: 'Hot Wallet',
						value: WalletTypes.Hot
					},


					{
						name: 'Cold Wallet',
						value: WalletTypes.Cold
					},


					{
						name: 'Liquidity Wallet',
						value: WalletTypes.External
					},
				]
			},

			{
				id: 'status',
				options: [

					{
						name: 'Created',
						value: 'created'
					},

					{
						name: 'Completed',
						value: 'completed'
					},
					{
						name: 'In Progress',
						value: 'in_progress'
					},
					{
						name: 'Failed',
						value: 'failed'
					},
					{
						name: 'Canceled',
						value: 'canceled'
					},
					{
						name: 'Rejected',
						value: 'rejected'
					},

				]
			}

		]
	}
	return (
		<SimpleGrid
			//topTabs={filterTabs}
			containerId="balances"
			additionalInitialParams={{}}
			arrayFieldName="internalTransfers"
			immutableId="id"
			onDataReceived={handleDataRecieved}
			filters={filters}
			//  onRowClick={handleRowClick}
			additionalHeight={0}
			pullUpPaginationBy={32}
			initialAction={BalancesActions.GetBalanceHistoryAction}
			messageName={MessageNames.SET_BALANCES_HISTORY_DATA}
			externalData={balancesHistory}
			staticRows={staticRows}
		/>
	)
}

export default TransferHistory
