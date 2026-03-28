import {PayloadAction} from '@reduxjs/toolkit';
import {createSlice} from 'utils/@reduxjs/toolkit';

import {ContainerState} from './types';

// The initial state of the MarketTicks container
export const initialState: ContainerState = { marketTicksData: null, syncListData: null, isLoading: false, error: null };

const marketTicksSlice=createSlice({
	name: 'marketTicks',
	initialState,
	reducers: {
		GetMarketTicksAction(state,action: PayloadAction<Record<string, unknown>>) { },
		GetCurrencyPairsAction(state,action: PayloadAction<Record<string, unknown>>) { },
		SyncTicksAction(state,action: PayloadAction<Record<string, unknown>>) { },
		GetSyncListAction(state,action: PayloadAction<Record<string, unknown>>) { },
		setMarketTicksData(state,action: PayloadAction<Record<string,unknown>>) {
			state.marketTicksData=action.payload;
		},
		setSyncListData(state,action: PayloadAction<Record<string,unknown>>) {
			state.syncListData=action.payload;
		},
	},
});

export const {
	actions: MarketTicksActions,
	reducer: MarketTicksReducer,
	name: sliceKey,
}=marketTicksSlice;
