/* --- STATE --- */
export interface ExternalExchangeState {
  externalExchangeData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = ExternalExchangeState;
