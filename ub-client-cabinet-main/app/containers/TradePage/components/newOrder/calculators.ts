import { Currency } from 'containers/AddressManagementPage/types';
import { LocalStorageKeys } from 'services/constants';
import { storage } from 'utils/storage';
import { orderState, tradeInputLabel } from './types';

export const generateLabels = ({ newState }: { newState: orderState }) => {
  const currencies: Currency[] = storage.read(LocalStorageKeys.CURRENCIES);
  const pair = newState.selectedPair;
  const mainType = newState.activeMainType;
  const subType = newState.activeSubType;
  const pair0 = pair.split('-')[0];
  const pair1 = pair.split('-')[1];
  let showDigits0: number = 8;
  let showDigits1: number = 8;
  if (currencies) {
    currencies.forEach(item => {
      if (item.code === pair0) {
        showDigits0 = item.showDigits;
      }
      if (item.code === pair1) {
        showDigits1 = item.showDigits;
      }
    });
  }
  //total label
  const totalL: tradeInputLabel = {
    placeHolder: 'Total',
    endLabel: pair1,
    showDigits: showDigits1,
  };
  newState.totalInputLabel = totalL;
  if (subType === 'limit' || subType === 'stop_limit') {
    ///////////////////////////
    const amountL: tradeInputLabel = {
      placeHolder: 'Amount',
      endLabel: pair0,
      showDigits: showDigits0,
    };
    newState.amountInputLabel = amountL;
    const priceL: tradeInputLabel = {
      placeHolder: 'Price',
      endLabel: pair1,
      showDigits: showDigits1,
    };
    newState.priceInputLabel = priceL;

    const stopL: tradeInputLabel = {
      placeHolder: `${
        mainType === 'buy' ? 'Buy if price reaches' : 'Sell if price reaches'
      }`,
      endLabel: pair1,
      showDigits: showDigits1,
    };
    newState.stopPriceInputLabel = stopL;

    if (subType === 'stop_limit') {
      newState.priceInputLabel.placeHolder =
        mainType === 'buy' ? 'Buy at' : 'Sell at';
    }
    //////////////////////////
  } else if (subType === 'market') {
    const amountL: tradeInputLabel = {
      placeHolder: 'Amount',
      endLabel: mainType === 'buy' ? pair1 : pair0,
      showDigits: mainType === 'buy' ? showDigits1 : showDigits0,
    };
    newState.amountInputLabel = amountL;
  }
};
export const resetAllInputs = (newState: orderState) => {
  newState.amount = '';
  newState.total = '';
  newState.tradeFee = '';
  newState.price = '';
  newState.stopPrice = '';
  newState.youGet = '';
};

export const resetAllAndGenerateLabels = (newState: orderState) => {
  newState.sliderValue = 0;
  newState.amountWarning = '';
  newState.priceWarning = '';
  newState.stopPriceWarning = '';
  resetAllInputs(newState);
  generateLabels({ newState });
};
