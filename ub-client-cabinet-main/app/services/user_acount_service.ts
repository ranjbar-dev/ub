import { apiService } from './api_service';
import { RequestTypes } from './constants';
export const getCountriesAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'main-data/country-list',
    requestType: RequestTypes.GET,
  });
};
export const requestSMSAPI = (data: { phone: string }) => {
  return apiService.fetchData({
    data: data,
    url: 'user/sms-send',
    requestType: RequestTypes.POST,
  });
};
export const verifyCodeAPI = (data: {
  phone: string;
  code: string;
  '2fa_code'?: string;
  password?: string;
}) => {
  return apiService.fetchData({
    data: data,
    url: 'user/sms-enable',
    requestType: RequestTypes.POST,
  });
};
export const getUserProfileAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'user/get-user-profile',
    requestType: RequestTypes.GET,
  });
};
export const updateUserProfileAPI = (data: any) => {
  return apiService.fetchData({
    data: data,
    url: 'user/set-user-profile',
    requestType: RequestTypes.POST,
  });
};
export const get2faQrcodeAPIAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'user/google-2fa-barcode',
    requestType: RequestTypes.GET,
  });
};
export const deleteUserImageAPI = (data: { id: number }) => {
  return apiService.fetchData({
    data,
    url: 'user-profile-image/delete',
    requestType: RequestTypes.POST,
  });
};
