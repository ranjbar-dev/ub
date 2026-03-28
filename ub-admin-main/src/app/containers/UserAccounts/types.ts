/* --- STATE --- */
export interface UserAccountsState {
  userAccountsData: { count: number; users: User[] };
  isLoading: boolean;
}

export type ContainerState = UserAccountsState;
export interface User {
  addressConfirmationStatus: string;
  country: string;
  email: string;
  firstName: string;
  id: number;
  identityConfirmationStatus: string;
  lastName: string;
  manager: string;
  phoneConfirmationStatus: string;
  profileStatus: string;
  referKey: string;
  referralId: string;
  registeredIP: string;
  registrationDate: string;
  status: string;
}
export interface InitialUserDetails {
  accountStatus: string;
  address: string;
  addressConfirmationStatus: string;
  birthDate: string;
  city: string;
  country: string;
  countryId: number;
  email: string;
  fullName: string;
  gender: string;
  groupId: string;
  groupName: string;
  id: number;
  identityConfirmationStatus: string;
  managerId: string;
  managerName: string;
  metaData: {
    userLevels: Array<{ id: string; name: string }>;
    userGroups: Array<{ id: string; name: string }>;
    userStatuses: Array<{ id: string; name: string }>;
    userProfileStatuses: Array<{ id: string; name: string }>;
    countries: Array<{ id: string; name: string; code?: string }>;
  };
  mobile: string;
  phoneConfirmationStatus: string;
  postalCode: string;
  profileStatus: string;
  referKey: string;
  referralId: string;
  registeredIp: string;
  registrationDate: string;
  status: string;
  systemId: number;
  totalBalance: string;
  totalCommissions: string;
  totalDeposit: string;
  totalOnTrade: string;
  totalWithdraw: string;
  trustLevel: number;
  userLevelId: number;
  userLevelName: string;
}
export enum ConfirmationStatus {
  NotConfirmed = 'not_confirmed',
  Confirmed = 'confirmed',
  Incomplete = 'incomplete',
  Rejected = 'rejected',
}
