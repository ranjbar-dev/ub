export enum RequestTypes {
	PUT = 'PUT',
	POST = 'POST',
	GET = 'GET',
	DELETE = 'DELETE',
}
export interface RequestParameters<T = Record<string, unknown>> {
	requestType: RequestTypes;
	url: string;
	data: T;
	isRawUrl?: boolean;
	requestName?: string;
	signal?: AbortSignal;
}
export enum LocalStorageKeys {
	ACCESS_TOKEN = 'access_token',
	REFRESH_TOKEN = 'refresh_token',
	CSRF_TOKEN = 'csrf_token',
	CURRENCIES = 'currencies',
	Theme = 'theme',
	COUNTRIES = 'countries',
	Managers = 'Managers',
	LAYOUT_NAME = 'ln',
	PAIRS = 'pairs',
	SELECTED_COIN = 'selectedCoin',
	FUND_PAGE = 'fp',
	SHOW_TOP_INFO = 'sti',
	TRADELAYOUT = 'tl',
	FAV_PAIRS = 'fps',
	FAV_COIN = 'fc',
	SHOW_FAVS = 'sf',
	TIME_FRAME = 'timeframe',
	CHANNEL = 'chan',
	VERIFICATION_WINDOW_TYPE = 'VERIFICATION_WINDOW_TYPE',
	VISIBLE_ORDER_SECTION = 'vos',
}
export interface StandardResponse<T = unknown> {
	status: boolean;
	message: string;
	data: T;
	token?: string;
	errors?: Record<string, string[]>;
}
export enum UploadUrls {
	USER_PROFILE_IMAGE = 'user-profile-image/upload',
}

export const appUrl = 'https://admin.unitedbit.com';
export const BaseUrl: string =
	process.env.REACT_APP_API_BASE_URL || appUrl + '/api/v1/';
const prefix = process.env.NODE_ENV === 'development' ? 'dev-' : ''
export const webAppAddress = `https://${prefix}app.unitedbit.com/api/v1/`

/** Default request timeout in milliseconds */
export const API_TIMEOUT_MS = 30000;

/** Maximum retry attempts for idempotent requests */
export const API_MAX_RETRIES = 3;

/** Retryable HTTP status codes */
export const RETRYABLE_STATUS_CODES = [408, 429, 502, 503, 504];