/**
 *
 * Deposits
 *
 */

import { ValueFormatterParams, RowClickedEvent } from 'ag-grid-community';
import React,{useMemo,useState,useRef,useCallback} from 'react';
import {MessageNames,GridNames} from 'services/messageService';
import {useTranslation} from 'react-i18next';
import styled from 'styled-components/macro';

import {useInjectReducer,useInjectSaga} from 'utils/redux-injectors';
import {DepositsReducer,sliceKey,DepositsActions} from './slice';
import {depositsSaga} from './saga';
import {FullWidthWrapper} from 'app/components/wrappers/FullWidthWrapper';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import {translations} from 'locales/i18n';
import {SimpleGrid} from 'app/components/SimpleGrid/SimpleGrid';
import PopupModal from 'app/components/materialModal/modal';
import DepositModal from '../Billing/components/DepositModal';
import {DepositSaveData} from '../Billing/types';
import {useDispatch,useSelector} from 'react-redux';
import {stateStyler,cellColorAndNameFormatter} from 'utils/stylers';
import {CurrencyFormater} from 'utils/formatters';
import {FilterArrayElement} from 'locales/types';
import {LocalStorageKeys} from 'services/constants';
import {selectDepositsData} from './selectors';

interface Props { }

export function Deposits(props: Props) {
useInjectReducer({key: sliceKey,reducer: DepositsReducer});
useInjectSaga({key: sliceKey,saga: depositsSaga});
const dispatch=useDispatch();
const depositsData = useSelector(selectDepositsData);
const [IsModalOpen,setIsModalOpen]=useState(false);
const ModalData = useRef<RowClickedEvent | null>(null);
const {t}=useTranslation();
const staticRows=useMemo(
() => [
{
headerName: t(translations.CommonTitles.ID()),
field: 'id',
maxWidth: 130,
},
{
headerName: t(translations.CommonTitles.Method()),
field: 'currencyCode',
maxWidth: 130,
valueFormatter: (params: ValueFormatterParams) => {
return (params.data.currencyCode)+(params.data.blockchainNetwork? ` (${params.data.blockchainNetwork})`:'');
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
headerName: t(translations.CommonTitles.FromAddress()),
field: 'fromAddress',
},
{
headerName: t(translations.CommonTitles.ToAddress()),
field: 'toAddress',
},
{
headerName: t(translations.CommonTitles.TransactionId()),
field: 'txId',
},
{
headerName: t(translations.CommonTitles.CreationDate()),
field: 'createdAt',
maxWidth: 140,
},
{
headerName: t(translations.CommonTitles.LastUpdate()),
field: 'updatedAt',
maxWidth: 140,
},
{
headerName: t(translations.CommonTitles.Status()),
field: 'status',
maxWidth: 130,
...cellColorAndNameFormatter('status'),
},
],
[],
);
const handleRowClick=useCallback(e => {
ModalData.current=e;
setIsModalOpen(true);
},[]);
const handleUpdate=(data: DepositSaveData) => {
dispatch(DepositsActions.UpdateDepositsAction(data));
};
const filters: FilterArrayElement={
dateCols: ['createdAt','updatedAt'],
dropDownCols: [
{
id: 'status',
options: [
{
name: 'Created',
value: 'created',
},
{
name: 'In Progress',
value: 'in_progress',
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
name: 'Canceled',
value: 'cancel',
},
{
name: 'Rejected',
value: 'reject',
},
],
},
{
id: 'currencyCode',
options: JSON.parse(
localStorage[LocalStorageKeys.CURRENCIES],
).currencies.map((item: { name: string; code: string },index: number) => {
return {name: item.name,value: item.code};
}),
},
],
};
return (
<FullWidthWrapper>
<PopupModal
onClose={() => {
setIsModalOpen(false);
}}
isOpen={IsModalOpen}
>
<DepositModal
onSave={handleUpdate}
onCancel={() => {
setIsModalOpen(false);
}}
row={ModalData.current!}
/>
</PopupModal>

<TitledContainer
id="deposits"
title={t(translations.CommonTitles.Deposits())}
>
<SimpleGrid
containerId="deposits"
gridName={GridNames.DEPOSITS_PAGE}
additionalInitialParams={{type: 'deposit'}}
arrayFieldName="payments"
immutableId="id"
filters={filters}
onRowClick={handleRowClick}
initialAction={DepositsActions.GetDepositsAction}
messageName={MessageNames.SET_DEPOSITS_DATA}
externalData={depositsData}
staticRows={staticRows}
/>
</TitledContainer>
</FullWidthWrapper>
);
}

const Wrapper=styled.div``;