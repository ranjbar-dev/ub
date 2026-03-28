/* --- STATE --- */
export interface FilledOrdersState {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  filledOrdersData: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = FilledOrdersState;
