import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the ExternalOrders container
export const initialState: ContainerState = { externalOrdersData: null, netQueueData: null, allQueueData: null, newQueueDetailList: null, isLoading: false, error: null };

const externalOrdersSlice = createSlice({
	name: 'externalOrders',
	initialState,
	reducers: {
		GetExternalOrderAction(state, action: PayloadAction<Record<string, unknown>>) { },
		GetNetQueueAction(state, action: PayloadAction<Record<string, unknown>>) { },
		GetAllQueueAction(state, action: PayloadAction<Record<string, unknown>>) { },
		ChangeNetQueueStatus(state, action: PayloadAction<Record<string, unknown>>) { },
		SubmitNetQueueAction(state, action: PayloadAction<Record<string, unknown>>) { },
		CancelNetQueueAction(state, action: PayloadAction<Record<string, unknown>>) { },
		GetListNetQueueAction(state, action: PayloadAction<Record<string, unknown>>) { },
		setExternalOrdersData(state, action: PayloadAction<Record<string, unknown>>) {
			state.externalOrdersData = action.payload;
		},
		setNetQueueData(state, action: PayloadAction<Record<string, unknown>>) {
			state.netQueueData = action.payload;
		},
		setAllQueueData(state, action: PayloadAction<Record<string, unknown>>) {
			state.allQueueData = action.payload;
		},
		setNewQueueDetailList(state, action: PayloadAction<Record<string, unknown> | null>) {
			state.newQueueDetailList = action.payload;
		},
	},
});

export const {
	actions: ExternalOrdersActions,
	reducer: ExternalOrdersReducer,
	name: sliceKey,
} = externalOrdersSlice;
