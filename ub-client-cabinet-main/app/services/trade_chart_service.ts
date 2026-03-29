import { apiService } from './api_service';
import { ChartApiPrefix, RequestTypes } from './constants';
export const getChartConfigAPI = () => {
  return apiService.fetchData({
    data: {},
    url: `${ChartApiPrefix}get-configuration`,
    isRawUrl: true,
    requestType: RequestTypes.GET,
  });
};
