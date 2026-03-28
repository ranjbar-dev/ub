import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { BalancePageData } from 'containers/FundsPage/types';

/* --- STATE --- */
interface TradeHeaderState {
  readonly default: any;
  readonly balancePageData: BalancePageData;
}

/* --- ACTIONS --- */
type TradeHeaderActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = TradeHeaderState;
type ContainerActions = TradeHeaderActions;

export { RootState, ContainerState, ContainerActions };
