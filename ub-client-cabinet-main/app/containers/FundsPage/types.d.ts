import { ActionType } from 'typesafe-actions';
import * as actions from './actions';
import { ApplicationRootState } from 'types';
import { DWStatus } from './constants';

/* --- STATE --- */
interface FundsPageState {
  readonly default: any;
  readonly balancePageData: BalancePageData;
  readonly depositAndWithDrawData: depositAndWithDrawData;
  readonly transactionHistoryData: Transaction[];
  readonly formerWithdrawAddresses: WithdrawAddress[];
  readonly isLoadingBalancePageData: boolean;
  readonly isLoadingDepositeAndWithdraw: boolean;
  readonly isLoadingTransactionHistory: boolean;
  readonly userData?: UserData;
}
interface Balance {
  totalAmount: string | number;
  availableAmount: string;
  inOrderAmount: string;
  equivalentTotalAmount: string;
  equivalentAvailableAmount: string;
  equivalentInOrderAmount: string;
  code: string;
  name: string;
  price: string;
  fee: string;
  subUnit: number;
  minimumWithdraw: string;
  btcTotalEquivalentAmount: string;
  btcAvailableEquivalentAmount: string;
  btcInOrderEquivalentAmount: string;
  image: string;
  backgroundImage: string;
}
interface BalancePageData {
  balances?: Balance[];
  totalSum?: string;
  availableSum?: string;
  inOrderSum?: string;
  btcTotalSum?: string;
  btcAvailableSum?: string;
  btcInOrderSum?: string;
  minimumOfSmallBalances?: string;
}
interface Transaction {
  id: number;
  status: DWStatus;
  type: string;
  amount: string;
  code: string;
  createdAt: string;
  address?: string;
  txId?: string;
  addressExplorerUrl?: string;
  txIdExplorerUrl?: string;
  isDetailsOpen: boolean;
}
interface IOtherNetworks {
  address: string;
  code: string;
  completedNetworkName: string;
  supportsDeposit: boolean;
  supportsWithdraw: boolean;
}
interface depositAndWithDrawData {
  walletAddress?: string;
  withdrawTransactions?: DWData[];
  depositTransactions?: DWData[];
  networksConfigsAndAddresses?: {
    address: string;
    code: string;
    completedNetworkName: string;
    fee: string;
    supportsDeposit: boolean;
    supportsWithdraw: boolean;
  }[];
  balance?: Balance;
  mainNetwork?: string;
  currencyExtraInfo?: {
    name: string;
    issueDate: string;
    totalAmount: string;
    circulation: string;
    links: [
      {
        name: string;
        link: string;
      },
    ];
    description: string;
    image: string;
  };
  balanceChart?: any[];
  supportsWithdraw?: boolean;
  supportsDeposit?: boolean;
  completedNetworkName?: string;
  otherNetworksConfigsAndAddresses?: IOtherNetworks[];
  isDepositPermissionGranted?: boolean;
  isWithdrawPermissionGranted?: boolean;
  withdrawComments?: string[];
  depositComments?: string[];
}
interface DWData {
  id: number;
  status: DWStatus;
  type: string;
  amount: string;
  code: string;
  createdAt: string;
  isDetailsOpen: boolean;
  isLoadingDetails: boolean;
  address?: string;
  details: OrderDetail;
}
interface WithdrawAddress {
  id: number;
  address: string;
  label: string;
  isFavorite: boolean;
  code: string;
  name: string;
}
interface OrderDetail {
  address: string;
  txId: string;
  rejectionReason?: string;
  addressExplorerUrl?: string;
  txIdExplorerUrl?: string;
}
interface WithdrawModel {
  label: string;
  code: string;
  amount: string;
  network?: string;
  address: string;
  G2fa_code?: string;
  email_code?: string;
}
interface InfiniteDwModel {
  code?: string;
  type: string;
  page_size: number;
  page: number;
  silent?: boolean;
}
/* --- ACTIONS --- */
type FundsPageActions = ActionType<typeof actions>;

/* --- EXPORTS --- */
type RootState = ApplicationRootState;
type ContainerState = FundsPageState;
type ContainerActions = FundsPageActions;
export interface UserData {
  google2faEnabled?: boolean;
  isAccountVerified?: boolean;
}
export {
  RootState,
  ContainerState,
  ContainerActions,
  Balance,
  BalancePageData,
  depositAndWithDrawData,
  DWData,
  OrderDetail,
  Transaction,
  WithdrawModel,
  InfiniteDwModel,
  IOtherNetworks,
};
