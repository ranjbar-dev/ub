import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface GoogleAuthenticationPageState {
  readonly default: any;
  readonly isLoading: boolean;
  readonly qrCode: QrCode;
}

/* --- ACTIONS --- */
type GoogleAuthenticationPageActions = ActionType<typeof actions>;
interface QrCode {
  url?: string;
  code?: string;
}
export interface SetG2FaModel {
  code: string;
  password: string;
  setEnable: boolean;
}
/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = GoogleAuthenticationPageState;
type ContainerActions = GoogleAuthenticationPageActions;

export { RootState, ContainerState, ContainerActions, QrCode };
