import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface ChangePasswordState {
  readonly default: any;
  readonly isChangingPassword: boolean;
}
export interface ChangePasswordModel {
  old_password: string;
  new_password: string;
  confirmed: string;
}
/* --- ACTIONS --- */
type ChangePasswordActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = ChangePasswordState;
type ContainerActions = ChangePasswordActions;

export { RootState, ContainerState, ContainerActions };
