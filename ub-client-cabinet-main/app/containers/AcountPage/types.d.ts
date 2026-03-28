import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { SecurityLevel, KycStatus } from './constants';

/* --- STATE --- */
interface AcountPageState {
  readonly default: any;
  readonly isLoading: boolean;
  readonly userData: UserData | null;
}
export interface UserData {
  email: string;
  ubId: string;
  phone: string;
  kycLevel: string;
  kycStatus: KycStatus;
  kycLevelMessage: string;
  securityLevel: SecurityLevel;
  securityLevelMessage: string;
  profileStatus: string;
  google2faEnabled: boolean;
  has2fa: boolean;
  isAccountVerified: boolean;
  channelName: string;
  themeId: number;
}
/* --- ACTIONS --- */
type AcountPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = AcountPageState;
type ContainerActions = AcountPageActions;

export { RootState, ContainerState, ContainerActions };
