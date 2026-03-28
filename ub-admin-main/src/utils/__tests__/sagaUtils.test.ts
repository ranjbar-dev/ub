import { call } from 'redux-saga/effects';
import { safeApiCall, showSuccessToast, showErrorToast } from '../sagaUtils';
import { ApiError } from 'services/apiService';
import { MessageService, MessageNames } from 'services/messageService';

jest.mock('services/messageService', () => ({
	MessageService: { send: jest.fn() },
	MessageNames: {
		TOAST: 'TOAST',
		SETLOADING: 'SETLOADING',
		SET_INPUT_ERROR: 'SET_INPUT_ERROR',
		AUTH_ERROR_EVENT: 'AUTH_ERROR_EVENT',
	},
}));

const mockSend = MessageService.send as jest.Mock;

describe('safeApiCall', () => {
	const mockApiFunc = jest.fn();
	const params = { id: 42 };

	beforeEach(() => {
		mockSend.mockClear();
		mockApiFunc.mockClear();
	});

	describe('success flow', () => {
		it('yields call effect and returns response when status is true', () => {
			const response = { status: true, data: { result: 'ok' } };
			const gen = safeApiCall(mockApiFunc, params);

			expect(gen.next().value).toEqual(call(mockApiFunc, params));

			const result = gen.next(response as any);
			expect(result.done).toBe(true);
			expect(result.value).toEqual(response);
		});

		it('does not show error toast on success', () => {
			const response = { status: true, data: {} };
			const gen = safeApiCall(mockApiFunc, params);
			gen.next();
			gen.next(response as any);

			const toastCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.TOAST,
			);
			expect(toastCalls).toHaveLength(0);
		});
	});

	describe('loading indicator', () => {
		it('shows loading before call and hides it after (success)', () => {
			const loadingId = 'myButton';
			const response = { status: true, data: {} };
			const gen = safeApiCall(mockApiFunc, params, { loadingId });

			// First next: SETLOADING(true) is called synchronously before the yield
			gen.next();
			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.SETLOADING,
				value: true,
				loadingId,
			});

			// Second next: SETLOADING(false) called in finally
			gen.next(response as any);
			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.SETLOADING,
				value: false,
				loadingId,
			});
		});

		it('hides loading even on failure', () => {
			const loadingId = 'myButton';
			const response = { status: false, message: 'Something went wrong' };
			const gen = safeApiCall(mockApiFunc, params, { loadingId });

			gen.next();
			gen.next(response as any);

			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.SETLOADING,
				value: false,
				loadingId,
			});
		});

		it('hides loading even on thrown error', () => {
			const loadingId = 'myButton';
			const gen = safeApiCall(mockApiFunc, params, { loadingId });

			gen.next();
			gen.throw!(new Error('Network error'));

			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.SETLOADING,
				value: false,
				loadingId,
			});
		});

		it('does not send SETLOADING when loadingId is not provided', () => {
			const response = { status: true, data: {} };
			const gen = safeApiCall(mockApiFunc, params);

			gen.next();
			gen.next(response as any);

			const loadingCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.SETLOADING,
			);
			expect(loadingCalls).toHaveLength(0);
		});
	});

	describe('API failure (status: false)', () => {
		it('shows error toast with response message', () => {
			const response = { status: false, message: 'Operation failed' };
			const gen = safeApiCall(mockApiFunc, params);

			gen.next();
			const result = gen.next(response as any);

			expect(result.done).toBe(true);
			expect(result.value).toBeUndefined();
			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.TOAST,
				value: 'Operation failed',
				type: 'error',
			});
		});

		it('uses errorMessage option when response has no message', () => {
			const response = { status: false };
			const gen = safeApiCall(mockApiFunc, params, {
				errorMessage: 'Custom error',
			});

			gen.next();
			gen.next(response as any);

			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.TOAST,
				value: 'Custom error',
				type: 'error',
			});
		});

		it('uses fallback message when no message or errorMessage', () => {
			const response = { status: false };
			const gen = safeApiCall(mockApiFunc, params);

			gen.next();
			gen.next(response as any);

			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.TOAST,
				value: 'Operation failed',
				type: 'error',
			});
		});

		it('does not show toast when toastOnError is false', () => {
			const response = { status: false, message: 'error' };
			const gen = safeApiCall(mockApiFunc, params, { toastOnError: false });

			gen.next();
			gen.next(response as any);

			const toastCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.TOAST,
			);
			expect(toastCalls).toHaveLength(0);
		});
	});

	describe('ApiError 401', () => {
		it('returns undefined without showing toast', () => {
			const gen = safeApiCall(mockApiFunc, params);
			const error = new ApiError('Unauthorized', 401);

			gen.next();
			const result = gen.throw!(error);

			expect(result.done).toBe(true);
			expect(result.value).toBeUndefined();

			const toastCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.TOAST,
			);
			expect(toastCalls).toHaveLength(0);
		});
	});

	describe('ApiError 422', () => {
		it('dispatches SET_INPUT_ERROR with validation errors', () => {
			const errors = { email: ['Email is invalid'], name: ['Name is required'] };
			const gen = safeApiCall(mockApiFunc, params);
			const error = new ApiError('Validation error', 422, errors);

			gen.next();
			const result = gen.throw!(error);

			expect(result.done).toBe(true);
			expect(result.value).toBeUndefined();
			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.SET_INPUT_ERROR,
				value: errors,
			});
		});

		it('does not show error toast for 422 validation errors', () => {
			const errors = { field: ['bad'] };
			const gen = safeApiCall(mockApiFunc, params);
			const error = new ApiError('Validation error', 422, errors);

			gen.next();
			gen.throw!(error);

			const toastCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.TOAST,
			);
			expect(toastCalls).toHaveLength(0);
		});
	});

	describe('generic Error', () => {
		it('shows generic error toast and returns undefined', () => {
			const gen = safeApiCall(mockApiFunc, params);

			gen.next();
			const result = gen.throw!(new Error('Network error'));

			expect(result.done).toBe(true);
			expect(result.value).toBeUndefined();
			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.TOAST,
				value: 'An unexpected error occurred',
				type: 'error',
			});
		});

		it('uses custom errorMessage for generic errors', () => {
			const gen = safeApiCall(mockApiFunc, params, {
				errorMessage: 'Something broke',
			});

			gen.next();
			gen.throw!(new Error('Network error'));

			expect(mockSend).toHaveBeenCalledWith({
				name: MessageNames.TOAST,
				value: 'Something broke',
				type: 'error',
			});
		});

		it('does not show toast when toastOnError is false', () => {
			const gen = safeApiCall(mockApiFunc, params, { toastOnError: false });

			gen.next();
			gen.throw!(new Error('Network error'));

			const toastCalls = mockSend.mock.calls.filter(
				([msg]: [any]) => msg.name === MessageNames.TOAST,
			);
			expect(toastCalls).toHaveLength(0);
		});
	});
});

describe('showSuccessToast', () => {
	beforeEach(() => {
		mockSend.mockClear();
	});

	it('calls MessageService.send with TOAST name and success type', () => {
		showSuccessToast('Operation succeeded');

		expect(mockSend).toHaveBeenCalledTimes(1);
		expect(mockSend).toHaveBeenCalledWith({
			name: MessageNames.TOAST,
			value: 'Operation succeeded',
			type: 'success',
		});
	});

	it('forwards the message string as value', () => {
		showSuccessToast('Custom success message');

		expect(mockSend).toHaveBeenCalledWith(
			expect.objectContaining({ value: 'Custom success message' }),
		);
	});
});

describe('showErrorToast', () => {
	beforeEach(() => {
		mockSend.mockClear();
	});

	it('calls MessageService.send with TOAST name and error type', () => {
		showErrorToast('Something went wrong');

		expect(mockSend).toHaveBeenCalledTimes(1);
		expect(mockSend).toHaveBeenCalledWith({
			name: MessageNames.TOAST,
			value: 'Something went wrong',
			type: 'error',
		});
	});

	it('forwards the message string as value', () => {
		showErrorToast('Custom error message');

		expect(mockSend).toHaveBeenCalledWith(
			expect.objectContaining({ value: 'Custom error message' }),
		);
	});
});
