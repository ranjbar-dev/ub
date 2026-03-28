import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto'
import { translations } from 'locales/i18n'
import React, { memo } from 'react'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'

import { BalancesActions } from '../slice'

interface Props { dispatch: (action: { type: string; payload?: unknown }) => void, type: string }

const UpdateAllBalancesButton = memo((props: Props) => {
    const { t } = useTranslation()

    const { dispatch, type } = props
    const handleUpdateClick = () => {
        dispatch(BalancesActions.UpdateAllBalancesAction({ type, loaderId: 'UpdateAllBalances' + type }))
    }
    return (
        <Wrapper>
            <IsLoadingWithTextAuto
                text={t(translations.CommonTitles.Update())}
                loadingId={'UpdateAllBalances' + type}
                onClick={() => handleUpdateClick()}
            />
        </Wrapper>
    )
}, ({ type: oldType }, { type: newType }) => oldType === newType)

const Wrapper = styled.div`
    position: absolute;
    right: 0;
    top: 0;
    .loadingCircle{
        top:8px !important;
    }
`

export default UpdateAllBalancesButton
