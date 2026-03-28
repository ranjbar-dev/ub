/* --- STATE --- */
export interface OpenOrdersState {
  openOrdersData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = OpenOrdersState;
export enum Sides {
  Buy = 'buy',
  BuyUpper = 'BUY',
  Sell = 'sell',
  SellUpper = 'SELL',
}