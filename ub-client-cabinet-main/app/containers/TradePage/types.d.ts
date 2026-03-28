import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface TradePageState {
  readonly default: any;
  pairsMap: {
    [key: string]: PairItem
  } | {}
}

/* --- ACTIONS --- */
type TradePageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = TradePageState;
type ContainerActions = TradePageActions;
enum eTransactionType {
  buy = 'buy',
  sell = 'sell',
}

interface IRegisteredUserNotificationPayload {
  amount: string;
  id: number;
  price: string;
  status: string;
  type: eTransactionType;
}

interface PairItem {
  basisCode: string
  dependentCode: string
  dependentId: number
  image: string
  pairId: number
  pairName: string
  showDigits: number
}

export {
  RootState,
  ContainerState,
  ContainerActions,
  IRegisteredUserNotificationPayload,
  eTransactionType,
  PairItem,
};
