/**
 *
 * MarketTicks
 *
 */

import GridTabs from 'app/components/GridTabs/GridTabs';
import {SimpleGrid} from 'app/components/SimpleGrid/SimpleGrid';
import TitledContainer from 'app/components/titledContainer/TitledContainer';
import {FullWidthWrapper} from 'app/components/wrappers/FullWidthWrapper';
import {translations} from 'locales/i18n';
import {FilterArrayElement} from 'locales/types';
import React,{memo,useMemo,useState} from 'react';
import {useTranslation} from 'react-i18next';
import {useSelector} from 'react-redux';
import styled from 'styled-components/macro';
import {useInjectReducer,useInjectSaga} from 'utils/redux-injectors';

import {MarketTicksTimeFrames} from './constants';
import {marketTicksSaga} from './saga';
import {selectMarketTicksData,selectSyncListData} from './selectors';
import {MarketTicksReducer,sliceKey,MarketTicksActions} from './slice';
import SyncPage from './SyncPage';



interface Props { }

export const MarketTicks=memo((props: Props) => {
	useInjectReducer({key: sliceKey,reducer: MarketTicksReducer});
	useInjectSaga({key: sliceKey,saga: marketTicksSaga});

	const {t}=useTranslation();

	const marketTicksData=useSelector(selectMarketTicksData);
	const syncListData=useSelector(selectSyncListData);

	const [ActivePage,setActivePage]=useState('ticks');
	const staticRows=useMemo(
		() => [
			{
				headerName: t(translations.CommonTitles.ID()),
				field: 'id',
			},
			{
				headerName: t(translations.CommonTitles.PairName()),
				field: 'pairCurrencyName',
			},
			{
				headerName: t(translations.CommonTitles.TimeFrame()),
				field: 'timeFrame',
			},
			{
				headerName: t(translations.CommonTitles.StartTime()),
				field: 'startTime',
			},
			{
				headerName: t(translations.CommonTitles.EndTime()),
				field: 'endTime',
			},
			{
				headerName: t(translations.CommonTitles.Open()),
				field: 'openPrice',
			},
			{
				headerName: t(translations.CommonTitles.Close()),
				field: 'closePrice',
			},
			{
				headerName: t(translations.CommonTitles.Low()),
				field: 'lowPrice',
			},
			{
				headerName: t(translations.CommonTitles.High()),
				field: 'highPrice',
			},
		],
		[],
	);

	/*
	
	endTime: "2021-01-18 00:00:00"
	id: 7484
	pairCurrencyName: "TRX-ETH"
	startTime: "2021-01-17 00:00:00"
	status: "done"
	timeFrame: "1day"
	
	*/
	const syncListStaticRows=[
		{
			headerName: t(translations.CommonTitles.ID()),
			field: 'id',
		},
		{
			headerName: t(translations.CommonTitles.PairName()),
			field: 'pairCurrencyName',
		},
		{
			headerName: t(translations.CommonTitles.StartTime()),
			field: 'startTime',
		},
		{
			headerName: t(translations.CommonTitles.EndTime()),
			field: 'endTime',
		},
		{
			headerName: t(translations.CommonTitles.TimeFrame()),
			field: 'timeFrame',
		},
		{
			headerName: t(translations.CommonTitles.Status()),
			field: 'status',

		},


	]


	const tabs=[
		{
			name: t(translations.CommonTitles.AllTicks()),
			callObject: {
				page: 'ticks',
			},
		},
		{
			name: t(translations.CommonTitles.Synchronize()),
			callObject: {
				page: 'sync',
			},
		},
		{
			name: t(translations.CommonTitles.SyncronizedList()),
			callObject: {
				page: 'list',
			},
		},
	];
	const handleTopTabChange = (e: Record<string, unknown>) => {
		//console.log(e);
		setActivePage(e.page as string);
	};

	const filters: FilterArrayElement=
	{
		dropDownCols: [
			{
				id: 'timeFrame',options: MarketTicksTimeFrames
			}
			,{
				id: 'status',options: [

					{
						name: 'Created',
						value: 'created'
					},

					{
						name: 'Processing',
						value: 'processing'
					},

					{
						name: 'Done',
						value: 'done'
					},

					{
						name: 'Error',
						value: 'error'
					},
				]
			}
		],
		dateCols: ['startTime','endTime'],
		//add00ToDate: true

	}



	return (
		<FullWidthWrapper>
			<TitledContainer
				id="marketTicks"
				title={t(translations.CommonTitles.MarketTicks())}
			>
				<GridTabs onChange={handleTopTabChange} tabs={tabs} />
				{ActivePage==='ticks'&&(
					<SimpleGrid
						containerId="marketTicks"
						additionalInitialParams={{}}
						arrayFieldName="data"
						immutableId="id"
						additionalHeight={0}
						pullUpPaginationBy={32}
						filters={filters}
						//  onRowClick={handleRowClick}
						initialAction={MarketTicksActions.GetMarketTicksAction}
						externalData={marketTicksData}
						staticRows={staticRows}
					/>
				)}
				{ActivePage==='sync'&&<SyncPage />}
				{ActivePage==='list'&&<SimpleGrid
					containerId="marketTicks"
					additionalInitialParams={{}}
					arrayFieldName="data"
					immutableId="id"
					additionalHeight={0}
					pullUpPaginationBy={32}
					filters={filters}
					//  onRowClick={handleRowClick}
					initialAction={MarketTicksActions.GetSyncListAction}
					externalData={syncListData}
					staticRows={syncListStaticRows}
				/>}
			</TitledContainer>
		</FullWidthWrapper>
	);
});

const Wrapper=styled.div``;
