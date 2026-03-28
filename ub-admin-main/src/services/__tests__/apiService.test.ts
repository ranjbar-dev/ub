import { ApiService, ApiError } from '../apiService';
import { RequestTypes, LocalStorageKeys, BaseUrl } from '../constants';
import { MessageService, MessageNames } from '../messageService';

// ── Mocks ─────────────────────────────────────────────────────────────────────

jest.mock('../messageService', () => ({
	MessageService: { send: jest.fn() },
	MessageNames: jest.requireActual('../messageService').MessageNames,
}));

const mockSend = MessageService.send as jest.Mock;

// ── Helpers ───────────────────────────────────────────────────────────────────

/** Build a minimal fetch Response-like object. */
const makeFetchResponse = (status: number, body: object = {}) => ({
	status,
	ok: status >= 200 && status < 300,
	json: jest.fn().mockResolvedValue(body),
	headers: { get: jest.fn().mockReturnValue(null) },
});

const SUCCESS_BODY = { status: true, message: 'ok', data: { id: 1 } };
const ADMIN_BASE = BaseUrl + 'admin/';

// ── Setup ─────────────────────────────────────────────────────────────────────

beforeEach(() => {
	// Reset the singleton so each test gets a fresh instance
	(ApiService as any).instance = undefined;

	// Seed localStorage with a valid access token
	localStorage.clear();
	localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, 'test-access-token');

	// Default: every fetch call succeeds
	global.fetch = jest.fn().mockResolvedValue(makeFetchResponse(200, SUCCESS_BODY));

	mockSend.mockClear();
});

afterEach(() => {
	jest.restoreAllMocks();
});

// ── Singleton ─────────────────────────────────────────────────────────────────

describe('singleton', () => {
	it('getInstance() always returns the same object', () => {
		const a = ApiService.getInstance();
		const b = ApiService.getInstance();
		expect(a).toBe(b);
	});

	it('reads the access token from localStorage on first fetchData call', async () => {
		const svc = ApiService.getInstance();
		await svc.fetchData({ requestType: RequestTypes.GET, url: 'ping', data: {} });
		expect(svc.token).toBe('test-access-token');
	});
});

// ── Successful requests ───────────────────────────────────────────────────────

describe('successful requests', () => {
	let svc: ApiService;
	beforeEach(() => { svc = ApiService.getInstance(); });

	it('GET: builds correct URL with query string and returns parsed JSON', async () => {
		const result = await svc.fetchData({
			requestType: RequestTypes.GET,
			url: 'users',
			data: { page: 1, limit: 10 },
		});

		const [calledUrl, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledUrl).toContain(ADMIN_BASE + 'users');
		expect(calledUrl).toContain('page=1');
		expect(calledOpts.method).toBe('GET');
		expect(calledOpts.headers['Authorization']).toBe('Bearer test-access-token');
		expect(result).toEqual(SUCCESS_BODY);
	});

	it('POST: sends JSON body and sets Content-Type', async () => {
		await svc.fetchData({
			requestType: RequestTypes.POST,
			url: 'users',
			data: { name: 'Alice' },
		});

		const [calledUrl, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledUrl).toBe(ADMIN_BASE + 'users');
		expect(calledOpts.method).toBe('POST');
		expect(calledOpts.headers['Content-Type']).toBe('application/json');
		expect(JSON.parse(calledOpts.body)).toEqual({ name: 'Alice' });
	});

	it('PUT: sends body', async () => {
		await svc.fetchData({
			requestType: RequestTypes.PUT,
			url: 'users/1',
			data: { name: 'Bob' },
		});

		const [, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledOpts.method).toBe('PUT');
		expect(JSON.parse(calledOpts.body)).toEqual({ name: 'Bob' });
	});

	it('DELETE: uses DELETE method', async () => {
		await svc.fetchData({
			requestType: RequestTypes.DELETE,
			url: 'users/1',
			data: { id: 1 },
		});

		const [calledUrl, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledUrl).toBe(ADMIN_BASE + 'users/1');
		expect(calledOpts.method).toBe('DELETE');
		// DELETE sends body (same as POST/PUT in this implementation)
		expect(JSON.parse(calledOpts.body)).toEqual({ id: 1 });
	});

	it('passes AbortSignal to fetch for timeout control', async () => {
		await svc.fetchData({ requestType: RequestTypes.GET, url: 'ping', data: {} });

		const [, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledOpts.signal).toBeInstanceOf(AbortSignal);
	});

	it('omits Authorization header when no token is stored', async () => {
		localStorage.clear(); // no access_token
		(ApiService as any).instance = undefined;
		svc = ApiService.getInstance();

		await svc.fetchData({ requestType: RequestTypes.GET, url: 'ping', data: {} });

		const [, calledOpts] = (global.fetch as jest.Mock).mock.calls[0];
		expect(calledOpts.headers['Authorization']).toBeUndefined();
	});
});

// ── Error handling ────────────────────────────────────────────────────────────

describe('error handling', () => {
	let svc: ApiService;
	beforeEach(() => { svc = ApiService.getInstance(); });

	it('400 response → throws ApiError with statusCode 400', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(
			makeFetchResponse(400, { message: 'Bad Request' }),
		);

		await expect(
			svc.fetchData({ requestType: RequestTypes.GET, url: 'bad', data: {} }),
		).rejects.toMatchObject({ name: 'ApiError', statusCode: 400 });
	});

	it('422 response → throws ApiError with validation errors', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(
			makeFetchResponse(422, {
				message: 'Validation failed',
				errors: { email: ['Invalid email'] },
			}),
		);

		const error = await svc
			.fetchData({ requestType: RequestTypes.POST, url: 'users', data: {} })
			.catch(e => e);

		expect(error).toBeInstanceOf(ApiError);
		expect(error.statusCode).toBe(422);
		expect(error.errors).toEqual({ email: ['Invalid email'] });
	});

	it('500 response → throws ApiError with statusCode 500', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(
			makeFetchResponse(500, { message: 'Internal Server Error' }),
		);

		await expect(
			svc.fetchData({ requestType: RequestTypes.GET, url: 'crash', data: {} }),
		).rejects.toMatchObject({ name: 'ApiError', statusCode: 500 });
	});
});

// ── Retry logic ───────────────────────────────────────────────────────────────
// Spy on Math.pow so backoff = Math.min(1000 * 0, 8000) = 0 ms — keeps tests fast
// while still exercising the full retry code-path.

describe('retry logic', () => {
	let svc: ApiService;

	beforeEach(() => {
		jest.spyOn(Math, 'pow').mockReturnValue(0);
		svc = ApiService.getInstance();
	});

	it('GET retries up to 3 times on 503 and returns on final success', async () => {
		let calls = 0;
		(global.fetch as jest.Mock).mockImplementation(() => {
			calls++;
			if (calls < 3) return Promise.resolve(makeFetchResponse(503));
			return Promise.resolve(makeFetchResponse(200, SUCCESS_BODY));
		});

		const result = await svc.fetchData({
			requestType: RequestTypes.GET,
			url: 'flaky',
			data: {},
		});

		expect(global.fetch).toHaveBeenCalledTimes(3);
		expect(result).toEqual(SUCCESS_BODY);
	});

	it('GET exhausts all 3 retries and throws the last ApiError', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(makeFetchResponse(503));

		await expect(
			svc.fetchData({ requestType: RequestTypes.GET, url: 'always-503', data: {} }),
		).rejects.toMatchObject({ name: 'ApiError', statusCode: 503 });

		expect(global.fetch).toHaveBeenCalledTimes(3);
	});

	it('POST with 503 does NOT retry (non-idempotent)', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(makeFetchResponse(503));

		await expect(
			svc.fetchData({ requestType: RequestTypes.POST, url: 'submit', data: {} }),
		).rejects.toMatchObject({ name: 'ApiError', statusCode: 503 });

		expect(global.fetch).toHaveBeenCalledTimes(1);
	});

	it('GET retries on network error and succeeds on retry', async () => {
		let calls = 0;
		(global.fetch as jest.Mock).mockImplementation(() => {
			calls++;
			if (calls < 2) return Promise.reject(new Error('Network failure'));
			return Promise.resolve(makeFetchResponse(200, SUCCESS_BODY));
		});

		const result = await svc.fetchData({
			requestType: RequestTypes.GET,
			url: 'net-error',
			data: {},
		});

		expect(global.fetch).toHaveBeenCalledTimes(2);
		expect(result).toEqual(SUCCESS_BODY);
	});

	it('GET with 408 retries', async () => {
		let calls = 0;
		(global.fetch as jest.Mock).mockImplementation(() => {
			calls++;
			if (calls < 2) return Promise.resolve(makeFetchResponse(408));
			return Promise.resolve(makeFetchResponse(200, SUCCESS_BODY));
		});

		await svc.fetchData({ requestType: RequestTypes.GET, url: 'timeout-once', data: {} });

		expect(global.fetch).toHaveBeenCalledTimes(2);
	});
});

// ── 401 token-refresh flow ────────────────────────────────────────────────────

describe('401 token refresh', () => {
	let svc: ApiService;
	beforeEach(() => { svc = ApiService.getInstance(); });

	it('401 + valid refresh token → refreshes, updates token, retries original request', async () => {
		localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, 'old-refresh');

		let callCount = 0;
		(global.fetch as jest.Mock).mockImplementation((url: string) => {
			callCount++;
			// First call: original request → 401
			if (callCount === 1) return Promise.resolve(makeFetchResponse(401));
			// Second call: refresh endpoint → success
			if (url.includes('auth/refresh')) {
				return Promise.resolve(
					makeFetchResponse(200, { token: 'new-token', refresh_token: 'new-refresh' }),
				);
			}
			// Third call: retried original request → 200
			return Promise.resolve(makeFetchResponse(200, SUCCESS_BODY));
		});

		const result = await svc.fetchData({
			requestType: RequestTypes.GET,
			url: 'protected',
			data: {},
		});

		expect(result).toEqual(SUCCESS_BODY);
		expect(localStorage.getItem(LocalStorageKeys.ACCESS_TOKEN)).toBe('new-token');
		expect(mockSend).not.toHaveBeenCalledWith(
			expect.objectContaining({ name: MessageNames.AUTH_ERROR_EVENT }),
		);
	});

	it('401 + refresh endpoint fails → emits AUTH_ERROR_EVENT', async () => {
		localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, 'old-refresh');

		let callCount = 0;
		(global.fetch as jest.Mock).mockImplementation((url: string) => {
			callCount++;
			if (callCount === 1) return Promise.resolve(makeFetchResponse(401));
			if (url.includes('auth/refresh')) {
				return Promise.resolve(makeFetchResponse(401, { message: 'Refresh invalid' }));
			}
			return Promise.resolve(makeFetchResponse(200, SUCCESS_BODY));
		});

		await expect(
			svc.fetchData({ requestType: RequestTypes.GET, url: 'protected', data: {} }),
		).rejects.toMatchObject({ statusCode: 401 });

		expect(mockSend).toHaveBeenCalledWith(
			expect.objectContaining({ name: MessageNames.AUTH_ERROR_EVENT }),
		);
	});

	it('401 + no refresh token → emits AUTH_ERROR_EVENT immediately', async () => {
		localStorage.removeItem(LocalStorageKeys.REFRESH_TOKEN);
		(global.fetch as jest.Mock).mockResolvedValue(makeFetchResponse(401));

		await expect(
			svc.fetchData({ requestType: RequestTypes.GET, url: 'protected', data: {} }),
		).rejects.toMatchObject({ statusCode: 401 });

		expect(mockSend).toHaveBeenCalledWith(
			expect.objectContaining({ name: MessageNames.AUTH_ERROR_EVENT }),
		);
	});

	it('refresh-token request getting 401 does NOT loop — emits AUTH_ERROR_EVENT', async () => {
		(global.fetch as jest.Mock).mockResolvedValue(makeFetchResponse(401));

		await expect(
			svc.fetchData({
				requestType: RequestTypes.POST,
				url: 'auth/refresh',
				data: {},
				requestName: 'refresh-token',
			}),
		).rejects.toMatchObject({ statusCode: 401 });

		// Must only call fetch once (no infinite refresh loop)
		expect(global.fetch).toHaveBeenCalledTimes(1);
		expect(mockSend).toHaveBeenCalledWith(
			expect.objectContaining({ name: MessageNames.AUTH_ERROR_EVENT }),
		);
	});
});
