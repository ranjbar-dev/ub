/* --- STATE --- */
export interface FinanceMethodsState {
  financeMethodsData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = FinanceMethodsState;
