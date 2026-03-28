interface NewOrderModel {
  type?: string;
  exchange_type?: 'market' | 'limit' | 'stop_limit';
  pair_currency_id?: number;
  price?: string;
  amount?: string;
  pair_name?: string;
  stop_point_price?: string;
  user_agent_info: { device: 'web'; browser: 'Chrome'; os: 'Win32' };
}

interface CurrencyPairDetails {
  pairBalances: PairBalance[];
  pairData: PairData;
  chart: ChartData[];
  sum: string;
  fee: Fee;
}
interface PairBalance {
  currencyId: number;
  currencyCode: string;
  currencyName: string;
  balance: string;
}
interface PairData {
  id: number;
  name: string;
  minimumOrderAmount: string;
}
interface ChartData {
  amount: string;
  equivalentAmount: string;
  name: string;
  percent: string;
}
interface Fee {
  makerFee: number;
  takerFee: number;
}

interface inputChangePayload {
  payload: {
    value: string;
    fromInput: boolean = false;
  };
}
export interface amountChangeAction extends inputChangePayload {
  type: 'amount';
}
export interface priceChangeAction extends inputChangePayload {
  type: 'price';
}
export interface totalChangeAction extends inputChangePayload {
  type: 'total';
}
export interface lastPriceAction extends inputChangePayload {
  type: 'lastPrice';
}
export interface setStopPriceAction extends inputChangePayload {
  type: 'stopPrice';
}
export interface percentChangeAction {
  type: 'slider';
  payload: number;
}
export interface setPairDetailsAction {
  type: 'pairDetails';
  payload: CurrencyPairDetails;
}
export interface resetPairDetailsAction {
  type: 'resetPairDetails';
}
export interface updateYouGetAction {
  type: 'updateMarketYouGet';
  payload:{
    lastPrice:string
  }
}
export interface setLabelsAction {
  type: 'setLabels';
}

export interface setSelectedPairNameAction {
  type: 'selectedPairName';
  payload: string;
}

export type mainType = 'buy' | 'sell';
export type subType = 'limit' | 'market' | 'stop_limit';
export interface tradeInputLabel {
  placeHolder: string;
  endLabel: string;
  showDigits: number;
}
export interface orderState {
  activeMainType: mainType;
  activeSubType: subType;
  amount: string;
  price: string;
  total: string;
  tradeFee: string;
  sliderValue: number;
  youGet: string;
  lastPrice: string;
  pairDetails?: CurrencyPairDetails;
  amountWarning: string;
  priceWarning: string;
  stopPriceWarning: string;
  stopPrice: string;
  totalInputLabel: tradeInputLabel;
  priceInputLabel: tradeInputLabel;
  amountInputLabel: tradeInputLabel;
  stopPriceInputLabel: tradeInputLabel;
  selectedPair: string;
}
export type orderActionType =
  | amountChangeAction
  | priceChangeAction
  | totalChangeAction
  | percentChangeAction
  | setMainTypeTabAction
  | setSubTypeTabAction
  | setPairDetailsAction
  | setStopPriceAction
  | resetPairDetailsAction
  | setSelectedPairNameAction
  | updateYouGetAction
  | setInputErrorAction
  | setLabelsAction
  | lastPriceAction;

export type setMainTypeTabAction = {
  type: 'setMainType';
  payload: mainType;
};

export type setInputErrorAction = {
  type: 'setError';
  payload: {
    inputName: 'amount' | 'price' | 'stopPrice';
    errorText: string;
  };
};

export type setSubTypeTabAction = {
  type: 'setSubType';
  payload: subType;
};

interface orderReducerProps {
  state: orderState;
  action: orderActionType;
}

export { NewOrderModel, CurrencyPairDetails, NewOrderModel };
