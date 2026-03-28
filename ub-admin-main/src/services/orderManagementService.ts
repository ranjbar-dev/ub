import { apiService } from './apiService';
import { RequestTypes, StandardResponse } from './constants';

export interface LiquidityOrdersParams {
  page?: number;
  per_page?: number;
  pair?: string;
}

export interface UpdateCommissionReportParams {
  id: number;
  [key: string]: unknown;
}

/**
 * Fetch paginated liquidity commission report orders.
 *
 * @param parameters - Pagination and optional pair filter
 * @returns Promise with liquidity orders
 * @endpoint GET exchange/order/commission-report
 */
export const GetLiquidityOrdersAPI = (parameters: LiquidityOrdersParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'exchange/order/commission-report',
    requestType: RequestTypes.GET,
  });
};

/**
 * Update a commission report entry.
 *
 * @param parameters - Report id and fields to update
 * @returns Promise with update result
 * @endpoint POST exchange/order/update-commission-report
 */
export const UpdateCommissionReportAPI = (parameters: UpdateCommissionReportParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'exchange/order/update-commission-report',
    requestType: RequestTypes.POST,
  });
};
