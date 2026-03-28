import { MenuItem, Select } from '@material-ui/core'
import { ColDef, ICellRendererParams, ValueFormatterParams } from 'ag-grid-community'
import { toast } from 'app/components/Customized/react-toastify'
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto'
import PopupModal from 'app/components/materialModal/modal'
import { CellRenderer } from 'app/components/renderer'
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid'
import { translations } from 'locales/i18n'
import { FilterArrayElement } from 'locales/types'
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useDispatch, useSelector } from 'react-redux'
import { LocalStorageKeys } from 'services/constants'
import styled from 'styled-components'
import { CurrencyFormater } from 'utils/formatters'
import { cellColorAndNameFormatter } from 'utils/stylers'

import NetQueueActions from '../components/netQueueActions'
import NetQueueDetailsListGrid from '../components/netQueueDetailsListGrid'
import { selectNetQueueData, selectNewQueueDetailList } from '../selectors'
import { ExternalOrdersActions } from '../slice'

interface Props { }

const statusOptions = [
	{
		name: 'Running',
		value: 'run',
	},
	{
		name: 'Paused',
		value: 'pause',
	},
	{
		name: 'Stopped',
		value: 'stop',
	},

]

const netQueActions = ({ data, dispatch }: { data: Record<string, unknown>; dispatch: ReturnType<typeof useDispatch> }) =>
	CellRenderer(
		<>
			<NetQueueActions dispatch={dispatch} data={data} />
		</>,
	)

const SelectWrapper = styled.div`
.MuiInputBase-root {
    font-size: 10px !important;
    font-weight: 500 !important;}
.MuiOutlinedInput-notchedOutline{
	border:none !important;
}
`

const StatusDropdown = ({ data, dispatch }: { data: Record<string, unknown>; dispatch: ReturnType<typeof useDispatch> }) => {
	const { pairCurrencyAggregationStatus } = data
	const [status, setStatus] = useState(pairCurrencyAggregationStatus)
	const handleChange = (e: React.ChangeEvent<{ value: unknown }>) => {
		const selectedValue = e.target.value
		setStatus(selectedValue)
		dispatch(ExternalOrdersActions.ChangeNetQueueStatus({ pair_currency_id: data.pairCurrencyId, status: selectedValue }))
	}

	return <SelectWrapper>
		<Select
			MenuProps={{
				getContentAnchorEl: null,
				anchorOrigin: {
					vertical: 'bottom',
					horizontal: 'left',
				},
			}}
			className="select"
			variant="outlined"
			margin="dense"
			value={status}
			onChange={handleChange}
		>
			{statusOptions.map((item, index) => {
				return (
					<MenuItem key={item.value} value={item.value}>
						{item.name}
					</MenuItem>
				);
			})}
		</Select>
	</SelectWrapper>
	// return <span>{data.pairCurrencyAggregationStatus}</span>
}

function NetQueue(props: Props) {
	const { } = props
	const dispatch = useDispatch()
	const { t } = useTranslation();
	const [IsModalOpen, setIsModalOpen] = useState<boolean>(false);
	const [ModalGridData, setModalGridData] = useState<Record<string, unknown>[]>([])
	const netQueueData = useSelector(selectNetQueueData);
	const newQueueDetailList = useSelector(selectNewQueueDetailList);

	useEffect(() => {
		if (newQueueDetailList === null) return;
		const list: Record<string, unknown>[] = (newQueueDetailList as Record<string, unknown>).externalExchangeOrders as Record<string, unknown>[] || [];
		if (list.length === 0) {
			toast.warn('no details on this order');
			return;
		}
		setModalGridData(list);
		setIsModalOpen(true);
	}, [newQueueDetailList]);


	const staticRows: ColDef[] = [

		{
			headerName: t(translations.CommonTitles.Pair()),
			field: 'pairCurrencyName',
			maxWidth: 140,
		},
		{
			headerName: t(translations.CommonTitles.Type()),
			field: 'type',
			maxWidth: 120,
			...cellColorAndNameFormatter('type'),
		},

		{
			headerName: t(translations.CommonTitles.Amount()),
			field: 'amount',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.amount);
			},
		},
		{
			headerName: t(translations.CommonTitles.Price()),
			field: 'price',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.price);
			},
		},
		{
			headerName: t(translations.CommonTitles.BuyAmount()),
			field: 'buyAmount',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.buy.amount);
			},
		},
		{
			headerName: t(translations.CommonTitles.BuyPrice()),
			field: 'buyPrice',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.buy.price);
			},
		},
		{
			headerName: t(translations.CommonTitles.SellAmount()),
			field: 'sellAmount',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.sell.amount);
			},
		},
		{
			headerName: t(translations.CommonTitles.SellPrice()),
			field: 'sellPrice',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.sell.price);
			},
		},
		{
			headerName: t(translations.CommonTitles.Status()),
			field: 'pairCurrencyAggregationStatus',
			cellRenderer: (params: ICellRendererParams) => CellRenderer(
				<>
					<StatusDropdown data={params.data} dispatch={dispatch} />
				</>,
			)
		},
		{
			headerName: '',
			field: 'actions',
			minWidth: 280,
			cellRenderer: (params: ICellRendererParams) => {
				return netQueActions({ data: params.data, dispatch })
			}
		},
	]



	const filters: FilterArrayElement = {
		disabledCols: ['createdAt', 'updatedAt', 'id', 'price', 'amount', 'buyAmount', 'sellAmount', 'total', 'buyPrice', 'sellPrice'],
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
				id: 'pairCurrencyAggregationStatus',
				options: statusOptions,
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
					setIsModalOpen(false);
					dispatch(ExternalOrdersActions.setNewQueueDetailList(null));
				}}
				isOpen={IsModalOpen}
			>
				<NetQueueDetailsListGrid data={ModalGridData} />
			</PopupModal>

			<SimpleGrid
				containerId="externalOrders"
				additionalInitialParams={{}}
				arrayFieldName="externalExchangeOrders"
				immutableId="pairCurrencyId"
				filters={filters}
				//  onRowClick={handleRowClick}
				pagination={false}
				additionalHeight={0}
				pullUpPaginationBy={32}
				initialAction={ExternalOrdersActions.GetNetQueueAction}
				externalData={netQueueData}
				staticRows={staticRows}
			/></>
	)
}

export default NetQueue
