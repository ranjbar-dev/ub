import { ValueFormatterParams } from 'ag-grid-community';
import {AgGridReact} from 'ag-grid-react';
import {rowHeight} from 'app/constants';
import {translations} from 'locales/i18n';
import React from 'react'
import {useTranslation} from 'react-i18next';
import styled from 'styled-components/macro';
import {CurrencyFormater} from 'utils/formatters';
import {cellColorAndNameFormatter} from 'utils/stylers';

interface Props {
	data: Record<string, unknown>[]

}

function NetQueueDetailsListGrid(props: Props) {
	const {data}=props


	const {t}=useTranslation();

	const staticRows=[

		{
			headerName: t(translations.CommonTitles.OrderID()),
			field: 'orderId',
		},
		{
			headerName: t(translations.CommonTitles.UserId()),
			field: 'userId',
		},
		{
			headerName: t(translations.CommonTitles.Date()),
			field: 'date',
			maxWidth: 150,
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


	return (
		<GridWrapper className='ag-theme-balham' >
			<AgGridReact
				headerHeight={32}
				rowHeight={rowHeight}
				columnDefs={staticRows}
				defaultColDef={{suppressMenu: true,sortable: false}}
				rowData={
					data
				}
				immutableData={true}
				getRowNodeId={(data: Record<string, unknown>) => {
					return String(data.pairCurrencyId);
				}}
			></AgGridReact>


		</GridWrapper>
	)
}
const GridWrapper=styled.div`
width:1200px;
height:800px;

`
export default NetQueueDetailsListGrid
