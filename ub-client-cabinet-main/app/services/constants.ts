import { IClientOptions } from 'mqtt';

export enum RequestTypes {
  PUT = 'PUT',
  POST = 'POST',
  GET = 'GET',
  DELETE = 'DELETE',
}
export interface RequestParameters {
  requestType: RequestTypes;
  url: string;
  data: any;
  isRawUrl?: boolean;
  requestName?: string;
}
export enum LocalStorageKeys {
  ACCESS_TOKEN = 'access_token',
  REFRESH_TOKEN = 'refresh_token',
  USERNAME = 'username',
  RECAPTCHA = 'recall',
  SITEKEY = 'refreshToken',
  PASSWORD = 'password',
  CURRENCIES = 'currencies',
  CURRENCY_MAP = 'currencyMap',
  PAIRS_MAP = 'pairsMap',
  Theme = 'theme',
  CanSendNewEmail = 'csne',
  COUNTRIES = 'countries',
  LAYOUT_NAME = 'ln',
  SELECTED_COIN = 'selectedCoin',
  FUND_PAGE = 'fp',
  SHOW_TOP_INFO = 'sti',
  TRADELAYOUT = 'tl',
  FAV_PAIRS = 'fps',
  FAV_COIN = 'fc',
  TRADE_CONFIGS = 'tc',
  SHOW_FAVS = 'sf',
  TIME_FRAME = 'timeframe',
  CHANNEL = 'chan',
  VISIBLE_ORDER_SECTION = 'vos',
  SAVED_TRADE_PAIR = 'stp',
}
export enum SessionStorageKeys {
  SITE_KEY = 'SK',
}
export interface StandardResponse {
  status: boolean;
  message: string;
  data: any;
}
export enum UploadUrls {
  USER_PROFILE_IMAGE = 'user-profile-image/upload',
}
export const MqttProtocol:
  | 'wss'
  | 'ws'
  | 'mqtt'
  | 'mqtts'
  | 'tcp'
  | 'ssl'
  | 'wx'
  | 'wxs'
  | undefined = 'wss';
const productionAddress = 'app.unitedbit.com';
const devAddress = 'dev-app.unitedbit.com';

const mobileProductionAddress = 'https://m.unitedbit.com';
const mobileDevAddress = 'https://dev-m.unitedbit.com';

export const LandingPageAddress = 'https://www.unitedbit.com';
export const pre = 'https';
export const mainUrl =
  process.env.NODE_ENV == 'development' || process.env.IS_DEV_BUILD === 'true'
    ? devAddress
    : productionAddress;
export const minWidthToRedirectToMobile = 1000;
export const mobileUrl =
  process.env.NODE_ENV == 'development' || process.env.IS_DEV_BUILD === 'true'
    ? mobileDevAddress
    : mobileProductionAddress;

export const appUrl = `${pre}://${mainUrl}`;
export const BaseUrl = appUrl + '/api/v1/';
export const tradingView = `${appUrl}/tv/api/v1/main-route`;
export const ChartApiPrefix = `${appUrl}/tv/api/v1/js/`;
export const jsAPI = `${appUrl}/tv/api/v1/js`;
//TODO: change next line "productionAddress" to "mainUrl" when dev-mqtt is fixed
export const mqttServer = `${MqttProtocol}://${mainUrl}:8443`;

export const MqttAdditionalConfig: IClientOptions = {
  protocol: MqttProtocol,
  connectTimeout: 30 * 60 * 1000,
  reconnectPeriod: 2 * 1000,
  keepalive: 0,
};
