/**
*
* LiquidityOrders
*
*/

import { ValueFormatterParams, CellClassParams } from 'ag-grid-community';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { Buttons } from 'app/constants';
import { translations } from 'locales/i18n';
import { FilterArrayElement } from 'locales/types';
import React, { memo, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { LocalStorageKeys } from 'services/constants';
import styled from 'styled-components/macro';
import { CurrencyFormater } from 'utils/formatters';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
import { stateStyler } from 'utils/stylers';

import { liquidityOrdersSaga } from './saga';
import { selectLiquidityOrdersData } from './selectors';
import { LiquidityOrdersActions, LiquidityOrdersReducer, sliceKey } from './slice';

interface Props { }


export const LiquidityOrders = memo((props: Props) => {
	useInjectReducer({ key: sliceKey, reducer: LiquidityOrdersReducer });
	useInjectSaga({ key: sliceKey, saga: liquidityOrdersSaga });
	const dispatch = useDispatch();
	const liquidityOrdersData = useSelector(selectLiquidityOrdersData);
	const { t } = useTranslation();
	const pagesize = useRef<string>('21');
	const staticRows = useMemo(
		() => [
			{
				headerName: t(translations.CommonTitles.Type()),
				field: 'clientOrderId',
				maxWidth: 120,
				minWidth: 110,
			},
			{
				headerName: t(translations.CommonTitles.Date()),
				field: 'time',
				maxWidth: 140,
				minWidth: 140,
			},
			{
				headerName: t(translations.CommonTitles.Pair()),
				field: 'pairCurrency',
			},
			{
				headerName: t(translations.CommonTitles.Type()),
				field: 'orderType',
			},
			{
				headerName: t(translations.CommonTitles.Side()),
				field: 'side',
				cellStyle: (params: CellClassParams) => {
					return {
						color: stateStyler(params.data.side),
						textTransform: 'capitalize',
					};
				},
				maxWidth: 110,
				minWidth: 110,
			},

			// {
			{
				headerName: t(translations.CommonTitles.Price()),
				field: 'price',
				valueFormatter: (params: ValueFormatterParams) => {
					return CurrencyFormater(params.data.price);
				},
			},
			{
				headerName: t(translations.CommonTitles.Executed()),
				field: 'executedAmount',
				valueFormatter: (params: ValueFormatterParams) => {
					return CurrencyFormater(params.data.executedAmount);
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
				headerName: t(translations.CommonTitles.Total()),
				field: 'total',
				valueFormatter: (params: ValueFormatterParams) => {
					return CurrencyFormater(params.data.total + '');
				},
			},
			// {
			{
				headerName: t(translations.CommonTitles.Status()),
				field: 'status',
			},
		],
		[],
	);

	const filters: FilterArrayElement = {
		dropDownCols: [{
			id: 'clientOrderId',
			substituteId: 'client_type',
			options: [
				{
					name: 'Web',
					value: 'web'
				}, {
					name: 'Exchange',
					value: 'exchange'
				},
			]
		},
		{
			id: 'side',
			options: [{
				name: 'BUY',
				value: 'buy'
			},
			{
				name: 'SELL',
				value: 'sell'
			}
			]
		},
		{
			id: 'orderType',
			options: [{
				name: 'Stop Loss Limit',
				value: 'stop_loss_limit'
			},
			{
				name: 'Take Profit Limit',
				value: 'take_profit_limit'
			}
				,
			{
				name: 'Market',
				value: 'market'
			}
				,
			{
				name: 'Limit',
				value: 'limit'
			}
			]
		},
		{
			id: 'status',
			options: [{
				name: 'New',
				value: 'new'
			},
			{
				name: 'Partially Filled',
				value: 'partially_filled'
			},
			{
				name: 'Filled',
				value: 'filled'
			},
			{
				name: 'Canceled',
				value: 'canceled'
			},
			{
				name: 'Pending Cancel',
				value: 'pending_cancel'
			},
			{
				name: 'Rejected',
				value: 'rejected'
			},
			{
				name: 'Expired',
				value: 'expired'
			},

			]
		},
		{
			id: 'pairCurrency',
			substituteId: 'pair_currency_id',
			options: JSON.parse(
				localStorage[LocalStorageKeys.PAIRS])
		},
		],
		disabledCols: ['time', 'status', 'executed', 'price', 'total', 'fee'],


	}

	const handleUpdateCommissionClick = () => {
		dispatch(
			LiquidityOrdersActions.UpdateCommissionReportAction({
				page: 1,
				page_size: pagesize.current
			})
		)
	}
	const setPageSize = (pageSize: string) => {
		pagesize.current = pageSize
	}
	return (
		<FullWidthWrapper>
			<TitledContainer
				id={'liquidityOrders'}
				title={t(translations.CommonTitles.LiquidityOrders())}
			>
				<UpdateButtonWrapper>
					<IsLoadingWithTextAuto
						text={t(translations.CommonTitles.Update())}
						className={Buttons.SkyBlueButton}
						loadingId={'updateCommisionReportButton'}
						onClick={() =>
							handleUpdateCommissionClick()
						}
					/>
				</UpdateButtonWrapper>
				<SimpleGrid
					onPageSizeSet={setPageSize}
					containerId="liquidityOrders"
					additionalInitialParams={{}}
					arrayFieldName="orders"
					immutableId="id"
					filters={filters}
					//  onRowClick={handleRowClick}
					initialAction={LiquidityOrdersActions.GetLiquidityOrdersAction}
					externalData={liquidityOrdersData}
					staticRows={staticRows}
				/>
			</TitledContainer>
		</FullWidthWrapper>
	);

});

const UpdateButtonWrapper = styled.div`
    position: absolute;
    top: -55px;
    right: 0;
	.loadingCircle{
		top:8px !important;
	}

`