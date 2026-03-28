/**
 *
 * Balances
 *
 */

import { ValueFormatterParams, ICellRendererParams } from 'ag-grid-community';
import GridTabs from 'app/components/GridTabs/GridTabs';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import PopupModal from 'app/components/materialModal/modal';
import { CellRenderer } from 'app/components/renderer';
import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import { FullWidthWrapper } from 'app/components/wrappers/FullWidthWrapper';
import { Buttons } from 'app/constants';
import { translations } from 'locales/i18n';
import React, { memo, useState, useMemo, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import { MessageNames, Subscriber } from 'services/messageService';
import { CurrencyFormater, safeFinancialAdd } from 'utils/formatters';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import TransferHistory from './components/TransferHistory';
import TransferModal from './components/TransferModal';
import UpdateAllBalancesButton from './components/updateButton';
import { balancesSaga } from './saga';
import { BalancesReducer, sliceKey, BalancesActions } from './slice';
import { selectBalancesData } from './selectors';
import { IWallet, WalletTypes } from './types';


interface Props { }

export const Balances = memo((props: Props) => {
	useInjectReducer({ key: sliceKey, reducer: BalancesReducer });
	useInjectSaga({ key: sliceKey, saga: balancesSaga });
	const [ShowGrid, setShowGrid] = useState<boolean>(true);
	const [ShowTransferHistory, setShowTransferHistory] = useState<boolean>(false);
	const [ModalData, setModalData] = useState<{ data: IWallet; from: WalletTypes }>();
	const [IsModalOpen, setIsModalOpen] = useState<boolean>(false);
	const gridData = useRef<IWallet[] | null>(null);
	const dispatch = useDispatch();
	const balancesData = useSelector(selectBalancesData);
	const walletType = useRef(WalletTypes.Internal)
	const { t } = useTranslation();
	const handleTransferClick = ({ data, from }: { data: IWallet; from: WalletTypes }) => {
		setModalData({
			from,
			data
		})
		setIsModalOpen(true)
	}
	const commonBetweenHotAndCold = [
		{
			headerName: t(translations.CommonTitles.Name()),
			field: 'name',
			maxWidth: 120,
		},
		{
			headerName: t(translations.CommonTitles.Address()),
			field: 'address',
			minWidth: 435
		},
		{
			headerName: t(translations.CommonTitles.Currency()),
			field: 'code',
			valueFormatter: (params: ValueFormatterParams) => {
				return (params.data.code) + (params.data.network ? ` (${params.data.network})` : '');
			},
		},
		{
			headerName: t(translations.CommonTitles.TotalBalances()),
			field: 'free',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.free);
			},
		},
	]
	const InternalFields = [
		{
			headerName: t(translations.CommonTitles.Name()),
			field: 'name',
			maxWidth: 120,
		},
		{
			headerName: t(translations.CommonTitles.Currency()),
			field: 'code',
			maxWidth: 85,
		},
		{
			headerName: t(translations.CommonTitles.TotalBalances()),
			field: 'free',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(safeFinancialAdd(params.data.free, params.data.locked));
			},
		},

		{
			headerName: t(translations.CommonTitles.TotalAvailableUsersBalance()),
			field: 'free',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.free);
			},
		},
		{
			headerName: t(translations.CommonTitles.InOrderOrPendingAmount()),
			field: 'locked',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.locked);
			},

		},
	]


	const handleUpdateClick= (sendingData: Record<string, unknown>) => {

		dispatch(BalancesActions.UpdateAllBalancesAction(sendingData))
	}

	const ColdWalletStaticRows = [
		...commonBetweenHotAndCold,
		{
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<IsLoadingWithTextAuto
							text={t(translations.CommonTitles.Update())}

							className={Buttons.GreenButton}
							loadingId={'balanceUpdateButton' + data.id}
							onClick={() => {

								handleUpdateClick({
									type: WalletTypes.Cold,
									code: data.code,
									loaderId: 'balanceUpdateButton' + data.id,
									...(data.network && { network: data.network })
								});
							}}
						/>
					</>,
				),
		}
	]
	const HotWalletStaticRows = [
		...commonBetweenHotAndCold,
		{
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<IsLoadingWithTextAuto
							text={t(translations.CommonTitles.Transfer())}

							className={Buttons.SkyBlueButton}
							loadingId={'AllBalancesTransferButton' + data.id}
							onClick={() => {
								handleTransferClick({ data, from: WalletTypes.Hot });
							}}
						/>
						<IsLoadingWithTextAuto
							text={t(translations.CommonTitles.Update())}

							className={Buttons.GreenButton}
							loadingId={'balanceUpdateButton' + data.id}
							onClick={() => {
								handleUpdateClick({
									type: WalletTypes.Hot,
									code: data.code,
									loaderId: 'balanceUpdateButton' + data.id,
									...(data.network && { network: data.network })
								});
							}}
						/>
					</>,
				),
		}
	]

	const InternalStaticRows = [
		...InternalFields,
		{
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<IsLoadingWithTextAuto
							text={t(translations.CommonTitles.Update())}
							className={Buttons.GreenButton}
							loadingId={'balanceUpdateButton' + data.id}
							onClick={() => {
								handleUpdateClick({ type: WalletTypes.Internal, code: data.code, loaderId: 'balanceUpdateButton' + data.id });
							}}
						/>
					</>,
				),
		}
	]
	const ExternalStaticRows = [
		{
			headerName: t(translations.CommonTitles.Name()),
			field: 'name',
			maxWidth: 120,
		},
		{
			headerName: t(translations.CommonTitles.Currency()),
			field: 'code',
			maxWidth: 85,
		},
		{
			headerName: t(translations.CommonTitles.Address()),
			field: 'address',
			minWidth: 330,
		},
		{
			headerName: t(translations.CommonTitles.TotalBalances()),
			field: 'free',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(safeFinancialAdd(params.data.free, params.data.locked));
			},
		},

		{
			headerName: t(translations.CommonTitles.TotalAvailableUsersBalance()),
			field: 'free',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.free);
			},
		},
		{
			headerName: t(translations.CommonTitles.InOrderOrPendingAmount()),
			field: 'locked',
			valueFormatter: (params: ValueFormatterParams) => {
				return CurrencyFormater(params.data.locked);
			},

		},
		{
			cellRenderer: ({ data, rowIndex }: ICellRendererParams) =>
				CellRenderer(
					<>
						<IsLoadingWithTextAuto
							text={t(translations.CommonTitles.Transfer())}
							className={Buttons.SkyBlueButton}
							loadingId={'AllBalancesTransferButton' + data.id}
							onClick={() => {
								handleTransferClick({ data, from: WalletTypes.External });
							}}
						/>
					</>,
				),
		}

	]


	const [GridParameters, setGridParameters] = useState({
		callParameters: {
			type: WalletTypes.Internal,
		},
		staticRows: InternalStaticRows
	}
	);



	const filterTabs = [
		{
			name: t(translations.CommonTitles.ClientWallet()),
			callObject: {
				type: WalletTypes.Internal,
			},
		},
		{
			name: t(translations.CommonTitles.HotWallets()),
			callObject: {
				type: WalletTypes.Hot,
			},
		},
		{
			name: t(translations.CommonTitles.LiquidityWallet()),
			callObject: {
				type: WalletTypes.External,
			},
		},
		{
			name: t(translations.CommonTitles.ColdWallets()),
			callObject: {
				type: WalletTypes.Cold,
			},
		},
		{
			name: t(translations.CommonTitles.TransferHistory()),
			callObject: {
				type: 'transferHistory',
			},
		},
	];

	const handleTopTabChange = (e: Record<string, unknown>) => {
		const tabParams = e as { type: WalletTypes };
		walletType.current = tabParams.type
		setShowGrid((prev) => false)
		setShowTransferHistory(false);
		switch (tabParams.type as string) {
			case WalletTypes.Internal:
				setGridParameters({
					callParameters: tabParams,
					staticRows: InternalStaticRows
				})
				break;
			case WalletTypes.External:
				setGridParameters({
					callParameters: tabParams,
					// @ts-expect-error — staticRows type doesn't match expected type
					staticRows: ExternalStaticRows
				})
				break;
			case WalletTypes.Cold:
				setGridParameters({
					callParameters: tabParams,
					// @ts-expect-error — staticRows type doesn't match expected type
					staticRows: ColdWalletStaticRows
				})
				break;
			case WalletTypes.Hot:
				setGridParameters({
					callParameters: tabParams,
					// @ts-expect-error — staticRows type doesn't match expected type
					staticRows: HotWalletStaticRows
				})
				break;
			case 'transferHistory':
				setShowTransferHistory(true)
				break;

			default:
				break;
		}

		setTimeout(() => {
			setShowGrid((prev) => true)
		}, 0);

	}

	const handleModalSubmit = (data: Record<string, unknown>) => {
		dispatch(BalancesActions.InternalTransferAction(data))
	}
	const handleDataRecieved = (data: Record<string, unknown>) => {
		gridData.current = (data as { balances: IWallet[] }).balances

	}

	useEffect(() => {
		const Subscription = Subscriber.subscribe((message: { name: string }) => {
			if (message.name === MessageNames.CLOSE_POPUP) {
				setIsModalOpen(false)
			}
		})
		return () => {
			Subscription.unsubscribe()
		}
	}, []);



	return (
		<FullWidthWrapper>

			<PopupModal
				onClose={() => {
					setIsModalOpen(false);
				}}
				isOpen={IsModalOpen}
			>
			{ModalData && (
				<TransferModal
					dispatch={dispatch}
					allData={gridData.current ?? []}
					onCancel={() => { setIsModalOpen(false) }}
					onSubmit={handleModalSubmit}
					data={ModalData.data}
					from={ModalData.from}
				/>
			)}
			</PopupModal>

			<TitledContainer
				id="balances"
				title={t(translations.CommonTitles.Balances())}
			>
				<GridTabs onChange={handleTopTabChange} tabs={filterTabs} />
				{!ShowTransferHistory && <UpdateAllBalancesButton dispatch={dispatch} type={GridParameters.callParameters.type} />}
				{ShowGrid && (ShowTransferHistory === false ? <SimpleGrid
					//topTabs={filterTabs}
					containerId="balances"
					additionalInitialParams={GridParameters.callParameters}
					arrayFieldName="balances"
					immutableId="id"
					onDataReceived={handleDataRecieved}
					//filters={{}}
					//  onRowClick={handleRowClick}
					additionalHeight={0}
					pullUpPaginationBy={32}
					initialAction={BalancesActions.GetBalancesAction}
					messageName={MessageNames.SET_BALANCES_DATA}
					externalData={balancesData}
					staticRows={GridParameters.staticRows}
				/> : <TransferHistory />)}
			</TitledContainer>
		</FullWidthWrapper>
	);
});
