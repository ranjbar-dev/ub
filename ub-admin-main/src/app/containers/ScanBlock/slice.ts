import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the ScanBlock container
export const initialState: ContainerState = { isLoading: false, error: null };

const scanBlockSlice = createSlice({
  name: 'scanBlock',
  initialState,
  reducers: {
    Scan(state, action: PayloadAction<{network:string,block_number:number}>) {},
  },
});

export const { actions:ScanBlockActions, reducer:ScanBlockReducer, name: sliceKey } = scanBlockSlice;
