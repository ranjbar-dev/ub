/* --- STATE --- */
export interface CurrencyPairsState {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  currencyPairsData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = CurrencyPairsState;
