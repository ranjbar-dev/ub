import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface UpdatePasswordPageState {
  readonly default: any;
}

/* --- ACTIONS --- */
type UpdatePasswordPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = UpdatePasswordPageState;
type ContainerActions = UpdatePasswordPageActions;

interface UpdatePasswordModel {
  code: string;
  password: string;
  confirmed: string;
  email: string;
  '2fa_code'?: string;
}
export { RootState, ContainerState, ContainerActions, UpdatePasswordModel };
