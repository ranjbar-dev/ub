/* --- STATE --- */
export interface ReportsState {
  adminReports: Report[] | null;
  withdrawalComments: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState = ReportsState;
export interface Report {
  adminFullName: string;
  adminId: number;
  comment: string;
  createdAt: string;
  id: number;
  isDeleted: false;
  updatedAt: string;
  userFullName: string;
  userId: number;
}
