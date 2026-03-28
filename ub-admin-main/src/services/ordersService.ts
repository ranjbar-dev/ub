import { apiService } from './apiService';
import { RequestTypes, StandardResponse } from './constants';

export interface CancelOrderParams {
	order_id: number;
}

export interface FullFillOrderParams {
	order_id: number;
}

export interface UpdateDepositParams {
	id: number;
	status: string;
}

export interface BalancePaginationParams {
	page?: number;
	per_page?: number;
}

export interface InternalTransferParams {
	from_currency: string;
	to_currency: string;
	amount: number;
	user_id: number;
}

/**
 * Cancel an open order.
 *
 * @param parameters - Order identifier
 * @returns Promise with cancellation result
 * @endpoint POST order/cancel
 */
export const CancelOrderAPI = (parameters: CancelOrderParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'order/cancel',
		requestType: RequestTypes.POST,
	});
};

/**
 * Fulfill an open order.
 *
 * @param parameters - Order identifier
 * @returns Promise with fulfillment result
 * @endpoint POST order/fulfill
 */
export const FullFillOrderAPI = (parameters: FullFillOrderParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'order/fulfill',
		requestType: RequestTypes.POST,
	});
};

/**
 * Update the status of a deposit.
 *
 * @param parameters - Deposit id and new status
 * @returns Promise with update result
 * @endpoint POST payment/update-deposit
 */
export const UpdateDepositAPI = (parameters: UpdateDepositParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment/update-deposit',
		requestType: RequestTypes.POST,
	});
};

/**
 * Fetch paginated crypto balances.
 *
 * @param parameters - Pagination options
 * @returns Promise with balance records
 * @endpoint GET crypto-balance
 */
export const GetBalancesAPI = (parameters: BalancePaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'crypto-balance',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated internal transfer history.
 *
 * @param parameters - Pagination options
 * @returns Promise with transfer history records
 * @endpoint GET crypto-internal-transfer/list
 */
export const GetBalanceHistoryAPI = (parameters: BalancePaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'crypto-internal-transfer/list',
		requestType: RequestTypes.GET,
	});
};

/**
 * Trigger a full balance update for all crypto balances.
 *
 * @param parameters - Flexible update payload
 * @returns Promise with update result
 * @endpoint POST crypto-balance/update-balance
 */
export const UpdateAllBalancesAPI = (parameters: Record<string, unknown>): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters,
		url: 'crypto-balance/update-balance',
		requestType: RequestTypes.POST,
	});
};

/**
 * Create an internal crypto transfer between currencies.
 *
 * @param parameters - Source/target currencies, amount, and user id
 * @returns Promise with transfer result
 * @endpoint POST crypto-internal-transfer/create
 */
export const InternalTransferAPI = (parameters: InternalTransferParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'crypto-internal-transfer/create',
		requestType: RequestTypes.POST,
	});
};
