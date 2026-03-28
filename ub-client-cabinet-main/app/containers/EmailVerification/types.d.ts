import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface EmailAuthenticationState {
  readonly default: any;
}

/* --- ACTIONS --- */
type EmailAuthenticationActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = EmailAuthenticationState;
type ContainerActions = EmailAuthenticationActions;

export { RootState, ContainerState, ContainerActions };
