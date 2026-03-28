/* --- STATE --- */
export interface BillingState {
  billingData: Record<string, unknown> | null;
  depositsData: Record<string, unknown> | null;
  withdrawalsData: Record<string, unknown> | null;
  allTransactionsData: Record<string, unknown> | null;
  selectedPaymentDetails: { rowData: Payment; details: PaymentDetails } | null;
  commissions: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = BillingState;
export interface Payment {
  [key: string]: unknown;
  amount: string;
  createdAt: string;
  currencyCode: string;
  currencyId: number;
  fromAddress: string;
  id: number;
  status: 'completed' | 'canceled' | 'rejected' | 'in_progress' | 'user_canceled' | 'created' | 'pending';
  toAddress: string;
  txId: string;
  type: string;
  updatedAt: string;
  userEmail: string;
  userId: number;
  should_deposit?: boolean;
}
export interface IAdminComment {
  adminName: string;
  comment: string;
  id: number;
  updatedAt: string;
}
export interface PaymentDetails {
  adminComments: IAdminComment[];
  adminStatus: string;
  amount: string;
  autoTransfer: boolean;
  country: string;
  createdAt: string;
  currencyCode: string;
  currencyName: string;
  fee: string;
  fromAddress: string;
  id: number;
  ip: string;
  level: string;
  name: string;
  status: string;
  tag: string;
  toAddress: string;
  totalAmount: string;
  txId: string;
  type: string;
  userEmail: string;
  userId: number;
  rejectionReason?: string;
}
export interface DepositSaveData {
  amount: string;
  from_address: string;
  id: number;
  should_deposit: boolean;
  status: string;
  to_address: string;
  tx_id: string;
}
