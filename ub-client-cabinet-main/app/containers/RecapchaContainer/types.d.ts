import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface RecapchaContainerState {
  readonly default: any;
  readonly recapcha: string;
}

/* --- ACTIONS --- */
type RecapchaContainerActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = RecapchaContainerState;
type ContainerActions = RecapchaContainerActions;

export { RootState, ContainerState, ContainerActions };
