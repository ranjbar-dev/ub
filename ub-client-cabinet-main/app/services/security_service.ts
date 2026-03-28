import { apiService } from './api_service';
import { RequestTypes } from './constants';
import { ChangePasswordModel } from 'containers/ChangePassword/types';
import { RegisterModel } from 'containers/SignupPage/types';
import { SetG2FaModel } from 'containers/GoogleAuthenticationPage/types';

export const loginAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'auth/login',
    requestType: RequestTypes.POST,
    requestName: 'login',
  });
};
export const getUserDataAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'user/user-data',
    requestType: RequestTypes.GET,
  });
};
export const getNewVerificationEmailAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'user/send-verification-email',
    requestType: RequestTypes.POST,
  });
};
export const changePasswordAPI = (data: ChangePasswordModel) => {
  return apiService.fetchData({
    data: data,
    url: 'user/change-password',
    requestType: RequestTypes.POST,
  });
};
export const registerAPI = (data: RegisterModel) => {
  return apiService.fetchData({
    data: data,
    url: 'auth/register',
    requestType: RequestTypes.POST,
  });
};
export const set2FaAPI = (data: SetG2FaModel) => {
  return apiService.fetchData({
    data: { code: data.code, password: data.password },
    url: `user/google-2fa-${data.setEnable === true ? 'enable' : 'disable'}`,
    requestType: RequestTypes.POST,
  });
};
export const getRecapchaKeyAPI = () => {
  return apiService.fetchData({
    data: {},
    url: `main-data/common`,
    requestType: RequestTypes.GET,
  });
};
export const acountActivationAPI = (data: { code: string }) => {
  return apiService.fetchData({
    data: data,
    url: `auth/verify`,
    requestType: RequestTypes.POST,
  });
};
export const forgotPasswordAPI = (data: { email: string }) => {
  return apiService.fetchData({
    data: data,
    url: `auth/forgot-password`,
    requestType: RequestTypes.POST,
  });
};
export const resetPasswordAPI = (data: { email: string; code: string }) => {
  return apiService.fetchData({
    data: data,
    url: `auth/forgot-password/update`,
    requestType: RequestTypes.POST,
  });
};
