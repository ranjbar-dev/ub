import { apiService } from './apiService';
import { RequestTypes, StandardResponse } from './constants';

export interface UpdateProfileImageStatusParams {
  user_id: number;
  status: string;
}

/**
 * Update the approval status of a user's profile image.
 *
 * @param parameters - User id and new image status
 * @returns Promise with update result
 * @endpoint POST user/profile-image/update
 */
export const UpdateProfileImageStatusAPI = (parameters: UpdateProfileImageStatusParams): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'user/profile-image/update',
    requestType: RequestTypes.POST,
  });
};
