/*
 *
 * TradePage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/TradePage/DEFAULT_ACTION',
  SET_PAIR_MAP = 'app/TradePage/SET_PAIR_MAP',
  GET_CURRENCIES = 'app/TradePage/GET_CURRENCIES',
  GET_INITIAL_MARKET_TRADES_DATA = 'GET_INITIAL_MARKET_TRADES_DATA',
  ADD_REMOVE_FAVORITE_PAIR = 'app/TradePage/ADD_REMOVE_FAVORITE_PAIR',
}
export enum OrderPage {
  OpenOrders = 'OpenOrders',
  OrderHistory = 'OrderHistory',
  TradeHistory = 'TradeHistory',
}
export default ActionTypes;
