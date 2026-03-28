import { GridLoading } from 'app/components/grid_loading/gridLoading'
import UbDropDown from 'app/components/UbDropDown'
import UBInput from 'app/components/UBInput/UBInput'
import { translations } from 'locales/i18n'
import React, { useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useSelector } from 'react-redux'
import { LocalStorageKeys } from 'services/constants'
import styled from 'styled-components'
import { CurrencyFormater } from 'utils/formatters'

import { BalancesActions } from '../slice'
import { selectTransferModalBalancesData } from '../selectors'
import { WalletTypes } from '../types'
import { Content, RowWrapper, Title } from './styledComponents'

interface Props {
	fromTo: 'from' | 'to',
	defaultWallet: WalletTypes,
	defaultCode: string,
	address?: string,
	dispatch: (action: { type: string; payload?: unknown }) => void;
	network?: string
	balance?: string,
	onWalletChange: ({ fromTo, value }: { fromTo: 'from' | 'to', value: string }) => void,
	onCoinChange?: ({ fromTo, value }: { fromTo: 'from' | 'to', value: string }) => void,
	onAddressChange?: ({ fromTo, value }: { fromTo: 'from' | 'to', value: string }) => void
}
const fromWallets = [
	{
		name: 'Hot Wallet',
		value: 'hot',
	},
	{
		name: 'Liquidity Wallet',
		value: 'external',
	},
]
const toWallets = [
	{
		name: 'Hot Wallet',
		value: 'hot',
	},
	{
		name: 'Liquidity Wallet',
		value: 'external',
	},
	{
		name: 'Cold Wallet',
		value: 'cold',
	},
]


function FromToComponent(props: Props) {

	const currencies = localStorage[LocalStorageKeys.CURRENCIES] ? JSON.parse(localStorage[LocalStorageKeys.CURRENCIES]).currencies.map((item: { name: string; code: string }, index: number) => {
		return {
			name: item.name,
			value: item.code
		}
	}) : []



	const { fromTo, defaultWallet, defaultCode, address, balance, onWalletChange, onAddressChange, onCoinChange, dispatch, network } = props

	const selectedWallet = useRef(defaultWallet)
	const selectedCoin = () => {
		let coinName = ''
		for (const coin of currencies) {
			if (coin.value === defaultCode) {
				coinName = coin.name + (network ? ` (${network})` : '');
				break
			}
		}
		return coinName
	}
	const [Address, setAddress] = useState<string>(address ?? '');
	const [Balance, setBalance] = useState<string>(balance ?? '');
	const [IsLoading, setIsLoading] = useState<boolean>(false);

	const transferModalData = useSelector(selectTransferModalBalancesData);
	useEffect(() => {
		if (transferModalData && transferModalData.type === selectedWallet.current) {
			setIsLoading(false);
			for (const b of transferModalData.balances) {
				if (b.code === defaultCode) {
					setAddress(b.address);
					setBalance(b.free);
					break;
				}
			}
		}
	}, [transferModalData, defaultCode]);



	const { t } = useTranslation()
	const handleWalletSelect = (value: WalletTypes) => {
		setIsLoading(true)
		selectedWallet.current = value
		dispatch(BalancesActions.GetBalancesForTransferModalAction({ type: value }))
		onWalletChange({ fromTo, value })
	}
	const handleCoinSelect = (value: string) => {
		onCoinChange && onCoinChange({ fromTo, value })
		//dispatch(BalancesActions.GetBalancesForTransferModalAction({type: selectedWallet.current}))
	}
	const handleAddressChange = (value: string) => {
		onAddressChange && onAddressChange({ fromTo, value })
		onAddressChange && setAddress(value)
	}
	return (
		<Wrapper>
			{IsLoading && <GridLoading />}
			<RowWrapper>
				<Title>{t(translations.CommonTitles.Wallet())}:</Title>
				<Content>
					<UbDropDown
						style={{ marginBottom: '0', maxWidth: '50%' }}
						options={fromTo === 'from' ? fromWallets : toWallets}
						initialValue={defaultWallet}
						onSelect={(e) => handleWalletSelect(e as WalletTypes)}
					/>
					{fromTo === 'from' && <CoinNameWrapper>{selectedCoin()}</CoinNameWrapper>}
					{/*{onCoinChange&&<UbDropDown
						style={{marginBottom: '0'}}
						options={currencies}
						initialValue={defaultCode}
						onSelect={(e) => handleCoinSelect(e)}
					/>}*/}
				</Content>
			</RowWrapper>
			{Address && <RowWrapper>
				<Title>{t(translations.CommonTitles.Address())}:</Title>
				<Content>
					<AddressWrapper  ><input onChange={(e) => handleAddressChange(e.target.value)} value={Address} /></AddressWrapper>
				</Content>
			</RowWrapper>}
			{Balance !== '' && <RowWrapper>
				<Title>{t(translations.CommonTitles.Balance())}:</Title>
				<Content>
					<BalanceWrapper>
						{CurrencyFormater(Balance + '')}
					</BalanceWrapper>
				</Content>
			</RowWrapper>}
		</Wrapper>
	)
}
const BalanceWrapper = styled.div`
background:#DEDEDE;
border-radius:8px;
padding:3px 5px;
font-size: 14px;
font-weight: 600;
`
const CoinNameWrapper = styled.div`

padding:3px 5px;
font-size: 14px;
font-weight: 600;
`


const AddressWrapper = styled.div`
background:#ebebeb;
border-radius:8px;
input{
    width: 300px;
    font-size: 12px;
    height: 32px;
    border: 1px solid #cecece;
    border-radius: 6px;
	font-weight: 600;
}

`




const Wrapper = styled.div`
width:100%;
border-bottom:1px solid #D5D5D5;

`

export default FromToComponent
