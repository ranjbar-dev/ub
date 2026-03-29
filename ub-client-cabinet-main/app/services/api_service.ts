import { RequestTypes, RequestParameters, BaseUrl } from './constants';
import { queryStringer } from 'utils/formatters';

import { toast } from 'components/Customized/react-toastify';
import { MessageService, MessageNames } from './message_service';
import { cookies, CookieKeys, cookieConfig } from './cookie';
export class ApiService {
  private static instance: ApiService;
  private isRefreshing = false;
  private refreshSubscribers: Array<(token: string) => void> = [];
  private constructor () {}
  public static getInstance (): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService();
    }
    return ApiService.instance;
  }

  private baseUrl = BaseUrl;
  public token: string = '';
  public async fetchData (params: RequestParameters) {
    const url = params.isRawUrl ? params.url : this.baseUrl + params.url;
    //critical to use it this way,
    this.token = cookies.get(CookieKeys.Token) ?? '';

    if (process.env.NODE_ENV !== 'production') {
      console.log(
        `🚀 %c${params.requestType} %crequest to: %c${this.baseUrl}${params.url}\n✉%c:`,
        'color:green;',
        'color:black;',
        'color:green;',
        'color:black;',
        params.data,
      );
    }
    switch (params.requestType) {
      case RequestTypes.GET:
        let query = '';
        if (params.data && Object.keys(params.data).length > 0) {
          query = queryStringer(params.data);
        }
        const rawRes = await fetch(url + query, {
          method: 'GET',
          mode: 'cors',
          credentials: 'omit',
          headers: this.setHeaders(),
        });
        return await this.handleRawResponse(rawRes, params);
      default:
        const rawResponse = await fetch(url, {
          mode: 'cors',
          method: params.requestType,
          headers: this.setHeaders(),
          credentials: 'omit',
          body: JSON.stringify(params.data),
        });
        return await this.handleRawResponse(rawResponse, params);
    }
  }
  handleRawResponse (rawResponse: Response, params: RequestParameters) {
    if (!rawResponse.ok) {
      //if user-pass is wrong or token is expired
      if (rawResponse.status === 401) {
        MessageService.send({ name: MessageNames.SETLOADING, payload: false });
        MessageService.send({ name: MessageNames.AUTH_ERROR_EVENT });
      } else if (rawResponse.status === 403) {
        if (!this.isRefreshing) {
          return this.retryWithNewToken(params);
        }
        // Refresh already in-flight — queue this request to retry once the new token arrives
        return new Promise(resolve => {
          this.refreshSubscribers.push((_token: string) => {
            resolve(this.fetchData(params));
          });
        });
      } else if (rawResponse.status === 500) {
        toast.error('Something Went Wrong!');
      }

      else {
        return rawResponse.json();
      }
    }
    if (process.env.NODE_ENV !== 'production') {
      if (rawResponse.ok) {
        rawResponse
          .clone()
          .json()
          .then(response => {
            console.log(
              `✅ %csuccess %c${params.requestType} %crequest to: %c${this.baseUrl}${params.url}\n✉%c:`,
              'color:green;font-size:15px;',
              'color:blue;',
              'color:black;',
              'color:green;',
              'color:black;',
              params.data,
              '\n',
              ' response 👇',
              response,
            );
          });
      } else {
        console.log(
          `⛔ %cError %c${params.requestType} %crequest to: %c${this.baseUrl}${params.url}\n✉%c:`,
          'color:red;font-size:15px;',
          'color:green;',
          'color:black;',
          'color:green;',
          'color:black;',
          params.data,
        );
        return new Error(`❌ Error calling ${this.baseUrl}${params.url}`);
      }
    }
    return rawResponse.json();
  }

  private setHeaders (): Record<string, string> {
    return {
      'Content-Type': 'application/json',
      ...(this.token !== '' && { Authorization: `Bearer ${this.token}` }),
    };
  }

  private async retryWithNewToken (params: RequestParameters) {
    this.isRefreshing = true;
    try {
      const refreshResponse = await this.fetchData({
        data: { refresh: cookies.get(CookieKeys.RefreshToken) },
        url: 'auth/refresh',
        requestType: RequestTypes.POST,
      });
      if (refreshResponse?.token.length > 0) {
        this.token = refreshResponse.token;
        cookies.set(CookieKeys.Token, refreshResponse.token, cookieConfig());
        cookies.set(
          CookieKeys.RefreshToken,
          refreshResponse.refreshToken,
          cookieConfig(),
        );
        // Resolve all queued requests with the new token
        this.refreshSubscribers.forEach(cb => cb(this.token));
        this.refreshSubscribers = [];
        return await this.fetchData(params);
      }
      // Refresh endpoint returned no token — real auth failure
      this.refreshSubscribers = [];
      MessageService.send({ name: MessageNames.AUTH_ERROR_EVENT });
    } catch {
      this.refreshSubscribers = [];
      MessageService.send({ name: MessageNames.AUTH_ERROR_EVENT });
    } finally {
      this.isRefreshing = false;
    }
  }
}
export const apiService = ApiService.getInstance();
