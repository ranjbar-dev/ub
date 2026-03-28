import {toast} from 'app/components/Customized/react-toastify';
import {GridLoading} from 'app/components/grid_loading/gridLoading';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UBInput from 'app/components/UBInput/UBInput';
import {Buttons} from 'app/constants';
import {translations} from 'locales/i18n';
import React,{useRef,useState} from 'react'
import {useTranslation} from 'react-i18next';
import styled from 'styled-components/macro'

import {BalancesActions} from '../slice';
import {IWallet,WalletTypes} from '../types'
import FromToComponent from './FromToComponent';
import {Content,RowWrapper,Title} from './styledComponents';

interface Props {
	dispatch: (action: { type: string; payload?: unknown }) => void,
	from: WalletTypes,
	data: IWallet,
	allData: IWallet[],
	onCancel: () => void,
	onSubmit: (data: Record<string, unknown>) => void
}

function TransferModal(props: Props) {
	const {from,data,onCancel,onSubmit,allData,dispatch}=props
	const [AllData,setAllData]=useState<IWallet[]>(allData);


	const {t}=useTranslation()

	const [PopupData,setPopupData]=useState<IWallet>(data);
	const dataToSend=useRef<{
		loaderId: string;
		code: string;
		from: WalletTypes;
		fee: string;
		to: WalletTypes;
		amount: string;
		to_custom_address?: string;
		fromAddress?: string;
	}>(
		{
			loaderId: 'AllBalancesTransferButton'+data.code,
			code: data.code,
			from,
			fee: '0',
			to: from,
			amount: '0',
			...(data.network&&{network: data.network})
		}
	);


	const handleAmountChange = (e: string) => {
		dataToSend.current.amount=e
	}



	const handleSubmit=() => {
		if(dataToSend.current.amount==='0'||!dataToSend.current.amount) {
			toast.warn('amount is Reqired')
			return
		}
		if(dataToSend.current.from===dataToSend.current.to) {
			if(!dataToSend.current.to_custom_address) {
				toast.warn('from and to Wallets cant be the same Wallet')
				return
			}
		}
		//console.log(dataToSend.current)
		onSubmit(dataToSend.current)
	}

	const onWalletChange=({fromTo,value}: {fromTo: 'from' | 'to'; value: string}) => {
		dataToSend.current[fromTo==='from'? 'from':'to']=value as WalletTypes
	}
	const onAddressChange=({fromTo,value}: {fromTo: 'from' | 'to'; value: string}) => {
		dataToSend.current[fromTo==='from'? 'fromAddress':'to_custom_address']=value
	}
	const onCoinChange=({fromTo,value}: {fromTo: 'from' | 'to'; value: string}) => {
		dataToSend.current.code=value
		for(let i=0;i<AllData.length;i++) {
			if(AllData[i].code===value) {
				setPopupData(AllData[i])
			}
		}
	}





	return (
		<>

			<Wrapper >
				<TitleWrapper>

					<MainTitle>{` Transfet From ${from} Wallet`}</MainTitle>
				</TitleWrapper>
				<AllWrapper>
					<FromToWrapper>
						<div className="from">
							from:
					</div>
						<div className="to">to:</div>
					</FromToWrapper>
					<EditableWrapper>
						<FromToComponent
							network={PopupData.network}
							dispatch={dispatch}
							defaultCode={PopupData.code}
							defaultWallet={from}
							address={PopupData.address}
							balance={PopupData.free}
							fromTo='from'
							onWalletChange={onWalletChange}

							onCoinChange={onCoinChange}
						/>
						<FromToComponent
							dispatch={dispatch}
							defaultCode={PopupData.code}
							defaultWallet={from}
							address={PopupData.address}
							balance={PopupData.free}
							fromTo='to'
							onWalletChange={onWalletChange}
							onAddressChange={onAddressChange}
						/>
						<AmountAndFeeWrapper>
							<RowWrapper>
								<Title>{t(translations.CommonTitles.Amount())}:</Title>
								<Content>
									<UBInput onChange={handleAmountChange} initialValue={'0'} />
								</Content>
							</RowWrapper>

							{/* <RowWrapper>
								<Title>{t(translations.CommonTitles.Fee())}:</Title>
								<Content>
									<input style={{ width: '60%' }} onChange={(e) => handleFeeChange(e.target.value)} />
									<IsLoadingWithTextAuto
										text={t(translations.CommonTitles.Estimate())}
										className={Buttons.BlackButton}
										loadingId={'EstimateFee'}
										onClick={() => {
											handleEstimateClick();
										}}
									/>
								</Content>
							</RowWrapper> */}
						</AmountAndFeeWrapper>

						<CancelAndSubmitWrapper>
							<IsLoadingWithTextAuto
								text={t(translations.CommonTitles.Cancel())}
								className={Buttons.BlackButton}
								loadingId={'cancel'}
								onClick={() => {
									onCancel();
								}}
							/>
							<IsLoadingWithTextAuto
								text={t(translations.CommonTitles.Submit())}
								className={Buttons.SkyBlueButton}
								loadingId={'AllBalancesTransferButton'+data.code}
								onClick={() => {
									handleSubmit();
								}}
							/>
						</CancelAndSubmitWrapper>
					</EditableWrapper>

				</AllWrapper>
			</Wrapper>
		</>
	)
}

const CancelAndSubmitWrapper=styled.div`
    display: flex;
    justify-content: flex-end;
    padding: 0 12px;
    gap: 20px;
	.loadingCircle {
    top: 8px !important;
}
`
const AmountAndFeeWrapper=styled.div`


`

const EditableWrapper=styled.div`
flex:1;
`

const AllWrapper=styled.div`
display: flex;
`

const TitleWrapper=styled.div`
width:100%;
height:50px;
margin-bottom:30px;
display: flex;
align-items: center;
justify-content: center;
`
const MainTitle=styled.div`
    height: 75%;
	border-radius:7px;
	font-weight:600;
    width: 430px;
    background: rgb(213, 217, 227);
    display: flex;
    align-items: center;
    justify-content: center;
`

const FromToWrapper=styled.div`

height: 325px;

width:20%;
border-right:1px solid #D5D5D5;
.from,.to{
	text-align: end;
    padding: 0 10px;
}
.from{
	margin-bottom: 140px;
}

`

const Wrapper=styled.div`
width:550px;
height:500px;
background:#EBEBEB;
`
export default TransferModal
