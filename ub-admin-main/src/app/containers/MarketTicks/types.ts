/* --- STATE --- */
export interface MarketTicksState {
  marketTicksData: Record<string, unknown> | null;
  syncListData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = MarketTicksState;
