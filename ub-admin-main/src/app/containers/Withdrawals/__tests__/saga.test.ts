import { put } from 'redux-saga/effects';
import { GetWithdrawals, GetWithdrawalDetails } from '../saga';
import { WithdrawalsActions } from '../slice';
import { MessageService, MessageNames } from 'services/messageService';
import { safeApiCall } from 'utils/sagaUtils';

jest.mock('services/userManagementService', () => ({
  GetPaymentAPI: jest.fn(),
  GetWithdrawDetailAPI: jest.fn(),
}));

jest.mock('services/messageService', () => ({
  MessageService: { send: jest.fn() },
  MessageNames: {
    SET_MAIN_WITHDRAWALS_ITEM_DETAILS: 'SET_MAIN_WITHDRAWALS_ITEM_DETAILS',
  },
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
const mockSend = MessageService.send as jest.Mock;

describe('Withdrawals saga', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('GetWithdrawals', () => {
    it('dispatches setWithdrawalsData on successful response', () => {
      const responseData = { withdrawals: [{ id: 1, amount: 100 }] };
      const mockResponse = { status: true, data: responseData };

      mockSafeApiCall.mockImplementation(function* () {
        return mockResponse;
      });

      const action = {
        type: 'withdrawals/GetWithdrawals',
        payload: { page: 1, page_size: 10 },
      };
      const gen = GetWithdrawals(action);
      const step1 = gen.next();

      expect(step1.value).toEqual(
        put(WithdrawalsActions.setWithdrawalsData(responseData)),
      );
    });

    it('does not dispatch when response is falsy', () => {
      mockSafeApiCall.mockImplementation(function* () {
        return undefined;
      });

      const action = {
        type: 'withdrawals/GetWithdrawals',
        payload: { page: 1 },
      };
      const gen = GetWithdrawals(action);
      const result = gen.next();

      expect(result.done).toBe(true);
    });
  });

  describe('GetWithdrawalDetails', () => {
    it('sends details message on successful response', () => {
      const responseData = { details: 'withdrawal info' };
      const mockResponse = { status: true, data: responseData };

      mockSafeApiCall.mockImplementation(function* () {
        return mockResponse;
      });

      // Mock getElementById to return null (no loading element in test)
      jest.spyOn(document, 'getElementById').mockReturnValue(null);

      const action = {
        type: 'withdrawals/GetWithdrawalDetailAction',
        payload: { id: '123' },
      };
      const gen = GetWithdrawalDetails(action);

      // Run through the generator
      let result = gen.next();
      while (!result.done) {
        result = gen.next();
      }

      expect(mockSend).toHaveBeenCalledWith({
        name: MessageNames.SET_MAIN_WITHDRAWALS_ITEM_DETAILS,
        payload: { rowData: action.payload, details: responseData },
      });
    });
  });
});