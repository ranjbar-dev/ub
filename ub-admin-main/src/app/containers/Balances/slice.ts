import {PayloadAction} from '@reduxjs/toolkit';
import {createSlice} from 'utils/@reduxjs/toolkit';

import {ContainerState, IWallet} from './types';

// The initial state of the Balances container
export const initialState: ContainerState = {
  balances: null,
  transferModalBalances: null,
  balancesHistory: null,
  isLoading: false,
  error: null,
};

const balancesSlice=createSlice({
	name: 'balances',
	initialState,
	reducers: {
		GetBalancesAction(state,action: PayloadAction<Record<string, unknown>>) { },
		GetBalancesForTransferModalAction(state,action: PayloadAction<Record<string, unknown>>) { },
		UpdateAllBalancesAction(state,action: PayloadAction<Record<string, unknown>>) { },
		InternalTransferAction(state,action: PayloadAction<Record<string, unknown>>) { },
		GetBalanceHistoryAction(state,action: PayloadAction<Record<string, unknown>>) { },
		setBalancesData(state, action: PayloadAction<Record<string, unknown>>) {
			state.balances = action.payload;
		},
		setTransferModalBalancesData(state, action: PayloadAction<{ balances: IWallet[]; type: string }>) {
			state.transferModalBalances = action.payload;
		},
		setBalancesHistoryData(state, action: PayloadAction<Record<string, unknown>>) {
			state.balancesHistory = action.payload;
		},
	},
});

export const {
	actions: BalancesActions,
	reducer: BalancesReducer,
	name: sliceKey,
}=balancesSlice;
