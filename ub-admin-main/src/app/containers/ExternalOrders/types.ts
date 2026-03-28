/* --- STATE --- */
export interface ExternalOrdersState {
  externalOrdersData: Record<string, unknown> | null;
  netQueueData: Record<string, unknown> | null;
  allQueueData: Record<string, unknown> | null;
  newQueueDetailList: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = ExternalOrdersState;
