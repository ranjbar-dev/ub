/* --- STATE --- */
export interface OrdersState {
  openOrdersData: Record<string, unknown> | null;
  orderHistoryData: Record<string, unknown> | null;
  tradeHistoryData: Record<string, unknown> | null;
  orderHistory: Order[] | null;
  tradeHistory: Order[] | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = OrdersState;
export interface Order {
  amount: string;
  createdAt: string;
  executed: string;
  fee: string;
  id: number;
  isMaker: boolean;
  isTradedByBot: boolean;
  pair: string;
  price: string;
  type: string;
  updateAt: string;
  userEmail: string;
  userId: number;
}
