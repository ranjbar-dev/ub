import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface ContactUsPageState {
  readonly default: any;
  readonly counterValue: number;
  readonly inputValue: string;
}

/* --- ACTIONS --- */
type ContactUsPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = ContactUsPageState;
type ContainerActions = ContactUsPageActions;

export { RootState, ContainerState, ContainerActions };
