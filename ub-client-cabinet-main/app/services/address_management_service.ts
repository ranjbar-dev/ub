import { apiService } from './api_service';
import { RequestTypes } from './constants';
export const getCurrenciesAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'currencies',
    requestType: RequestTypes.GET,
  });
};
export const getWithDrawAddressesAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'withdraw-address',
    requestType: RequestTypes.GET,
  });
};
export const addNewWithDrawAddressAPI = (data: {
  address: string;
  code: string;
  label: string;
  network?: string;
}) => {
  return apiService.fetchData({
    data: data,
    url: 'withdraw-address/new',
    requestType: RequestTypes.POST,
  });
};
export const deleteWithDrawAddressAPI = (data: { ids: number[] }) => {
  return apiService.fetchData({
    data: data,
    url: 'withdraw-address/delete',
    requestType: RequestTypes.POST,
  });
};
export const setFavoriteWithDrawAddressAPI = (data: {
  action: string;
  id: number;
}) => {
  return apiService.fetchData({
    data: data,
    url: 'withdraw-address/favorite',
    requestType: RequestTypes.POST,
  });
};
