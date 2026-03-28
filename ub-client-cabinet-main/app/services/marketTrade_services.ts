import { apiService } from './api_service';
import { RequestTypes } from './constants';

export const getMarketTradesAPI = ({ pair }: { pair: string }) => {
  return apiService.fetchData({
    data: { pair },
    url: 'trade-book',
    requestType: RequestTypes.GET,
  });
};
