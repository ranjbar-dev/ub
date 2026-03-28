/* --- STATE --- */
export interface DepositsState {
  depositsData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = DepositsState;
export enum DepositStatusStrings {
  Completed = 'completed',
  CompletedUpper = 'COMPLETED',
  InProgress = 'in progress',
  Rejected = 'reject',
  RejectedUpper = 'REJECTED',
  Confirmed = 'CONFIRMED',
  Created = 'created',
}
