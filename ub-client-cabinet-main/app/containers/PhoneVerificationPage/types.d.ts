import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { PhoneVerificationSteps } from './constants';

/* --- STATE --- */
interface PhoneVerificationPageState {
  readonly default: any;
  readonly isCountriesLoading: any;
  readonly countries: any[];
  readonly isLoading: boolean;
  readonly activeStep: PhoneVerificationSteps;
  readonly enteredPhoneNumber: string;
}
interface Step {
  title: any;
  description: any;
  isSelected?: boolean;
}

/* --- ACTIONS --- */
type PhoneVerificationPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = PhoneVerificationPageState;
type ContainerActions = PhoneVerificationPageActions;

export { RootState, ContainerState, ContainerActions, Step };
