import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { Currency } from 'containers/AddressManagementPage/types';

/* --- STATE --- */
interface OrdersPageState {
  readonly default: any;
  readonly openOrders: Order[];
  readonly orderHistory: Order[];
  readonly tradeHistory: Order[];
  readonly currencies: Currency[];

  readonly isLoadingOpenOrders: boolean;
  readonly isLoadingOrderHistory: boolean;
  readonly isLoadingTradeHistory: boolean;
}
interface Order {
  mainType: string;
  type: string;
  id: number;
  pair: string;
  side: string;
  price: string;
  subUnit: number;
  averagePrice: string;
  amount: string;
  executed: string;
  total: string;
  createdAt: string;
  updatedAt: string;
  triggerCondition: string;
  status: string;
  details?: any;
  createdAtToFilter: string;
  isDetailsOpen: boolean;
}
interface WithdrawAddress {
  id: number;
  address: string;
  label: string;
  isFavorite: boolean;
  code: string;
  name: string;
}
interface OrderDetail {
  createdAt: string;
  pair: string;
  type: string;
  subUnit: number;
  price: string;
  executed: string;
  fee: string;
  amount: string;
}
interface OrderHistorySearchModel {
  period?: string;
  start_date?: string;
  end_date?: string;
  pair_currency_name?: string;
  type?: string;
}
interface FilterModel {
  period: string;
  start_date: string;
  end_date: string;
  pair_currency_name: string;
  type: string;
  hideCancelledOrders: boolean;
  code: string;
  address: string;
  dwType: string;
  silent?: boolean;
}
interface StreamOrder {
  amount: string;
  id: string;
  price: string;
  status: 'open' | 'filled';
}
/* --- ACTIONS --- */
type OrdersPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = OrdersPageState;
type ContainerActions = OrdersPageActions;

export {
  RootState,
  ContainerState,
  ContainerActions,
  OrdersPageState,
  Order,
  OrderHistorySearchModel,
  OrderDetail,
  FilterModel,
  StreamOrder,
};
