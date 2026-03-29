/*
 * AppConstants
 * Each action has a corresponding type, which the reducer knows and picks up on.
 * To avoid weird typos between the reducer and the actions, we save them as
 * constants here. We prefix them with 'yourproject/YourComponent' so we avoid
 * reducers accidentally picking up actions they shouldn't.
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'App/DEFAULT_ACTION',
  LOGGED_IN_ACTION = 'App/LOGGED_IN_ACTION',
  SET_CURRENCIES_ACTION = 'App/SET_CURRENCIES_ACTION',
  SET_COUNTRIES_ACTION = 'App/SET_COUNTRIES_ACTION',
}
enum Themes {
  DARK = 'darkTheme',
  LIGHT = 'lightTheme',
}
enum AppPages {
  AcountPage = '/account',
  AddressManagement = '/address-management',
  ChangePassword = '/change-password',
  DocumentVerification = '/document-verification',
  Funds = '/funds',
  GoogleAuthentication = '/google-authentication',
  HomePage = '/',
  LoginPage = '/login',
  NotFoundPage = '',
  Orders = '/orders',
  PhoneVerification = '/phone-verification',
  SignupPage = '/signup',
  TradePage = '/trade',
  UpdatePassword = '/auth/forgot-password/update',
  UserInfo = '/userInfo',
  VerifyEmail = '/auth/verify',
  ContactUs = '/contactus',
}

export enum Buttons {
  SimpleRoundButton = 'simpleRoundButton',
  TransParentRoundButton = 'transParentRoundButton',
  SimpleGreyButton = 'simpleGreyButton',
  RoundedRedButton = 'roundedRedButton',
  CancelButton = 'cancelButton',
  Underlined = 'underlined',
  GreenOutlined = 'greenOutlined',
  DensePrimary = 'densePrimary',
}
export enum GridFilterTypes {
  Equals = 'equals',
  Not_Equals = 'notEqual',
  Contains = 'contains',
  Not_Contains = 'notContains',
  Starts_With = 'startsWith',
  Ends_With = 'endsWith',
  Less_Than = 'lessThan',
  Less_Than_or_Equal = 'lessThanOrEqual',
  Greater_Than = 'greaterThan',
  Greater_Than_or_Equal = 'greaterThanOrEqual',
  In_Range = 'inRange',
  Empty = 'empty',
}
export enum CentrifugoChannels {
  MarketTradePrefix = 'trade:trade-book:',
  OrderBookPrefix = 'trade:order-book:',
  TickerChannel = 'trade:ticker',
  TradeChartPrefix = 'trade:kline:',
}
export const GlobalTranslateScope = 'app.globalTitles';
export const GridHeaderNames = 'GridHeaderNames';
export default ActionTypes;
export { Themes, AppPages };
