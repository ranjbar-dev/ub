import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto'
import {Buttons} from 'app/constants'
import {translations} from 'locales/i18n'
import React from 'react'
import {useTranslation} from 'react-i18next'
import styled from 'styled-components/macro'

import {ExternalOrdersActions} from '../slice'
import ListIcon from './listIcon'

interface Props {
	data: Record<string, unknown>
	dispatch: (action: { type: string; payload?: unknown }) => void
}

function NetQueueActions(props: Props) {
	const {data,dispatch}=props
	const {t}=useTranslation()
	const handleSubmit=() => {
		dispatch(ExternalOrdersActions.SubmitNetQueueAction({
			pair_currency_id: data.pairCurrencyId,
			action: 'submit'
		}))
	}
	const handleReset=() => {
		dispatch(ExternalOrdersActions.CancelNetQueueAction({
			pair_currency_id: data.pairCurrencyId,
			action: 'cancel'
		}))
	}
	const handleShowList=() => {
		dispatch(ExternalOrdersActions.GetListNetQueueAction({
			pair_currency_id: data.pairCurrencyId,
		}))
	}

	return (
		<Wrapper >
			<IsLoadingWithTextAuto
				text={t(translations.CommonTitles.Submit())}
				className={Buttons.SkyBlueButton}
				loadingId={'SubmitNetQueueRowButton'+data.pairCurrencyId}
				onClick={() => {
					handleSubmit();
				}}
			/>
			<IsLoadingWithTextAuto
				text={t(translations.CommonTitles.Reset())}
				className={Buttons.BlackButton}
				loadingId={'ResetNetQueueRowButton'+data.pairCurrencyId}
				onClick={() => {
					handleReset();
				}}
			/>
			<IsLoadingWithTextAuto
				text={<ListIcon />}
				className={Buttons.LightGreenButton}
				loadingId={'NetQueueRowShowListButton'+data.pairCurrencyId}
				onClick={() => {
					handleShowList();
				}}
			/>
		</Wrapper>
	)
}
const Wrapper=styled.div``
export default NetQueueActions
