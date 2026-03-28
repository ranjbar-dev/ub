/* --- STATE --- */
export interface WithdrawalsState {
  withdrawalsData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = WithdrawalsState;