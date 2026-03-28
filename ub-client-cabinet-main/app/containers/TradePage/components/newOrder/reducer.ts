import {
  amountChangeAction,
  orderActionType,
  orderState,
  percentChangeAction,
  priceChangeAction,
  setInputErrorAction,
  setMainTypeTabAction,
  setPairDetailsAction,
  setSelectedPairNameAction,
  setStopPriceAction,
  setSubTypeTabAction,
  totalChangeAction,
} from './types';
import { divide, multiply, subtract } from 'precise-math';
import { resetAllAndGenerateLabels } from './calculators';

const updateMarketYouGet = ({ state }: { state: orderState }) => {
  if (state.activeSubType === 'market') {
    const fee = state.pairDetails?.fee.takerFee;
    if (state.amount !== '' && state.lastPrice !== '') {
      const amount = +state.amount;
      let eqEmount: number;
      if (state.activeMainType === 'buy') {
        eqEmount = divide(amount, +state.lastPrice);
      } else {
        eqEmount = multiply(amount, +state.lastPrice);
      }
      const marketFeeValue = multiply(eqEmount, fee ?? 0);
      state.tradeFee = marketFeeValue.toString();
      state.youGet = subtract(eqEmount, marketFeeValue).toString();
    } else {
      state.tradeFee = '';
      state.youGet = '';
    }
  }
};

const setValueToYouGet = ({
  state,
  amount,
}: {
  state: orderState;
  amount: string;
}) => {
  const fee =
    state.activeSubType === 'market'
      ? state.pairDetails?.fee?.takerFee
      : state.pairDetails?.fee?.makerFee;

  if (state.activeSubType !== 'market') {
    if (state.activeMainType === 'buy') {
      state.youGet = subtract(+amount, multiply(+amount, fee ?? 0)).toString();
      state.tradeFee = multiply(+amount, fee ?? 0).toString();
    }
  } else {
    updateMarketYouGet({ state });
  }
};

export const ordersReducer = (state: orderState, action: orderActionType) => {
  const newState: orderState = { ...state };

  switch (action.type) {
    case 'setLabels':
      resetAllAndGenerateLabels(newState);
      return { ...newState };

    case 'setMainType':
      changeMainType(newState, action);
      return { ...newState };

    case 'setSubType':
      changeSubType(newState, action);
      return { ...newState };

    case 'amount':
      changeAmount(newState, action);
      return { ...newState };

    case 'slider':
      changeSlider(newState, action);
      return { ...newState };

    case 'price':
      changePrice(action, newState);
      return { ...newState };

    case 'stopPrice':
      changeStopPrice(newState, action);
      return { ...newState };

    case 'total':
      changeTotal(action, newState);
      return { ...newState };

    case 'pairDetails':
      changePairDetails(newState, action);
      return { ...newState };

    case 'resetPairDetails':
      resetPairDetails(newState);
      return { ...newState };

    case 'selectedPairName':
      setSelectedPairName(newState, action);
      return { ...newState };
    case 'updateMarketYouGet':
      newState.lastPrice = action.payload.lastPrice;
      updateMarketYouGet({ state: newState });
      return { ...newState };

    case 'setError':
      setInputError(newState, action);
      return { ...newState };

    default:
      throw new Error();
  }
};

const setInputError = (newState: orderState, action: setInputErrorAction) => {
  switch (action.payload.inputName) {
    case 'amount':
      newState.amountWarning = action.payload.errorText;
      break;
    case 'price':
      newState.priceWarning = action.payload.errorText;
      break;
    case 'stopPrice':
      newState.stopPriceWarning = action.payload.errorText;
      break;

    default:
      break;
  }
};

const setSelectedPairName = (
  newState: orderState,
  action: setSelectedPairNameAction,
) => {
  newState.selectedPair = action.payload;
};

const resetPairDetails = (newState: orderState) => {
  newState.lastPrice = '';
  newState.pairDetails = undefined;
};

const changeStopPrice = (newState: orderState, action: setStopPriceAction) => {
  const v = applyCorrections(newState, action.payload.value, false);
  if (v == '') {
    newState.stopPrice = '';
    return;
  }
  const stopPrice = +v;
  if (stopPrice > 0) {
    newState.stopPriceWarning = '';
  }
  newState.stopPrice = action.payload.value;
};
const changePairDetails = (
  newState: orderState,
  action: setPairDetailsAction,
) => {
  newState.pairDetails = action.payload;
  newState.selectedPair = action.payload.pairData.name;
  resetAllAndGenerateLabels(newState);
};

const changeSubType = (newState: orderState, action: setSubTypeTabAction) => {
  newState.activeSubType = action.payload;
  resetAllAndGenerateLabels(newState);
};

const changeMainType = (newState: orderState, action: setMainTypeTabAction) => {
  newState.activeMainType = action.payload;
  resetAllAndGenerateLabels(newState);
};

const countCharacters = (v: string, char: string): number => {
  return v.split(char).length - 1;
};

const applyCorrections = (
  newState: orderState,
  value: string,
  resetOthersIfEmpty = true,
) => {
  let v = value.replace(/,/g, '');
  if (v == '' && resetOthersIfEmpty) {
    newState.total = '';
    newState.youGet = '';
    newState.tradeFee = '';
    return v;
  }
  if (countCharacters(v, '.') > 1) {
    v = '0';
  }
  if (v.startsWith('.') && v != '.') {
    v = '0' + v;
  }
  if (v.includes('..')) {
    v = '0';
  }
  return v;
};

function changeAmount(newState: orderState, action: amountChangeAction) {
  if (action.payload.fromInput) {
    newState.sliderValue = 0;
  }
  const v = applyCorrections(newState, action.payload.value);
  if (v == '') {
    newState.amount = '';
    if (newState.activeSubType === 'market') {
      newState.youGet = '';
      newState.tradeFee = '';
    }
    return;
  }
  if (isNaN(Number(v))) {
    return;
  }
  const amount = Number(v);
  if (amount > 0) {
    newState.amountWarning = '';
  }
  if (newState.price != '' || newState.activeSubType === 'market') {
    let total: number;
    if (newState.activeSubType !== 'market') {
      total = multiply(amount, +newState.price);
      newState.total = total.toString();
    }
    if (
      newState.activeMainType === 'sell' &&
      newState.activeSubType !== 'market'
    ) {
      total = multiply(amount, +newState.price);

      const fee = newState.pairDetails?.fee.makerFee;

      const newTradeFee = multiply(total, fee ?? 0);
      newState.tradeFee = newTradeFee.toString();
      newState.youGet = subtract(total, fee ?? 0).toString();
    } else {
      // updateInputValue(stream: amountValue, newValue: v);
      newState.amount = v;
      setValueToYouGet({ state: newState, amount: v });
    }
  } else {
    //totalValue.value = '';
  }
  newState.amount = v;
}

function changePrice(action: priceChangeAction, newState: orderState) {
  const v = applyCorrections(newState, action.payload.value);
  if (v == '') {
    newState.price = '';
    return;
  }

  const price = +v;
  if (price > 0) {
    newState.priceWarning = '';
  }
  if (newState.amount !== '') {
    const total = multiply(+newState.amount, price);
    newState.total = total.toString();

    if (
      newState.activeMainType === 'sell' &&
      newState.activeSubType !== 'market'
    ) {
      const fee = newState.pairDetails?.fee.makerFee;
      const newTradeFee = multiply(total, fee ?? 0);
      const newYouGet = subtract(total, fee ?? 0);
      newState.tradeFee = newTradeFee.toString();
      newState.youGet = newYouGet.toString();
    } else {
      setValueToYouGet({ state: newState, amount: newState.amount });
    }
  } else {
    newState.total = '';
    //!next comments will calculate amount correctly, uncomment when you can detect if amount was changed here or before by user
    // if (newState.total != '') {
    //   const total = +newState.total;
    //   newState.amount = divide(total, price).toString();
    //   if (
    //     newState.activeMainType === 'sell' &&
    //     newState.activeSubType !== 'market'
    //   ) {
    //     const fee = newState.pairDetails?.fee.makerFee;
    //     let newTradeFee = multiply(total, fee ?? 0);
    //     let newYouGet = subtract(total, fee ?? 0);
    //     newState.tradeFee = newTradeFee.toString();
    //     newState.youGet = newYouGet.toString();
    //   } else {
    //     setValueToYouGet({ state: newState, amount: newState.amount });
    //   }
    // }
  }
  newState.price = v;
}

function changeTotal(action: totalChangeAction, newState: orderState) {
  const v = applyCorrections(newState, action.payload.value);
  if (v == '') {
    newState.amount = '';
    return;
  }
  const total = +v;
  if (newState.price != '') {
    const price = +newState.price;
    const amount = divide(total, price);
    const newAmount = amount.toString();
    newState.amount = newAmount;
    if (
      newState.activeMainType === 'sell' &&
      newState.activeSubType !== 'market'
    ) {
      const fee = newState.pairDetails?.fee.makerFee;
      const newTradeFee = multiply(total, fee ?? 0);
      newState.tradeFee = newTradeFee.toString();
      newState.youGet = subtract(total, newTradeFee).toString();
    } else {
      setValueToYouGet({ amount: newAmount, state: newState });
    }
  } else {
    //amountValue.value = '';
  }
  if (newState.amount) {
    newState.amountWarning = '';
  }
  if (newState.price) {
    newState.priceWarning = '';
  }
  newState.total = v;
}

function changeSlider(newState: orderState, action: percentChangeAction) {
  let balance: string;
  if (newState.activeSubType === 'market') {
    balance =
      newState.pairDetails?.pairBalances[
        newState.activeMainType === 'buy' ? 0 : 1
      ].balance ?? '0';
  } else {
    balance =
      newState.pairDetails?.pairBalances[
        newState.activeMainType === 'buy' ? 0 : 1
      ].balance ?? '0';
  }
  //calculate valid amount to buy for example BTC in btc-usdt pair when we have some usdt
  if (
    newState.activeMainType === 'buy' &&
    newState.activeSubType === 'limit' &&
    newState.lastPrice !== '' &&
    +(newState.pairDetails?.pairBalances[0].balance ?? 0) > 0
  ) {
    const currentPrice = newState.lastPrice;
    balance = divide(+balance, +currentPrice).toString();
  }

  if (+balance == 0) {
    // _toastToDeposit();
    return;
  }

  if (newState.pairDetails?.sum != null) {
    if (action.payload === 0) {
      newState.amount = '';
      newState.total = '';
      newState.sliderValue = 0;
      newState.youGet = '';
      newState.tradeFee = '';
      return;
    }
    newState.sliderValue = action.payload;
    const amountNumber = multiply(+balance, divide(action.payload, 100));
    const stringAmount = amountNumber.toString();
    changeAmount(newState, {
      type: 'amount',
      payload: { fromInput: false, value: stringAmount },
    });
    //update youGet and trade fee when market
    if (newState.activeSubType === 'market') {
      updateMarketYouGet({ state: newState });
    }
  }
  // update();
  if (newState.sliderValue !== 0) {
    newState.amountWarning = '';
    newState.priceWarning = '';
    newState.stopPriceWarning = '';
  }
}
