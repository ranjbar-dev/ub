import {queryStringer} from 'utils/formatters';

import {
	RequestTypes,
	RequestParameters,
	LocalStorageKeys,
	BaseUrl,
	StandardResponse,
	API_TIMEOUT_MS,
	API_MAX_RETRIES,
	RETRYABLE_STATUS_CODES,
} from './constants';
import {MessageService,MessageNames} from './messageService';

/** Get CSRF token from meta tag OR from previous API response header */
function getCsrfToken(): string | null {
	// Try meta tag first (server-rendered)
	const metaTag = document.querySelector('meta[name="csrf-token"]');
	if (metaTag) {
		return metaTag.getAttribute('content');
	}
	// Fall back to stored token from last response
	return localStorage.getItem(LocalStorageKeys.CSRF_TOKEN);
}

/** Represents an API error with status code and optional validation errors. */
export class ApiError extends Error {
	constructor(
		message: string,
		public readonly statusCode: number,
		public readonly errors?: Record<string, string[]>,
		public readonly rawResponse?: unknown,
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export class ApiService {
	private static instance: ApiService;
	private constructor() { }
	public static getInstance(): ApiService {
		if(!ApiService.instance) {
			ApiService.instance=new ApiService();
		}
		return ApiService.instance;
	}

	public baseUrl=BaseUrl;
	public token: string='';

	public async fetchData<T = unknown>(
		params: RequestParameters,
	): Promise<StandardResponse<T>> {
		this.token=localStorage[LocalStorageKeys.ACCESS_TOKEN]
			? localStorage[LocalStorageKeys.ACCESS_TOKEN]
			:'';

		const baseUrl = params.isRawUrl ? BaseUrl : BaseUrl + 'admin/';
		const url = params.requestType === RequestTypes.GET
			? baseUrl + params.url + queryStringer(params.data)
			: baseUrl + params.url;

		if(process.env.NODE_ENV!=='production') {
			console.log(
				`🚀 %c${params.requestType} %crequest to: %c${url}\n✉%c:`,
				'color:green;',
				'color:black;',
				'color:green;',
				'color:black;',
				params.data,
			);
		}

		const content: RequestInit = {
			method: params.requestType,
			headers: this.setHeaders(params.requestType),
			body: params.requestType !== RequestTypes.GET
				? JSON.stringify(params.data)
				: undefined,
		};

		let response: Response | undefined;
		const isIdempotent = params.requestType === RequestTypes.GET || params.requestType === RequestTypes.PUT;
		const maxAttempts = isIdempotent ? API_MAX_RETRIES : 1;

		for (let attempt = 1; attempt <= maxAttempts; attempt++) {
			const controller = new AbortController();
			const timeoutId = setTimeout(() => controller.abort(), API_TIMEOUT_MS);

			try {
				response = await fetch(url, {
					...content,
					signal: params.signal || controller.signal,
				});
				clearTimeout(timeoutId);

				if (RETRYABLE_STATUS_CODES.includes(response.status) && attempt < maxAttempts) {
					const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
					await new Promise(resolve => setTimeout(resolve, backoffMs));
					continue;
				}

				break;
			} catch (fetchError) {
				clearTimeout(timeoutId);

				if (fetchError instanceof DOMException && fetchError.name === 'AbortError') {
					if (attempt < maxAttempts) {
						const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
						await new Promise(resolve => setTimeout(resolve, backoffMs));
						continue;
					}
					throw new ApiError(
						`Request timeout after ${API_TIMEOUT_MS}ms: ${params.requestType} ${params.url}`,
						408,
					);
				}

				if (attempt < maxAttempts) {
					const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
					await new Promise(resolve => setTimeout(resolve, backoffMs));
					continue;
				}

				throw new ApiError(
					`Network error calling ${params.requestType} ${params.url}`,
					0,
				);
			}
		}

		if (!response) {
			throw new ApiError(`No response received: ${params.requestType} ${params.url}`, 0);
		}

		// Store CSRF token if returned in response headers
		const newCsrfToken = response.headers.get('X-CSRF-Token');
		if (newCsrfToken) {
			localStorage.setItem(LocalStorageKeys.CSRF_TOKEN, newCsrfToken);
		}

		const json = await response.json().catch(() => null);

		if (response.status === 200) {
			return json as StandardResponse<T>;
		}

		if (response.status === 401) {
			// Don't try refresh for the refresh endpoint itself or login (avoid infinite loop)
			if (params.requestName === 'refresh-token' || params.url === 'auth/login') {
				MessageService.send({name: MessageNames.SETLOADING, payload: false});
				MessageService.send({name: MessageNames.AUTH_ERROR_EVENT});
				throw new ApiError('Authentication required', 401);
			}

			// Attempt token refresh
			try {
				const refreshToken = localStorage.getItem(LocalStorageKeys.REFRESH_TOKEN) || '';
				if (refreshToken) {
					const refreshResponse = await fetch(BaseUrl + 'admin/auth/refresh', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ refresh_token: refreshToken }),
					});

					if (refreshResponse.ok) {
						const refreshJson = await refreshResponse.json();
						if (refreshJson.token) {
							localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, refreshJson.token);
							if (refreshJson.refresh_token) {
								localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, refreshJson.refresh_token);
							}
							this.token = refreshJson.token;
							// Retry original request with new token
							return this.fetchData<T>(params);
						}
					}
				}
			} catch {
				// Refresh failed — fall through to auth error
			}

			MessageService.send({name: MessageNames.SETLOADING, payload: false});
			MessageService.send({name: MessageNames.AUTH_ERROR_EVENT});
			throw new ApiError('Authentication required', 401);
		}

		if (response.status === 422) {
			throw new ApiError(
				json?.message || 'Validation failed',
				422,
				json?.errors,
				json,
			);
		}

		// All other error codes (500, 403, etc.)
		throw new ApiError(
			json?.message || `API error: ${response.status}`,
			response.status,
			undefined,
			json,
		);
	}

	public setHeaders(requestType?: RequestTypes): Record<string, string> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
		};

		if (this.token !== '') {
			headers.Authorization = `Bearer ${this.token}`;
		}

		// Add CSRF token for non-GET requests
		if (requestType && requestType !== RequestTypes.GET) {
			const csrfToken = getCsrfToken();
			if (csrfToken) {
				headers['X-CSRF-Token'] = csrfToken;
			}
		}

		return headers;
	}
}
export const apiService=ApiService.getInstance();
