import {apiService} from './apiService';
import {RequestTypes, StandardResponse} from './constants';

export interface PaginationParams {
	page?: number;
	per_page?: number;
	sort_by?: string;
	sort_dir?: 'asc' | 'desc';
}

export interface UserIdParam {
	user_id: number;
}

export interface UpdateUserPermissionsParams {
	user_id: number;
	permissions: string[];
}

export interface AdminReportsParams extends PaginationParams {
	user_id?: number;
}

export interface LoginHistoryParams extends PaginationParams {
	user_id: number;
}

export interface SyncTicksParams {
	pair: string;
	from?: number;
	to?: number;
}

export interface ScanBlockParams {
	coin?: string;
}

export interface MarketTicksParams {
	pair?: string;
	timeframe?: string;
}

export interface WithdrawDetailParams {
	id: number;
}

export interface AddPaymentCommentParams {
	payment_id: number;
	comment: string;
}

export interface UpdateWithdrawParams {
	id: number;
	status: string;
}

/**
 * Fetch paginated list of user accounts.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with list of user accounts
 * @endpoint GET user/
 */
export const GetUserAccountsAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch initial data for a single user.
 *
 * @param parameters - User identifier
 * @returns Promise with user details
 * @endpoint GET user/show
 */
export const GetInitialUserDataAPI = (parameters: UserIdParam): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/show',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch balance information for a user.
 *
 * @param parameters - User identifier
 * @returns Promise with user balances
 * @endpoint GET user/balances
 */
export const GetUserBalancesAPI = (parameters: UserIdParam): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/balances',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch whitelisted withdrawal addresses for a user.
 *
 * @param parameters - User identifier
 * @returns Promise with withdrawal address list
 * @endpoint GET user/withdraw-addresses
 */
export const GetUserWhiteAddressesAPI = (parameters: UserIdParam): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/withdraw-addresses',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch permissions assigned to a user.
 *
 * @param parameters - User identifier
 * @returns Promise with user permissions
 * @endpoint GET user/permissions
 */
export const GetUserPermissionsAPI = (parameters: UserIdParam): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/permissions',
		requestType: RequestTypes.GET,
	});
};

/**
 * Update permissions for a user.
 *
 * @param parameters - User id and new permissions array
 * @returns Promise with update result
 * @endpoint POST user/update-permissions
 */
export const UpdateUserPermissionsAPI = (parameters: UpdateUserPermissionsParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/update-permissions',
		requestType: RequestTypes.POST,
	});
};

/**
 * Fetch paginated billing/payment grid data.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with billing records
 * @endpoint GET payment/
 */
export const GetBillingGridDataAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch admin reports/comments list.
 *
 * @param parameters - Pagination options and optional user filter
 * @returns Promise with admin reports
 * @endpoint GET user/admin-comment/list
 */
export const GetAdminReportsAPI = (parameters: AdminReportsParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/admin-comment/list',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated list of open orders.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with open orders
 * @endpoint GET order/
 */
export const GetOpenOrdersAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'order/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated order history (closed/filled orders).
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with order history records
 * @endpoint GET order/history
 */
export const GetOrderHistoryAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'order/history',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated trade history.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with trade history records
 * @endpoint GET trade/
 */
export const GetTradeHistoryAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'trade/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch profile images for a user.
 *
 * @param parameters - User identifier
 * @returns Promise with user profile images
 * @endpoint GET user/profile-images
 */
export const GetUserImagesAPI = (parameters: UserIdParam): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/profile-images',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch login history for a user.
 *
 * @param parameters - User id with pagination options
 * @returns Promise with login history records
 * @endpoint GET user/login-history
 */
export const GetLoginHistoryAPI = (parameters: LoginHistoryParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'user/login-history',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch paginated payment records.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with payment records
 * @endpoint GET payment/
 */
export const GetPaymentAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch available finance/currency methods.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with currency methods
 * @endpoint GET currency/
 */
export const GetFinanceMethodsAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'currency/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch available currency pairs.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with currency pairs
 * @endpoint GET currency/pair
 */
export const GetCurrencyPairsAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'currency/pair',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch OHLC sync list.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with sync list data
 * @endpoint GET ohlc/sync
 */
export const GetSyncListAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'ohlc/sync',
		requestType: RequestTypes.GET,
	});
};

/**
 * Trigger a blockchain block scan.
 *
 * @param parameters - Optional coin to scan
 * @returns Promise with scan result
 * @endpoint POST wallet/block/scan
 */
export const ScanBlockAPI = (parameters: ScanBlockParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'wallet/block/scan',
		requestType: RequestTypes.POST,
	});
};

/**
 * Create an OHLC tick sync job.
 *
 * @param parameters - Pair name and optional time range
 * @returns Promise with sync job result
 * @endpoint POST ohlc/create-sync
 */
export const SyncTicksAPI = (parameters: SyncTicksParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'ohlc/create-sync',
		requestType: RequestTypes.POST,
	});
};

/**
 * Fetch external exchange data.
 *
 * @param parameters - Pagination and sorting options
 * @returns Promise with external exchange records
 * @endpoint GET exchange/
 */
export const GetExternalExchangeAPI = (parameters: PaginationParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'exchange/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Fetch market OHLC ticks.
 *
 * @param parameters - Optional pair and timeframe filters
 * @returns Promise with market tick data
 * @endpoint GET ohlc/
 */
export const GetMarketTicksAPI = (parameters: MarketTicksParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'ohlc/',
		requestType: RequestTypes.GET,
	});
};

/**
 * Update user data from admin page.
 *
 * @param parameters - Flexible user data fields to update
 * @returns Promise with update result
 * @endpoint POST user/on-page-update
 */
export const UpdateUserDataAPI = (parameters: Record<string, unknown>): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters,
		url: 'user/on-page-update',
		requestType: RequestTypes.POST,
	});
};

/**
 * Fetch withdrawal detail by id.
 *
 * @param parameters - Withdrawal record identifier
 * @returns Promise with withdrawal details
 * @endpoint GET payment/withdraw-detail
 */
export const GetWithdrawDetailAPI = (parameters: WithdrawDetailParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment/withdraw-detail',
		requestType: RequestTypes.GET,
	});
};

/**
 * Add a comment to a payment record.
 *
 * @param parameters - Payment id and comment text
 * @returns Promise with comment creation result
 * @endpoint POST payment-comment/new
 */
export const AddPaymentCommentAPI = (parameters: AddPaymentCommentParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment-comment/new',
		requestType: RequestTypes.POST,
	});
};

/**
 * Update a withdrawal record status.
 *
 * @param parameters - Withdrawal id and new status
 * @returns Promise with update result
 * @endpoint POST payment/update-withdraw
 */
export const UpdateWithdrawAPI = (parameters: UpdateWithdrawParams): Promise<StandardResponse> => {
	return apiService.fetchData({
		data: parameters as unknown as Record<string, unknown>,
		url: 'payment/update-withdraw',
		requestType: RequestTypes.POST,
	});
};
