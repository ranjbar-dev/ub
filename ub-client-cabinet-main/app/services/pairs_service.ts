import { apiService } from './api_service';
import { RequestTypes } from './constants';

export const addRemoveFavPairAPI = (data: { pair_currency_id: number, action: 'add' | 'remove' }) => {
  return apiService.fetchData({
    data,
    url: 'currencies/favorite',
    requestType: RequestTypes.POST,
  });
};
export const getFavPairAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'currencies/favorite-pairs',
    requestType: RequestTypes.GET,
  });
};
export const getPairsListAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'currencies/pairs-list',
    requestType: RequestTypes.GET,
  });
};
