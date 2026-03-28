import { ColDef, ICellRendererParams, ValueFormatterParams } from 'ag-grid-community';
import PopupModal from 'app/components/materialModal/modal';
import { CellRenderer } from 'app/components/renderer';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React, { useState } from 'react'
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { LocalStorageKeys } from 'services/constants';
import styled from 'styled-components';
import { CurrencyFormater } from 'utils/formatters';
import { cellColorAndNameFormatter } from 'utils/stylers';

import { selectExternalOrdersData } from '../selectors';
import { ExternalOrdersActions } from '../slice';

interface Props { }

function Executed(props: Props) {

	const { t } = useTranslation();
	const externalOrdersData = useSelector(selectExternalOrdersData);

	const [popupInfo, setPopupInfo] = useState({
		isOpen: false,
		data: [{ name: '', value: '' }]
	})

	const handleExpandClick = (data: Record<string, unknown>) => {

		setPopupInfo({
			isOpen: true, data: [
				{ name: "Exception message", value: data.exceptionMessage as string },
				{ name: "Fail reason", value: data.failReason as string }

			]
		})

	}

	const staticRows: ColDef[] = [
		{
			headerName: t(translations.CommonTitles.ID()),
			field: 'id',
			maxWidth: 120,
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
			headerName: t(translations.CommonTitles.Pair()),
			field: 'pairCurrencyName',
			maxWidth: 100,
			minWidth: 100,

		},
		{
			headerName: t(translations.CommonTitles.Type()),
			field: 'type',
			maxWidth: 80,
			minWidth: 80,
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
		{
			headerName: t(translations.CommonTitles.BuyAmount()),
			field: 'buyAmount',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.buyAmount);
			},
		},
		{
			headerName: t(translations.CommonTitles.SellAmount()),
			field: 'sellAmount',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.sellAmount);
			},
		},
		{
			headerName: t(translations.CommonTitles.Total()),
			field: 'total',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.total + '');
			},
		},
		{
			headerName: t(translations.CommonTitles.Fee()),
			field: 'fee',
			valueFormatter: (params: ValueFormatterParams) => {
				return params.data.externalExchangeOtherInfo ?
					(CurrencyFormater(params.data.externalExchangeOtherInfo.fee.cost + '') + ' ' + params.data.externalExchangeOtherInfo.fee.currency)
					: '-';
			},
		},

		// {
		{
			headerName: 'E.M',
			field: 'exceptionMessage',
			cellRenderer: (params: ICellRendererParams) => CellRenderer(<>
				{params.data.exceptionMessage ?
					<ExpandButton onClick={() => { handleExpandClick(params.data) }} > {params.data.exceptionMessage.substring(0, 8)} </ExpandButton>

					: '-'

				}
			</>)
		},
		{
			headerName: t(translations.CommonTitles.Status()),
			field: 'status',
			...cellColorAndNameFormatter('status'),
			maxWidth: 130,
			minWidth: 100
		},
	]

	const filters: FilterArrayElement = {
		disabledCols: ['createdAt', 'updatedAt', 'id', 'price', 'amount', 'buyAmount', 'sellAmount', 'total'],
		//add00ToDate: true,
		dropDownCols: [
			{
				id: 'type',
				substituteId: 'type',
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
			{
				id: 'status',
				options: [
					{
						name: 'All',
						value: 'removeFilter',
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
						name: 'Open',
						value: 'open',
					},
					{
						name: 'Canceled',
						value: 'canceled',
					},
					{
						name: 'Expired',
						value: 'expired',
					},
					{
						name: 'Processing',
						value: 'processing',
					},
					{
						name: 'Admin Canceled',
						value: 'admin_canceled',
					},
					{
						name: 'Admin Postponed',
						value: 'admin_postponed',
					},
					{
						name: 'Admin Submitted',
						value: 'admin_submitted',
					},
				],
			},
			{
				id: 'pairCurrencyName',
				substituteId: 'pair_currency_id',
				options: JSON.parse(localStorage[LocalStorageKeys.PAIRS])
			}
		],
	};


	return (
		<>
			<PopupModal
				wrapperClassName='min1200'
				onClose={() => {
					setPopupInfo({ isOpen: false, data: [] });
				}}
				isOpen={popupInfo.isOpen}
			>
				<InfoPopup>
					{popupInfo.data.map((item, index) => {
						return <div key={index}><PopupSpan>{item.name}</PopupSpan>:<PopupSpanValue>{item.value}</PopupSpanValue></div>
					})}
				</InfoPopup>
			</PopupModal>

			<SimpleGrid
				containerId="externalOrders"
				additionalInitialParams={{ status: 'completed' }}
				arrayFieldName="externalExchangeOrders"
				immutableId="id"
				filters={filters}
				//  onRowClick={handleRowClick}
				additionalHeight={0}
				pullUpPaginationBy={32}
				initialAction={ExternalOrdersActions.GetExternalOrderAction}
				externalData={externalOrdersData}
				staticRows={staticRows}
			/>
		</>
	)
}

export default Executed
const InfoPopup = styled.div`
height:300px;
padding:12px;

`
const PopupSpan = styled.span`
font-size: 15px;
    font-weight: 700;
`
const PopupSpanValue = styled.span`
    font-size: 12px;
    font-weight: 600;
    margin-left: 7px;
`

const ExpandButton = styled.button`
cursor:pointer;

    width: 93px;
    height: 25px;
    margin-top: 7px;
    border-radius: 6px;
    border: none;
    background: #fe9f11 !important;
	display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 0px;
	font-size: 11px;
    font-weight: 600;



`