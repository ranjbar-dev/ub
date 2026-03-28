/* --- STATE --- */
export interface AdminsState {
  adminsData: unknown[] | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = AdminsState;
