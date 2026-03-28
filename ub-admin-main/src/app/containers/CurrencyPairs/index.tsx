/**
 *
 * CurrencyPairs
 *
 */

import { ColDef } from 'ag-grid-community';
import ConstructiveModal from 'app/components/ConstructiveModal/ConstructiveModal';
import PopupModal from 'app/components/materialModal/modal';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { translations } from 'locales/i18n';
import React, { memo, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { MessageNames, GridNames } from 'services/messageService';
import styled from 'styled-components/macro';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { currencyPairsSaga } from './saga';
import { CurrencyPairsReducer, sliceKey, CurrencyPairsActions } from './slice';
import { selectCurrencyPairsData } from './selectors';
import { IConstructiveModalElement } from '../FinanceMethods';


interface Props { }

export const CurrencyPairs = memo((props: Props) => {
	useInjectReducer({ key: sliceKey, reducer: CurrencyPairsReducer });
	useInjectSaga({ key: sliceKey, saga: currencyPairsSaga });
	const [IsModalOpen, setIsModalOpen] = useState(false);
	const [ModalData, setModalData] = useState({});
	const { t } = useTranslation();
	const dispatch = useDispatch();
	const currencyPairsData = useSelector(selectCurrencyPairsData);
	const staticRows: ColDef[] = useMemo(
		() => [
			{
				headerName: t(translations.CommonTitles.ID()),
				field: 'id',
				maxWidth: 70
			},
			{
				headerName: t(translations.CommonTitles.PairName()),
				field: 'name',
			},
			{
				headerName: t(translations.CommonTitles.Status()),
				field: 'isActive',
			},
			{
				headerName: t(translations.CommonTitles.Trade()),
				field: 'tradeStatus',
				valueFormatter: ({ data }: { data: { tradeStatus: string } }) => {
					return data.tradeStatus.replace("_", ' ');
				},
				maxWidth: 110

			},
			{
				headerName: "Show digits",
				field: 'showDigits',
				maxWidth: 150
			},
			{
				headerName: t(translations.CommonTitles.Type()),
				field: 'botSpreadType',
			},
			{
				headerName: t(translations.CommonTitles.BotSpread()),
				field: 'botSpread',
			},
			{
				headerName: t(translations.CommonTitles.MakerFee()),
				field: 'makerFee',
			},
			{
				headerName: t(translations.CommonTitles.TakerFee()),
				field: 'takerFee',
			},
			{
				headerName: t(translations.CommonTitles.SpreadPercent()),
				field: 'ohlcSpread',
			},
			{
				headerName: t(translations.CommonTitles.MinOrder()),
				field: 'minimumOrderAmount',
			},
			{
				headerName: "Our Max Limit",
				field: 'maxOurExchangeLimit',
			}, {
				headerName: "Order Agg Time",
				field: 'botOrdersAggregationTime',
				valueFormatter: ({ data }: { data: Record<string, unknown> }) => {
					return `${data.botOrdersAggregationTime} minute`;
				}
			},
		],
		[],
	);
	const handleRowClick = ({ data }: { data: Record<string, unknown> }) => {
		setModalData(data);
		setIsModalOpen(true);
	};
	const handleSubmit = (data: Record<string, unknown>) => {
		dispatch(CurrencyPairsActions.UpdateCurrencyPairAction(data));
	};
	/*
    
   botOrdersAggregationTime: 2
  botSpread: "0.0002"
  botSpreadType: "percentage"
  id: 1
  isActive: true
  makerFee: 0.01
  maxOurExchangeLimit: "3.0"
  minimumOrderAmount: "10"
  name: "BTC-USDT"
  ohlcSpread: 0
  showDigits: 6
  takerFee: 0.01
  tradeStatus: "full_trade"
	*/
	const modalFields: IConstructiveModalElement[] = [
		{
			name: 'Pair Name',
			field: 'name',
			editable: false,
		},
		{
			name: 'Status',
			field: 'isActive',
			type: 'dropDown',
			editable: true,
			options: [
				{
					name: 'Enabled',
					value: true,
				},
				{
					name: 'Disabled',
					value: false,
				},
			],
		},
		{
			name: 'Trade',
			field: 'tradeStatus',
			editable: true,
			type: 'dropDown',
			options: [
				{
					name: 'Full Trade',
					value: 'full_trade',
				},
				{
					name: 'Sell Only',
					value: 'sell_only',
				},
				{
					name: 'Buy Only',
					value: 'buy_only',
				},
				{
					name: 'Close Only',
					value: 'close_only',
				},
			],
		},
		{
			name: 'Type',
			field: 'botSpreadType',
			type: 'dropDown',
			editable: true,
			options: [
				{
					name: 'Percentage',
					value: 'percentage',
				},
				{
					name: 'Const',
					value: 'const',
				},
			],
		},
		{
			name: 'Show Digits',
			field: 'showDigits',
			editable: true,
		},
		{
			name: 'Bot Spread',
			field: 'botSpread',
			editable: true,
		},
		{
			name: 'Maker Fee',
			field: 'makerFee',
			editable: true,
		},
		{
			name: 'Taker Fee',
			field: 'takerFee',
			editable: true,
		},
		{
			name: 'OHLC Spread',
			field: 'ohlcSpread',
			editable: true,
		},
		{
			name: 'Min Order',
			field: 'minimumOrderAmount',
			editable: true,
		}, {
			name: 'Max Order',
			field: 'maxOurExchangeLimit',
			editable: true,
		},
		{
			name: 'Order Agg Time(min)',
			field: 'botOrdersAggregationTime',
			editable: true,
		},
	];
	return (
		<FullWidthWrapper>
			<PopupModal
				onClose={() => {
					setIsModalOpen(false);
				}}
				isOpen={IsModalOpen}
			>
				<ConstructiveModal
					onCancel={() => {
						setIsModalOpen(false);
					}}
					onSubmit={handleSubmit}
					// @ts-expect-error — ModalData type (empty object initially) doesn't match ConstructiveModal prop type
					initialData={ModalData}
					modalFields={modalFields}
				/>
			</PopupModal>
			<TitledContainer
				id="currencyPairs"
				title={"Currency Pairs"}
			>
				<SimpleGrid
					containerId="currencyPairs"
					additionalInitialParams={{}}
					arrayFieldName="data"
					flashCellUpdate={true}
					gridName={GridNames.CURRENCY_PAIRS_PAGE}
					immutableId="id"
					// filters={{}}
					onRowClick={handleRowClick}
					initialAction={CurrencyPairsActions.GetCurrencyPairsAction}
					messageName={MessageNames.SET_CURRENCYPAIRS_DATA}
					externalData={currencyPairsData}
					staticRows={staticRows}
				/>
			</TitledContainer>
		</FullWidthWrapper>
	);
});

const Wrapper = styled.div``;
