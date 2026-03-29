import { Subject, ReplaySubject } from 'rxjs';
enum MessageNames {
  SETLOADING = 'SETLOADING',
  SET_POPUP_LOADING = 'SET_POPUP_LOADING',
  LOGGED_IN = 'LOGGED_IN',
  LOGGED_OUT = 'LOGGED_OUT',
  CLOSE_MODAL = 'CLOSE_MODAL',
  SET_RECAPTCHA = 'SET_RECAPTCHA',
  RESET_RECAPTCHA = 'RESET_RECAPTCHA',
  RESET_SITE_KEY = 'RESET_SITE_KEY',
  SET_SITE_KEY = 'SET_SITE_KEY',
  OPEN_G2FA = 'OPEN_G2FA',
  OPEN_WITHDRAW_VERIFICATION_POPUP = 'OPEN_WITHDRAW_VERIFICATION_POPUP',
  OPEN_TWOFA_AND_EMAILCODE_POPUP = 'OPEN_TWOFA_AND_EMAILCODE_POPUP',
  OPEN_ALERT = 'OPEN_ALERT',
  SET_PAGE_LOADING_WITH_ID = 'SET_PAGE_LOADING_WITH_ID',
  SET_PAGE_FILTERS_WITH_ID = 'SET_PAGE_FILTERS_WITH_ID',

  SETWITHDRAWLOADING = 'SETWITHDRAWLOADING',
  SET_GRID_FILTER = 'SETGRIDFILTER',
  SET_UPLOADER_STATE = 'SET_UPLOADER_STATE',
  UNLOCK_SUBTYPE_SELECT = 'UNLOCK_SUBTYPE_SELECT',
  RESET_IMAGES = 'RESET_IMAGES',
  TOGGLE_SEND_IMAGE_BUTTON = 'TOGGLE_SEND_IMAGE_BUTTON',
  SET_GRID_DATA = 'SET_GRID_DATA',
  RESET_GRID_FILTER = 'RESET_GRID_FILTER',
  ADD_DATA_ROW_TO_GRID = 'ADD_DATA_ROW_TO_GRID',
  ADD_DATA_ROW_TO_WITHDRAWS = 'ADD_DATA_ROW_TO_WITHDRAWS',
  SET_INFINITE_DW_PAGE_DATA = 'SET_INFINITE_DW_PAGE_DATA',
  SET_INITIAL_INFINITE_DW_PAGE_DATA = 'SET_INITIAL_INFINITE_DW_PAGE_DATA',
  SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING = 'SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING',
  SET_DATA_TO_INFINITE_BOTTOM = 'SET_DATA_TO_INFINITE_BOTTOM',
  RESET_INFINITE_SCROLL = 'RESET_INFINITE_SCROLL',
  ADDITIONAL_ACTION = 'ADDITIONAL_ACTION',

  SET_FAVIORITE_ADDRESS = 'SET_FAVIORITE_ADDRESS',
  DELETE_GRID_ROW = 'DELETE_GRID_ROW',

  SET_OPEN_ORDERS_DATA = 'SET_OPEN_ORDERS_DATA',
  NEW_ORDER_NOTIFICATION = 'NEW_ORDER_NOTIFICATION',
  RECONNECT_EVENT = 'RECONNECT_EVENT',
  HIDE_ORDER_NOTIFICATION = 'HIDE_ORDER_NOTIFICATION',
  IS_CANCELING_ORDER = 'IS_CANCELING_ORDER',
  SET_ORDER_HISTORY_DATA = 'SET_ORDER_HISTORY_DATA',
  SET_PAGINATED_ORDER_HISTORY_DATA = 'SET_PAGINATED_ORDER_HISTORY_DATA',
  SET_PAGINATED_TRADE_HISTORY_DATA = 'SET_PAGINATED_TRADE_HISTORY_DATA',
  SET_TRADE_HISTORY_DATA = 'SET_TRADE_HISTORY_DATA',

  SET_BALANCE_PAGE_DATA = 'SET_BALANCE_PAGE_DATA',

  SET_DEPOSIT_PAGE_DATA = 'SET_DEPOSIT_PAGE_DATA',
  SET_FORMER_WITHDRAW_ADDRESSES = 'SET_FORMER_WITHDRAW_ADDRESSES',
  SET_PAGE_LOADING = 'SET_PAGE_LOADING',
  SET_ORDER_DETAIL = 'SET_ORDER_DETAIL',

  UPLOAD_PERCENTAGE = 'UPLOAD_PERCENTAGE',
  SET_UPLOADED_IMAGE = 'SET_UPLOADED_IMAGE',
  DELETE_UPLOADED_IMAGE = 'DELETE_UPLOADED_IMAGE',
  SET_IMAGE_PREVIEW = 'SET_IMAGE_PREVIEW',
  PROFILE_FILE_LOADED = 'PROFILE_FILE_LOADED',

  SET_STEP = 'SET_STEP',
  SET_TAB = 'SET_TAB',

  SET_ERROR = 'SET_ERROR',

  SET_INPUT_ERROR = 'SET_INPUT_ERROR',

  RESIZE = 'RESIZE',
  CHANGE_THEME = 'CHANGE_THEME',
  CHANGE_LAYOUT = 'CHANGE_LAYOUT',

  /////////////
  SET_TRADE_PAGE_CURRENCY_PAIR = 'SET_TRADE_PAGE_CURRENCY_PAIR',
  SELECT_ORDERBOOK_ROW = 'SELECT_ORDERBOOK_ROW',
  IS_LOADING_BUY_SELL = 'IS_LOADING_BUY_SELL',
  SET_CURRENCY_PAIR_DETAILS = 'SET_CURRENCY_PAIR_DETAILS',
  OPEN_LOGIN_POPUP = 'OPEN_LOGIN_POPUP',
  LAYOUT_RESIZE = 'LAYOUT_RESIZE',
  LAYOUT_CHANGE = 'LAYOUT_CHANGE',
  MAIN_CHART_SUMMARY = 'MAIN_CHART_SUMMARY',
  MAIN_CHART_LAST_PRICE = 'MAIN_CHART_LAST_PRICE',
  AUTH_ERROR_EVENT = 'AUTH_ERROR_EVENT',
  /////new order
  ADD_ONE_ORDER_TO_HISTORY = 'ADD_ONE_ORDER_TO_HISTORY',
  SET_DOCUMENT_IMAGES = 'SET_DOCUMENT_IMAGES',

  SET_LOADING_TEST = 'SET_LOADING_TEST',
  SET_LOADING_END = 'SET_LOADING_END',
}

enum DataInjectMessageNames {
  MARKET_TRADES_INITIAL_DATA = 'MARKET_TRADES_INITIAL_DATA',
}

export enum EventMessageNames {
  OPEN_TOOLTIP = 'OPEN_TOOLTIP',
  REFRESH_ORDER_GRID = 'REFRESH_ORDER_GRID',
  GOT_FAV_PAIRS = 'GOT_FAV_PAIRS',
}
interface BroadcastMessage {
  name: MessageNames | EventMessageNames | string;
  value?: any;
  id?: number | string;
  payload?: any;
  additional?: any;
  errorId?: string;
}
interface DataInjectMessage {
  name: DataInjectMessageNames;
  data: any;
}

const Subscriber = new Subject();
const MessageService = {
  send: (msg: BroadcastMessage) => {
    Subscriber.next(msg);
  },
};

const dataInjectSubscriber = new Subject();
const DataInjectMessageService = {
  send: (msg: DataInjectMessage) => {
    dataInjectSubscriber.next(msg);
  },
};

const RepaySubscriber3 = new ReplaySubject(3);
const ReplayMessageService3 = {
  send: (msg: BroadcastMessage) => {
    RepaySubscriber3.next(msg);
  },
};
const SideSubscriber: Subject<any> = new Subject();
const SideMessageService = {
  send: (msg: BroadcastMessage) => {
    SideSubscriber.next(msg);
  },
};
const EventSubscriber: Subject<any> = new Subject();
const EventMessageService = {
  send: (msg: BroadcastMessage) => {
    EventSubscriber.next(msg);
  },
};
const RegisteredUserSubscriber: Subject<any> = new Subject();
const RegisteredUserMessageService = {
  send: (msg: BroadcastMessage) => {
    RegisteredUserSubscriber.next(msg);
  },
};
const MarketTradeSubscriber: Subject<any> = new Subject();
const MarketTradeMessageService = {
  send: (msg: BroadcastMessage) => {
    MarketTradeSubscriber.next(msg);
  },
};
const OrderBookSubscriber: Subject<any> = new Subject();
const OrderBookMessageService = {
  send: (msg: BroadcastMessage) => {
    OrderBookSubscriber.next(msg);
  },
};
const MarketWatchSubscriber: Subject<any> = new Subject();
const MarketWatchMessageService = {
  send: (msg: BroadcastMessage) => {
    MarketWatchSubscriber.next(msg);
  },
};
const TradeChartSubscriber: Subject<any> = new Subject();
const TradeChartMessageService = {
  send: (msg: BroadcastMessage) => {
    TradeChartSubscriber.next(msg);
  },
};

export {
  MessageNames,
  DataInjectMessageNames,
  DataInjectMessageService,
  dataInjectSubscriber,
  BroadcastMessage,
  MessageService,
  Subscriber,
  ReplayMessageService3,
  RepaySubscriber3,
  SideMessageService,
  SideSubscriber,
  //registered user
  RegisteredUserMessageService,
  RegisteredUserSubscriber,
  //marketTrade service
  MarketTradeMessageService,
  MarketTradeSubscriber,
  //Order Book service
  OrderBookMessageService,
  OrderBookSubscriber,
  //Market Watch service
  MarketWatchMessageService,
  MarketWatchSubscriber,
  //Market Watch service
  TradeChartMessageService,
  TradeChartSubscriber,
  //other events
  EventSubscriber,
  EventMessageService,
};
