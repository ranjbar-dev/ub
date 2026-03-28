import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface SignupPageState {
  readonly default: any;
}

/* --- ACTIONS --- */
type SignupPageActions = ActionType<typeof actions>;
interface RegisterModel {
  email: string;
  password: string;
  recaptcha: string;
}
/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = SignupPageState;
type ContainerActions = SignupPageActions;

export { RootState, ContainerState, ContainerActions, RegisterModel };
