import { apiService } from "./apiService";
import { RequestTypes, StandardResponse } from "./constants";

export interface ExternalOrdersParams {
	page?: number;
	per_page?: number;
	pair?: string;
	status?: string;
}

export interface NetQueueParams {
	page?: number;
	per_page?: number;
}

export interface ChangeQueueStatusParams {
	id: number;
	status: string;
}

export interface CancelQueueParams {
	id: number;
}

export interface SubmitQueueParams {
	id: number;
}

/**
 * Fetch paginated external orders with optional filters.
 *
 * @param parameters - Pagination, pair, and status filters
 * @returns Promise with external orders list
 * @endpoint GET exchange/order
 */
export const GetExternalOrdersAPI = (parameters: ExternalOrdersParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/order',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated net queue entries.
 *
 * @param parameters - Pagination options
 * @returns Promise with net queue list
 * @endpoint GET exchange/order/queue
 */
export const GetNetQueueAPI = (parameters: NetQueueParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/order/queue',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch all queue entries with pagination.
 *
 * @param parameters - Pagination options
 * @returns Promise with all queue entries
 * @endpoint GET exchange/order/queue/all
 */
export const GetAllQueueAPI = (parameters: NetQueueParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/order/queue/all',
		requestType: RequestTypes.GET,
	});
};

/**
 * Change the status of a net queue aggregation entry.
 *
 * @param parameters - Entry id and new status
 * @returns Promise with status change result
 * @endpoint POST exchange/aggregation/change-status
 */
export const ChangeNetQueueStatusAPI = (parameters: ChangeQueueStatusParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/aggregation/change-status',
		requestType: RequestTypes.POST,
	});
};

/**
 * Cancel a net queue order.
 *
 * @param parameters - Order id to cancel
 * @returns Promise with cancellation result
 * @endpoint POST exchange/order/change-status
 */
export const CancelNetQueueAPI = (parameters: CancelQueueParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/order/change-status',
		requestType: RequestTypes.POST,
	});
};

/**
 * Submit a net queue order for processing.
 *
 * @param parameters - Order id to submit
 * @returns Promise with submission result
 * @endpoint POST exchange/order/change-status
 */
export const SubmitNetQueueAPI = (parameters: SubmitQueueParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/order/change-status',
		requestType: RequestTypes.POST,
	});
};