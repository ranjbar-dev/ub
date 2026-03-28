import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from '../../types';
import { Themes } from './constants';
import { Country } from 'containers/PhoneVerificationPage/constants';

/* --- STATE --- */

interface AppState {
  readonly loading: boolean;
  readonly error?: object | boolean;
  readonly loggedIn?: object | boolean;
  readonly currencies?: Currency[];
  readonly countries?: Country[];

  readonly theme?: Themes;
}
interface Currency {
  id: number;
  code: string;
  name: string;
  image: string;
  otherBlockChainNetworks?: any[];
  mainNetwork?: string;
  showDigits: number;
  backgroundImage: string;
}
/* --- ACTIONS --- */
type AppActions = ActionType<typeof actions>;

/* --- EXPORTS --- */

type RootState = ApplicationRootState;
type ContainerState = AppState;
type ContainerActions = AppActions;

export { RootState, ContainerState, ContainerActions, Currency };
