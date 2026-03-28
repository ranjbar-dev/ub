import React, { useState, useEffect, useRef, useReducer } from 'react';
import styled from 'styles/styled-components';
import translate from 'containers/TradePage/messages';
import TypeTabs from './typeTabs';
import SubTabs from './subTabs';
import { FormattedMessage } from 'react-intl';
import {
  TextField,
  Slider,
  Button,
  Tooltip,
  Zoom,
  Box,
} from '@material-ui/core';
import { createStructuredSelector } from 'reselect';
import { makeSelectLoggedIn } from 'containers/OrdersPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import {
  MessageNames,
  SideSubscriber,
  MessageService,
  Subscriber,
  EventSubscriber,
} from 'services/message_service';
import { NewOrderModel, orderState, mainType, subType } from './types';
import {
  getCurrencyPairInfoAction,
  createNewOrderAction,
} from 'containers/OrdersPage/actions';
import { useInjectReducer } from 'utils/injectReducer';
import { useInjectSaga } from 'utils/injectSaga';
import reducer from 'containers/OrdersPage/reducer';
import saga from 'containers/OrdersPage/saga';
import { CurrencyFormater, formatSmall } from 'utils/formatters';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { sliderMarks, sliderValuetext } from './sliderConfig';

import { divide } from 'precise-math';
import { currencyMap, savedPairID, savedPairName } from 'utils/sharedData';
import { ordersReducer } from './reducer';
import { cloneDeep } from 'lodash';
import { InputWithFormatter } from './inputWithFormatter';
import { makeSelectPairMap } from 'containers/TradePage/selectors';

const stateSelector = createStructuredSelector({
  loggedIn: makeSelectLoggedIn(),
  pairMap: makeSelectPairMap()
});
const EFields = {
  Amount: 'Amount',
  PriceValue: 'PriceValue',
  TotalValue: 'TotalValue',
  LiveUpdate: 'LiveUpdate',
};

const initialState: orderState = {
  activeMainType: 'buy',
  activeSubType: 'limit',
  sliderValue: 0,
  amount: '',
  tradeFee: '',
  price: '',
  lastPrice: '',
  total: '',
  youGet: '',
  amountWarning: '',
  priceWarning: '',
  stopPriceWarning: '',
  stopPrice: '',
  amountInputLabel: {
    endLabel: '',
    placeHolder: 'Amount',
    showDigits: 8,
  },
  priceInputLabel: {
    endLabel: '',
    placeHolder: 'Price',
    showDigits: 8,
  },
  selectedPair: savedPairName(),
  stopPriceInputLabel: {
    endLabel: '',
    placeHolder: '',
    showDigits: 8,
  },
  totalInputLabel: {
    endLabel: '',
    placeHolder: 'Total',
    showDigits: 8,
  },
};

const NewOrder = () => {
  const { loggedIn, pairMap } = useSelector(stateSelector);
  const [state, setNewState] = useReducer(
    ordersReducer,
    cloneDeep(initialState),
  );
  const coinDigits = {
    buy:
      currencyMap({ code: state.selectedPair.split('-')[1] })?.showDigits ?? 8,
    sell:
      currencyMap({ code: state.selectedPair.split('-')[0] })?.showDigits ?? 8,
  };

  const dispatch = useDispatch();

  const sendingData = useRef<NewOrderModel>({
    pair_currency_id: savedPairID(),
    user_agent_info: { device: 'web', browser: 'Chrome', os: 'Win32' },
  });

  useInjectReducer({ key: 'ordersPage', reducer: reducer });
  useInjectSaga({ key: 'ordersPage', saga: saga });

  const [IsLoadingBuySell, setIsLoadingBuySell] = useState(false);

  useEffect(() => {
    if (loggedIn === true) {
      setNewState({ type: 'resetPairDetails' });
      dispatch(
        getCurrencyPairInfoAction({
          pair_currency_id: savedPairID(),
        }),
      );
    } else {
      sendingData.current = {
        pair_currency_id: sendingData.current.pair_currency_id ?? savedPairID(),
        user_agent_info: { device: 'web', browser: 'Chrome', os: 'Win32' },
      };
      setNewState({ type: 'resetPairDetails' });

      setNewState({ type: 'setLabels' });
    }
    const BSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        sendingData.current.pair_currency_id = message.payload.id;
        sendingData.current.pair_name = message.payload.name;
        setNewState({
          type: 'selectedPairName',
          payload: message.payload.name,
        });

        setNewState({ type: 'resetPairDetails' });
        setNewState({ type: 'setLabels' });

        if (loggedIn === true) {
          dispatch(
            getCurrencyPairInfoAction({ pair_currency_id: message.payload.id }),
          );
        }
      }
    });
    return () => {
      BSubscription.unsubscribe();
    };
  }, [loggedIn]);

  useEffect(() => {
    const refreshData = () => {
      setNewState({ type: 'resetPairDetails' });
      setNewState({ type: 'setLabels' });
      dispatch(getCurrencyPairInfoAction({ pair_currency_id: savedPairID() }));
    };

    const eventSubscription = EventSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RECONNECT_EVENT) {
        refreshData();
      }
    });

    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.IS_LOADING_BUY_SELL) {
        setIsLoadingBuySell(message.payload);
      }

      if (message.name === MessageNames.NEW_ORDER_NOTIFICATION) {
        refreshData();
      }

      if (message.name === MessageNames.SET_CURRENCY_PAIR_DETAILS) {
        setNewState({ type: 'pairDetails', payload: message.payload });
      } else if (message.name === MessageNames.SELECT_ORDERBOOK_ROW) {
        setNewState({ type: 'setMainType', payload: message.payload.type });

        setNewState({
          type: 'price',
          payload: { fromInput: true, value: message.payload.data.price },
        });
      }
    });
    return () => {
      Subscription.unsubscribe();

      eventSubscription.unsubscribe();
    };
  }, []);
  //////////////////////subscribe to market price
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.MAIN_CHART_SUMMARY) {
        setNewState({
          type: 'updateMarketYouGet',
          payload: { lastPrice: message.payload.payload.closePrice },
        });
      }
      if (message.name === MessageNames.MAIN_CHART_LAST_PRICE) {
        setNewState({
          type: 'updateMarketYouGet',
          payload: { lastPrice: message.payload.price.toString() },
        });
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [state.activeSubType]);
  ////////////////////////////

  const handleMainTabChange = (e: mainType) => {
    setNewState({ type: 'setMainType', payload: e });
  };
  const handleSubTabChange = (e: subType) => {
    setNewState({ type: 'setSubType', payload: e });
  };
  const handleAmountChange = (val: string) => {
    setNewState({ type: 'amount', payload: { fromInput: true, value: val } });
  };
  const handlePriceChange = (e: string) => {
    setNewState({ type: 'price', payload: { fromInput: true, value: e } });
  };
  const handleTotalChange = (val: string) => {
    setNewState({ type: 'total', payload: { fromInput: true, value: val } });
  };
  const onSliderChange = (e: any, newValue: number) => {
    if (loggedIn === true) {
      setNewState({ type: 'slider', payload: newValue });
    }
  };
  const handleStopValueChange = (val: string) => {
    setNewState({
      type: 'stopPrice',
      payload: { value: val, fromInput: true },
    });
  };

  const openLoginPopup = () => {
    MessageService.send({ name: MessageNames.OPEN_LOGIN_POPUP });
  };

  const handleBuySell = () => {
    sendingData.current.exchange_type = state.activeSubType;
    sendingData.current.type = state.activeMainType;

    sendingData.current.amount = state.amount;

    if (!state.amount) {
      setNewState({
        type: 'setError',
        payload: { inputName: 'amount', errorText: 'Amount is required' },
      });
      return;
    }

    if (!sendingData.current.pair_name) {
      sendingData.current.pair_name = state.selectedPair;
    }
    if (state.activeSubType === 'market') {
      delete sendingData.current.price;
    } else {
      if (!state.price) {
        setNewState({
          type: 'setError',
          payload: { inputName: 'price', errorText: 'Price is required' },
        });
        return;
      }
      sendingData.current.price = state.price;
    }
    if (state.activeSubType !== 'stop_limit') {
      delete sendingData.current.stop_point_price;
    } else {
      if (!state.stopPrice) {
        setNewState({
          type: 'setError',
          payload: {
            inputName: 'stopPrice',
            errorText: 'Stop price is required',
          },
        });
        return;
      }
      sendingData.current.stop_point_price = state.stopPrice;
    }
    if (
      sendingData.current.price &&
      sendingData.current.amount &&
      state.total
    ) {
      sendingData.current.amount =
        divide(Number(state.total), Number(sendingData.current.price)) + '';
      // }
    }
    // console.log(sendingData);
    dispatch(createNewOrderAction(sendingData.current));
  };
  const isBuy = state.activeMainType === 'buy';


  const priceShowDigits = pairMap[state.selectedPair] ? pairMap[state.selectedPair].showDigits : 8;

  return (
    <NewOrderWrapper
      className={`${state.activeMainType === 'sell' ? 'sell' : ''}`}
    >
      <div
        className={`selectTypeTabsHeader ${loggedIn !== true ? 'disabled' : ''
          }`}
      >
        <TypeTabs onTabChange={handleMainTabChange} />
        <SubTabs onTabChange={handleSubTabChange} />
      </div>
      <div className={`content ${loggedIn !== true ? 'disabled' : ''}`}>
        <div
          className={`available ${state.pairDetails ? 'visible' : 'hiddenWithSizePreserved'
            }`}
        >
          <div className="availableTitle">
            <FormattedMessage {...translate.Available} />
          </div>
          <div className="availablevalue">
            {state.pairDetails &&
              CurrencyFormater(
                Number(
                  state.pairDetails.pairBalances[
                    state.activeMainType === 'buy' ? 0 : 1
                  ].balance,
                ).toFixed(coinDigits[state.activeMainType]),
              )}
            {' ' +
              state.selectedPair.split('-')[
              state.activeMainType === 'buy' ? 1 : 0
              ]}
          </div>
        </div>
        <div className="amountWrapper">
          <Tooltip
            PopperProps={{
              disablePortal: true,
            }}
            arrow
            TransitionComponent={Zoom}
            placement="top-start"
            open={state.amountWarning.length > 0}
            disableFocusListener
            disableHoverListener
            disableTouchListener
            className="amountTooltip"
            title={state.amountWarning}
          >
            <Box>
              <InputWithFormatter
                maxFraction={state.amountInputLabel.showDigits}
                error={state.amountWarning}
                label={state.amountInputLabel.endLabel}
                onChange={handleAmountChange}
                placeholder={state.amountInputLabel.placeHolder}
                value={state.amount}
              />
            </Box>
          </Tooltip>
          <Slider
            defaultValue={0}
            value={state.sliderValue}
            getAriaValueText={sliderValuetext}
            aria-labelledby="discrete-slider-always"
            step={25}
            marks={sliderMarks}
            onChange={onSliderChange}
          />
        </div>
        <div className="inputs">
          {state.activeSubType === 'limit' && (
            <div className="limitInputsWrapper">
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement="top-start"
                open={state.priceWarning.length > 0}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className="amountTooltip"
                title={state.priceWarning}
              >
                <Box>
                  <InputWithFormatter
                    maxFraction={priceShowDigits}
                    error={state.priceWarning}
                    label={state.priceInputLabel.endLabel}
                    onChange={handlePriceChange}
                    placeholder={state.priceInputLabel.placeHolder}
                    value={state.price}
                  />
                </Box>
              </Tooltip>
              <InputWithFormatter
                maxFraction={state.totalInputLabel.showDigits}
                error={''}
                label={state.totalInputLabel.endLabel}
                onChange={handleTotalChange}
                placeholder={state.totalInputLabel.placeHolder}
                value={state.total}
              />
            </div>
          )}
          {state.activeSubType === 'market' && (
            <div className="marketInputsWrapper">
              <TextField
                disabled
                className="disabledInput"
                fullWidth
                variant="outlined"
                margin="dense"
                // label={<FormattedMessage {...translate.AtBestMarketPrice} />}
                label={
                  <span>{`${isBuy ? 'Buy ' : 'Sell '} at market price ${
                    state.lastPrice ? '≈ ' : ''
                  }${
                    state.lastPrice &&
                    formatSmall(state.lastPrice, coinDigits.buy)
                  }`}</span>
                }
              ></TextField>
            </div>
          )}
          {state.activeSubType === 'stop_limit' && (
            <div className="stopLimitInputsWrapper">
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement="top-start"
                open={state.stopPriceWarning.length > 0}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className="amountTooltip"
                title={state.stopPriceWarning}
              >
                <Box>
                  <InputWithFormatter
                    maxFraction={state.stopPriceInputLabel.showDigits}
                    error={state.stopPriceWarning}
                    label={state.stopPriceInputLabel.endLabel}
                    onChange={handleStopValueChange}
                    placeholder={state.stopPriceInputLabel.placeHolder}
                    value={state.stopPrice}
                  />
                </Box>
              </Tooltip>
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement="top-start"
                open={state.priceWarning.length > 0}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className="amountTooltip"
                title={state.priceWarning}
              >
                <Box>
                  <InputWithFormatter
                    maxFraction={priceShowDigits}
                    error={state.priceWarning}
                    label={state.priceInputLabel.endLabel}
                    onChange={handlePriceChange}
                    placeholder={state.priceInputLabel.placeHolder}
                    value={state.price}
                  />
                </Box>
              </Tooltip>
              <InputWithFormatter
                maxFraction={state.totalInputLabel.showDigits}
                error={''}
                label={state.totalInputLabel.endLabel}
                onChange={handleTotalChange}
                placeholder={state.totalInputLabel.placeHolder}
                value={state.total}
              />
            </div>
          )}
        </div>
      </div>
      {state.pairDetails?.fee?.makerFee && (
        <div className="feeWrapper">
          <div className="tradeFee row">
            <div className="tit">
              <FormattedMessage {...translate.TradeFee} />
            </div>
            <div className="val">
              {state.youGet && Number(state.youGet) > 0
                ? formatSmall(
                  state.tradeFee,
                  coinDigits[isBuy ? 'sell' : 'buy'],
                )
                : '0.00'}{' '}
              {state.selectedPair.split('-')[isBuy ? 0 : 1]}
            </div>
          </div>

          <div
            className="youGet row"
            //style={{ opacity: YouGet.internal ? 1 : 0 }}
          >
            <div className="tit">
              <FormattedMessage {...translate.YouGet} />
            </div>

            <div className="val">
              {state.activeSubType === 'market' &&
                //@ts-ignore
                state.youGet &&
                //@ts-ignore
                Number(state.youGet) > 0 &&
                ' ≈ '}
              {state.youGet
                ? formatSmall(state.youGet, coinDigits[isBuy ? 'sell' : 'buy'])
                : '0.00'}{' '}
              {state.selectedPair.split('-')[isBuy ? 0 : 1]}
            </div>
          </div>
        </div>
      )}
      <div
        className={`SubmitWrapper ${loggedIn !== true ? 'toppedSubmit' : ''}`}
      >
        <Button
          onClick={() => {
            loggedIn == undefined || loggedIn === false
              ? openLoginPopup()
              : IsLoadingBuySell === false
                ? handleBuySell()
                : () => { };
          }}
          style={{
            cursor: IsLoadingBuySell === true ? 'not-allowed' : 'pointer',
          }}
          color={isBuy ? 'primary' : 'secondary'}
          className={`${loggedIn === true && isBuy
            ? 'green'
            : loggedIn === true
              ? 'red'
              : ''
            }`}
          fullWidth
          variant="contained"
        >
          {loggedIn == undefined || loggedIn === false ? (
            <FormattedMessage {...translate.LoginToSubmitNewOrder} />
          ) : (
            <IsLoadingWithText
              isLoading={IsLoadingBuySell}
              text={<FormattedMessage {...translate[isBuy ? 'Buy' : 'Sell']} />}
            />
          )}
          {loggedIn == undefined || loggedIn === false ? (
            <></>
          ) : (
            IsLoadingBuySell === false && (
              <div style={{ margin: '0 5px' }}>
                {' '}
                {state.selectedPair.split('-')[0]}
              </div>
            )
          )}
        </Button>
      </div>
    </NewOrderWrapper>
  );
};

export default NewOrder;

const NewOrderWrapper = styled.div`
  .MuiTooltip-tooltipPlacementTop {
    background: var(--orange) !important;
  }
  .MuiTooltip-arrow {
    color: var(--orange) !important;
    left: 0 !important;
  }

  --tabColor: var(--greenText);
  &.sell {
    --tabColor: var(--redText);
  }
  height: 100%;
  .selectTypeTabsHeader {
    width: calc(100% + 10px);
    margin-left: -5px;
  }
  .feeWrapper {
    position: absolute;
    bottom: 55px;
    display: flex;
    flex-direction: column;
    width: 100%;

    .row {
      display: flex;
      place-content: space-between;
      width: calc(100% - 25px);
      margin: 0 8px;
      align-items: center;
      min-height: 19px;
    }
    .tit,
    .val {
      font-size: 11px;
      font-weight: var(--tradeCellFontWeight);
      color: var(--blackText);
      span {
        font-size: 11px;
        font-weight: 500;
        color: var(--blackText);
      }
    }
    .youGet {
      background: var(--oddRows);
    }
  }
  .content {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    padding: 10px 12px;
    position: relative;

    .available {
      display: flex;
      place-content: space-between;
      .availableTitle {
        span {
          color: var(--textGrey);
          font-size: 12px;
          font-weight: 600;
        }
      }
      .availablevalue {
        color: var(--textGrey);
        font-size: 12px;
        font-weight: 600;
        span {
          color: var(--textGrey);
          font-size: 12px;
          font-weight: 500;
        }
      }
    }
  }
  .hiddenWithSizePreserved {
    opacity: 0;
  }
  .MuiSlider-markLabel {
    color: var(--blackText);
    font-size: 11px;
    font-weight: 600;
  }
  .MuiSlider-markLabelActive {
    color: var(--textBlue);
  }
  .MuiSlider-rail {
    background-color: var(--textGrey) !important;
  }
  .MuiSlider-mark {
    width: 10px;
    height: 10px;
    margin-top: -4px;
    margin-left: -5px;
    border-radius: 10px !important;
    background-color: var(--textGrey) !important;
    &.MuiSlider-markActive {
      background-color: var(--textBlue) !important;
    }
  }
  .MuiSlider-thumb.Mui-focusVisible,
  .MuiSlider-thumb:hover {
    box-shadow: 0px 0px 0px 3px rgba(57, 109, 224, 0.4);
  }
  .MuiSlider-root {
    max-width: calc(100% - 12px);
    margin: 0 6px;
    color: var(--textBlue);
  }
  .inputs {
    padding: 12px 0;
  }
  .stopLimitInputsWrapper {
    max-width: 100%;
  }
  .SubmitWrapper {
    position: absolute;
    width: calc(100% - 24px);
    bottom: calc(0% + 12px);
    transition: bottom 0.6s;
    transition-delay: 1.3s;
    left: 14px;
  }
  .endSpan {
    line-height: 1;
    font-size: 13px;
    font-weight: 500;
    margin-top: -1px;
    color: var(--textGrey);
  }
  .loadingCircle {
    top: 8px !important;
  }

  .green {
    background: var(--greenText) !important;
  }
  .red {
    background: var(--redText) !important;
  }
  .primary {
    background: var(--primary) !important;
  }
  .MuiButton-fullWidth {
    min-height: 38px;
  }
  .disabledInput {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .disabledInput:hover {
    fieldset {
      border-color: var(--inputBorderColor) !important;
    }
  }
  input {
    font-weight: 600;
  }
`;
