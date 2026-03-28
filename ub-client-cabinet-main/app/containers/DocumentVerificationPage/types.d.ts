import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';

/* --- STATE --- */
interface DocumentVerificationPageState {
  readonly default: any;
  readonly userProfileData: UserProfileData;
  readonly isLoadingData: boolean;
}
interface UserProfileData {
  region_and_city?: string;
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
  city?: string;
  country_id?: number;
  changed?: boolean;
}
interface ProfileImage {
  image?: string;
  type?: string;
  id?: number;
  idCardCode?: string;
  status?: string;
  rejectionReason?: string;
  imageId?: number;
  mainImageId?: number;
  isBack?: boolean;
}
export interface IUserProfileImage {
  createdAt: string;
  id: number;
  idCardCode: string;
  image: string;
  imageId: number;
  isBack: boolean;
  mainImageId: any;
  rejectionReason: string;
  status: string;
  subType: string;
  type: string;
}
export interface IUserProfileMetaData {
  types: { name: string; subTypes: { name: string; hasBack: boolean }[] }[];
}
/* --- ACTIONS --- */
type DocumentVerificationPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = DocumentVerificationPageState;
type ContainerActions = DocumentVerificationPageActions;

export {
  RootState,
  ContainerState,
  ContainerActions,
  UserProfileData,
  ProfileImage,
};
