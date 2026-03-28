import { apiService } from './api_service';
import { RequestTypes } from './constants';
export const getChartConfigAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'http://116.203.76.196/tv/api/v1/js/get-configuration',
    isRawUrl: true,
    requestType: RequestTypes.GET,
  });
};
