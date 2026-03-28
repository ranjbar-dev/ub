import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the LiquidityOrders container
export const initialState: ContainerState = { liquidityOrders: null, isLoading: false, error: null };

const liquidityOrdersSlice = createSlice({
  name: 'liquidityOrders',
  initialState,
  reducers: {
		GetLiquidityOrdersAction(state, action: PayloadAction<Record<string, unknown>>) {
			state.isLoading = true;
		},
		setLiquidityOrdersData(state, action: PayloadAction<Record<string, unknown>>) {
			state.liquidityOrders = action.payload;
			state.isLoading = false;
		},
		UpdateCommissionReportAction(state, action: PayloadAction<Record<string, unknown>>) {},
	},
});

export const { actions:LiquidityOrdersActions, reducer:LiquidityOrdersReducer, name: sliceKey } = liquidityOrdersSlice;
