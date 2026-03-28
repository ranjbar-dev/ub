/* --- STATE --- */
export interface LoginHistoryState {
  loginHistoryData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = LoginHistoryState;
export interface LoginHistoryData {
  createdAt: string;
  device: string;
  email: string;
  id: number;
  ip: string;
  password: string;
  type: StateStrings;
}
export enum StateStrings {
  Successful = 'SUCCESSFUL',
  Failed = 'FAILED',
}
