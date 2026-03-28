/* --- STATE --- */
export interface UserDetailsState {}

export type ContainerState = UserDetailsState;
export interface UserBalances {
  availableSum: string;
  balances: Balance[];
  btcAvailableSum: string;
  btcInOrderSum: string;
  btcTotalSum: string;
  inOrderSum: string;
  minimumOfSmallBalances: string;
  totalSum: string;
}
export interface Balance {
  availableAmount: string;
  backgroundImage: string;
  btcAvailableEquivalentAmount: string;
  btcInOrderEquivalentAmount: string;
  btcTotalEquivalentAmount: string;
  code: string;
  equivalentAvailableAmount: string;
  equivalentInOrderAmount: string;
  equivalentTotalAmount: string;
  fee: string;
  image: string;
  inOrderAmount: string;
  minimumWithdraw: string;
  name: string;
  price: string;
  subUnit: number;
  totalAmount: string;
}
export interface Address {
  address: string;
  code: string;
  id: number;
  isFavorite: boolean;
  label: string;
  name: string;
}
export interface Permission {
  id: number;
  name: string;
  userHasIt: boolean;
}
