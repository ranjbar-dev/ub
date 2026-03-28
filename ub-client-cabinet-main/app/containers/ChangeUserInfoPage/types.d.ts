import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface ChangeUserInfoPageState {
  readonly default: any;
  readonly userProfileData: UserProfileData;
  readonly isLoadingData: boolean;
}

/* --- ACTIONS --- */
type ChangeUserInfoPageActions = ActionType<typeof actions>;
interface UserProfileData {
  id?: number;
  updatedAt?: string;
  firstName?: string;
  lastName?: string;
  gender?: string;
  dateOfBirth?: string;
  address?: string;
  regionAndCity?: string;
  postalCode?: string;
  country?: number;
  countryName?: string;
  status?: string;
  adminComment?: string;
  userProfileImages?: ProfileImage[];
}
interface ProfileImage {
  image?: string;
  type?: string;
  idCardCode?: string;
  status?: string;
  rejectionReason?: string;
}
/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = ChangeUserInfoPageState;
type ContainerActions = ChangeUserInfoPageActions;

export { RootState, ContainerState, ContainerActions, UserProfileData };
