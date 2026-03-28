import { apiService } from './apiService';
import { RequestTypes, StandardResponse } from './constants';

export interface AddAdminCommentParams {
  user_id: number;
  comment: string;
}

export interface DeleteAdminCommentParams {
  id: number;
}

export interface EditAdminCommentParams {
  id: number;
  comment: string;
}

export interface WithdrawalCommentsParams {
  user_id: number;
  page?: number;
  per_page?: number;
}

export interface UpdateFinancialMethodParams {
  id: number;
  [key: string]: unknown;
}

export interface UpdateCurrencyPairParams {
  id: number;
  [key: string]: unknown;
}

export interface GetCommitionsParams {
  user_id: number;
}

/**
 * Add an admin comment to a user.
 *
 * @param parameters - User id and comment text
 * @returns Promise with comment creation result
 * @endpoint POST user/admin-comment/add
 */
export const AddAdminCommentAPI = (parameters: AddAdminCommentParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'user/admin-comment/add',
    requestType: RequestTypes.POST,
  });
};

/**
 * Delete an admin comment by id.
 *
 * @param parameters - Comment identifier
 * @returns Promise with deletion result
 * @endpoint POST user/admin-comment/delete
 */
export const DeleteAdminCommentAPI = (parameters: DeleteAdminCommentParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'user/admin-comment/delete',
    requestType: RequestTypes.POST,
  });
};

/**
 * Edit an existing admin comment.
 *
 * @param parameters - Comment id and updated text
 * @returns Promise with edit result
 * @endpoint POST user/admin-comment/update
 */
export const EditAdminCommentAPI = (parameters: EditAdminCommentParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'user/admin-comment/update',
    requestType: RequestTypes.POST,
  });
};

/**
 * Fetch withdrawal comments for a user.
 *
 * @param parameters - User id with optional pagination
 * @returns Promise with withdrawal comments
 * @endpoint GET payment/user-comments
 */
export const GetWithdrawalCommentsAPI = (parameters: WithdrawalCommentsParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'payment/user-comments',
    requestType: RequestTypes.GET,
  });
};

/**
 * Update a financial currency method.
 *
 * @param parameters - Method id and fields to update
 * @returns Promise with update result
 * @endpoint POST currency/update
 */
export const UpdateFinancialMethodAPI = (parameters: UpdateFinancialMethodParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'currency/update',
    requestType: RequestTypes.POST,
  });
};

/**
 * Update a currency pair configuration.
 *
 * @param parameters - Pair id and fields to update
 * @returns Promise with update result
 * @endpoint POST currency/update-pair
 */
export const UpdateCurrencyPairAPI = (parameters: UpdateCurrencyPairParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'currency/update-pair',
    requestType: RequestTypes.POST,
  });
};

/**
 * Fetch commission statistics for a user.
 *
 * @param data - User identifier
 * @returns Promise with user commission statistics
 * @endpoint GET statistic/user-statistic
 */
export const GetCommitionsAPI = (data: GetCommitionsParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: data as unknown as Record<string, unknown>,
    url: 'statistic/user-statistic',
    requestType: RequestTypes.GET,
  });
};
