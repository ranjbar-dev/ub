import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { Currency } from 'containers/App/types';

/* --- STATE --- */
interface AddressManagementPageState {
  readonly default: any;
  readonly isLoading: boolean;
  readonly currencies: Currency[];
  readonly withDrawArrdesses: WithdrawAddress[];
  // readonly isAddingAddress: boolean;
}

interface WithdrawAddress {
  id: number;
  address: string;
  label: string;
  isFavorite: boolean;
  code: string;
  name: string;
}
/* --- ACTIONS --- */
type AddressManagementPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = AddressManagementPageState;
type ContainerActions = AddressManagementPageActions;

export {
  RootState,
  ContainerState,
  ContainerActions,
  Currency,
  WithdrawAddress,
  // GridFilterTypes,
};
