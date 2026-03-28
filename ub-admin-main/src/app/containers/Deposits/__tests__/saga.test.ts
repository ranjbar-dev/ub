import { put } from 'redux-saga/effects';
import { GetDepositOrders } from '../saga';
import { DepositsActions } from '../slice';
import { MessageService, MessageNames, GridNames } from 'services/messageService';
import { safeApiCall } from 'utils/sagaUtils';

jest.mock('services/userManagementService', () => ({
  GetPaymentAPI: jest.fn(),
}));

jest.mock('services/ordersService', () => ({
  UpdateDepositAPI: jest.fn(),
}));

jest.mock('services/messageService', () => ({
  MessageService: { send: jest.fn() },
  MessageNames: {
    SET_BUTTON_LOADING: 'SET_BUTTON_LOADING',
    UPDATE_GRID_ROW: 'UPDATE_GRID_ROW',
    CLOSE_POPUP: 'CLOSE_POPUP',
  },
  GridNames: { DEPOSITS_PAGE: 'DEPOSITS_PAGE' },
}));

// safeApiCall is a generator — mock must also be a generator
jest.mock('utils/sagaUtils', () => ({
  safeApiCall: jest.fn(function* () {
    return undefined;
  }),
  showErrorToast: jest.fn(),
  showSuccessToast: jest.fn(),
}));

const mockSafeApiCall = safeApiCall as unknown as jest.Mock;

describe('Deposits saga', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('GetDepositOrders', () => {
    it('dispatches setDepositsData on successful response', () => {
      const responseData = { deposits: [{ id: 1, amount: 100 }] };
      const mockResponse = { status: true, data: responseData };

      mockSafeApiCall.mockImplementation(function* () {
        return mockResponse;
      });

      const action = {
        type: 'deposits/GetDepositsAction',
        payload: { page: 1, page_size: 10 },
      };
      const gen = GetDepositOrders(action);

      // yield* safeApiCall(...) delegates — after it resolves, we get the put
      const step1 = gen.next(); // starts the yield* delegation
      // The put action
      expect(step1.value).toEqual(
        put(DepositsActions.setDepositsData(responseData)),
      );
    });

    it('does not dispatch when response is falsy', () => {
      mockSafeApiCall.mockImplementation(function* () {
        return undefined;
      });

      const action = {
        type: 'deposits/GetDepositsAction',
        payload: { page: 1 },
      };
      const gen = GetDepositOrders(action);
      const result = gen.next();

      expect(result.done).toBe(true);
      expect(result.value).toBeUndefined();
    });
  });
});