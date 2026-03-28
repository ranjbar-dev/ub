import React, { useState, useEffect, useCallback, useRef } from 'react';
import styled from 'styles/styled-components';
import translate from 'containers/TradePage/messages';
import TypeTabs from './typeTabs';
import SubTabs from './subTabs';
import { FormattedMessage } from 'react-intl';
import {
  TextField,
  Slider,
  Button,
  InputAdornment,
  Tooltip,
  Zoom,
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
  EventMessageNames,
} from 'services/message_service';
import { NewOrderModel, CurrencyPairDetails } from './types';
import {
  getCurrencyPairInfoAction,
  createNewOrderAction,
} from 'containers/OrdersPage/actions';
import { useInjectReducer } from 'utils/injectReducer';
import { useInjectSaga } from 'utils/injectSaga';
import * as Decimal from 'utils/decimal.js';
import reducer from 'containers/OrdersPage/reducer';
import saga from 'containers/OrdersPage/saga';
import { CurrencyFormater } from 'utils/formatters';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { BS, LMS } from 'containers/OrdersPage/constants';
import FormattedInput from 'react-number-format';
import { savedPairID, savedPairName } from 'utils/sharedData';

const stateSelector = createStructuredSelector({
  loggedIn: makeSelectLoggedIn(),
});
interface NumberFormatCustomProps {
  inputRef: (instance: FormattedInput | null) => void;
  onChange: (event: { target: { name: string; value: string } }) => void;
  name: string;
}
const NumberFormatCustom = (props: NumberFormatCustomProps) => {
  const { inputRef, onChange, ...other } = props;

  return (
    <FormattedInput
      {...other}
      getInputRef={inputRef}
      onValueChange={values => {
        onChange({
          target: {
            name: props.name,
            value: values.value,
          },
        });
      }}
      allowNegative={false}
      thousandSeparator
      isNumericString
    />
  );
};
const NewOrder = () => {
  const firstTofixLength = useRef<number>(0);
  const secoundTofixLength = useRef(0);
  const lastPrice = useRef<number>(0);
  const sendingData = useRef<NewOrderModel>({
    pair_currency_id: savedPairID(),
    user_agent_info: { device: 'web', browser: 'Chrome', os: 'Win32' },
  });

  useInjectReducer({ key: 'ordersPage', reducer: reducer });
  useInjectSaga({ key: 'ordersPage', saga: saga });

  const [BuySell, setBuySell] = useState(BS.buy);
  const [LMSValue, setLMS] = useState(LMS.Limit);
  const [IsLoadingBuySell, setIsLoadingBuySell] = useState(false);
  const [PairName, setPairName] = useState(savedPairName());
  const [TotalValue, setTotalValue] = useState('');
  const [PriceValue, setPriceValue] = useState('');
  const [SliderValue, setSliderValue] = useState(0);
  const [TradeFee, setTradeFee] = useState(0);
  const [StopValue, setStopValue] = useState('');
  const [YouGet, setYouGet]: [any, any] = useState(0);

  const [IsAmountTooltipOpen, setIsAmountTooltipOpen] = useState<boolean>(
    false,
  );
  const [IsPriceTooltipOpen, setIsPriceTooltipOpen] = useState(false);
  const [IsStopPriceTooltipOpen, setIsStopPriceTooltipOpen] = useState(false);
  const amountTooltipText = useRef<string>('');
  const priceTooltipText = useRef<string>('');
  const stopPriceTooltipText = useRef<string>('');
  //@ts-ignore
  const [CurrencyPairDetails, setCurrencyPairDetails]: [
    CurrencyPairDetails,
    any,
  ] = useState({});
  const [Amount, setAmount] = useState('');

  const resetValues = useCallback(() => {
    delete sendingData.current.amount;
    delete sendingData.current.price;
    setTotalValue('');
    setAmount('');
    setStopValue('');
    setPriceValue('');
    setYouGet(0);
    setTradeFee(0);
    setSliderValue(0);
    requestAnimationFrame(() => {
      setIsAmountTooltipOpen(false);

      setIsPriceTooltipOpen(false);

      setIsStopPriceTooltipOpen(false);
    });
  }, []);
  const { loggedIn } = useSelector(stateSelector);
  const dispatch = useDispatch();
  useEffect(() => {
    return () => {
      sendingData.current = {
        pair_currency_id: savedPairID(),
        user_agent_info: { device: 'web', browser: 'Chrome', os: 'Win32' },
      };
    };
  }, []);
  useEffect(() => {
    if (loggedIn === true) {
      dispatch(getCurrencyPairInfoAction({ pair_currency_id: savedPairID() }));
    } else {
      sendingData.current = {
        pair_currency_id: savedPairID(),
        user_agent_info: { device: 'web', browser: 'Chrome', os: 'Win32' },
      };
      setTimeout(() => {
        resetValues();
      }, 0);
      setCurrencyPairDetails({});
    }
    const BSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        sendingData.current.pair_currency_id = message.payload.id;
        sendingData.current.pair_name = message.payload.name;
        setPairName(message.payload.name);
        setCurrencyPairDetails({});
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
    const EventSubscription = EventSubscriber.subscribe((message: any) => {
      if (message.name === EventMessageNames.OPEN_TOOLTIP) {
        if (message.id && message.id === 'amount') {
          amountTooltipText.current = message.payload;
          setIsAmountTooltipOpen(true);
        }
        if (message.id && message.id === 'price') {
          priceTooltipText.current = message.payload;
          setIsPriceTooltipOpen(true);
        }
        if (message.id && message.id === 'stop_point_price') {
          stopPriceTooltipText.current = message.payload;
          setIsStopPriceTooltipOpen(true);
        }
      }
    });
    return () => {
      EventSubscription.unsubscribe();
    };
  }, []);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.IS_LOADING_BUY_SELL) {
        setIsLoadingBuySell(message.payload);
      }
      if (message.name === MessageNames.SET_CURRENCY_PAIR_DETAILS) {
        firstTofixLength.current = message.payload.pairBalances[0].balance.split(
          '.',
        )[1].length;
        secoundTofixLength.current = message.payload.pairBalances[1].balance.split(
          '.',
        )[1].length;
        setCurrencyPairDetails(message.payload);
        //if (Amount) {
        //  setValueToYouGet({ amount: removeComma(Amount) + '' });
        //}
        resetValues();
      } else if (message.name === MessageNames.SELECT_ORDERBOOK_ROW) {
        setBuySell(message.payload.type);
        setTimeout(() => {
          setPriceValue(message.payload.data.price);
        }, 0);

        if (sendingData.current.amount) {
          setTotalValue(
            Decimal(sendingData.current.amount).mul(
              Number(message.payload.data.price),
            ).internal,
          );
          setValueToYouGet({ amount: sendingData.current.amount });
        }
        sendingData.current.type = message.payload.type;
        sendingData.current.price = message.payload.data.price;
      }
      if (message.name === MessageNames.MAIN_CHART_SUMMARY) {
        lastPrice.current = message.payload.payload[0].closePrice;
      }
    });
    return () => {
      Subscription.unsubscribe();
      lastPrice.current = 0;
    };
  }, [CurrencyPairDetails]);
  //////////////////////subscribe to market price
  useEffect(() => {
    resetValues();
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.MAIN_CHART_SUMMARY) {
        updateMarketYouGet();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [LMSValue, BuySell]);
  ////////////////////////////
  const setValueToYouGet = (data: { amount: string | number }) => {
    const fee =
      LMSValue == LMS.Market
        ? CurrencyPairDetails.fee?.takerFee
        : CurrencyPairDetails.fee?.makerFee;
    if (LMSValue !== LMS.Market) {
      if (BuySell === BS.buy) {
        setYouGet(
          Decimal(Number(data.amount)).sub(
            Decimal(Number(data.amount)).mul(fee),
          ).internal,
        );
        setTradeFee(Decimal(Number(data.amount)).mul(fee).internal);
      }
    } else {
      updateMarketYouGet();
    }
  };
  const updateMarketYouGet = () => {
    if (LMSValue === LMS.Market) {
      const fee =
        LMSValue == LMS.Market
          ? CurrencyPairDetails.fee?.takerFee
          : CurrencyPairDetails.fee?.makerFee;

      if (sendingData.current.amount) {
        const amount = sendingData.current.amount;
        let eqEmount;
        if (BuySell === BS.buy) {
          eqEmount = Decimal(Number(amount)).div(Number(lastPrice.current))
            .internal;
        } else {
          eqEmount = Decimal(Number(amount)).mul(Number(lastPrice.current))
            .internal;
        }
        const marketFeeValue = Decimal(eqEmount).mul(fee).internal;
        setTradeFee(marketFeeValue);
        setYouGet(Decimal(eqEmount).sub(Decimal(eqEmount).mul(fee)).internal);
      } else {
        setTradeFee(0);
        setYouGet(0);
      }
    }
  };
  const handleMainTabChange = e => {
    setBuySell(e);
    sendingData.current.type = e;
    resetValues();
  };
  const handleSubTabChange = e => {
    setLMS(e);
    sendingData.current.exchange_type = e;
  };
  const IsNumber = useCallback(val => {
    if (
      !Number(val) &&
      //  !removeComma(val) &&
      val !== '' &&
      val !== '0' &&
      !val.includes('.')
    ) {
      return false;
    }
    return true;
  }, []);
  //  const removeComma = useCallback(
  //    (val: any) => {
  //      let value = val.internal ?? val;
  //      return Number(value.replace(/,/g, ''));
  //    },
  //    [PriceValue, Amount, TotalValue],
  //  );
  const handleAmountChange = e => {
    if (IsAmountTooltipOpen) {
      setIsAmountTooltipOpen(false);
    }
    if (!IsNumber(e.target.value)) {
      return;
    }
    setAmount(e.target.value);

    if (PriceValue) {
      setTotalValue(Decimal(PriceValue).mul(Number(e.target.value)).internal);
      if (BuySell === BS.sell && LMSValue !== LMS.Market) {
        const fee = CurrencyPairDetails.fee?.makerFee;
        const total = Decimal(Number(e.target.value)).mul(PriceValue).internal;
        const newTradeFee = Decimal(total).mul(fee).internal;
        setTradeFee(newTradeFee);
        const newYouGet = Decimal(total).sub(fee).internal;
        setYouGet(newYouGet);
      } else {
        setValueToYouGet({ amount: e.target.value });
      }
    } else {
      setTotalValue('');
    }
    sendingData.current.amount = e.target.value;
    if (LMSValue === LMS.Market) {
      updateMarketYouGet();
    }
  };
  const handlePriceChange = e => {
    if (IsPriceTooltipOpen) {
      setIsPriceTooltipOpen(false);
    }
    if (!IsNumber(e)) {
      return;
    }
    setPriceValue(e);
    if (sendingData.current.amount) {
      setTotalValue(
        Decimal(Number(sendingData.current.amount)).mul(e).internal,
      );
      sendingData.current.price = e;
      if (BuySell === BS.sell && LMSValue !== LMS.Market) {
        const fee = CurrencyPairDetails.fee?.makerFee;
        const total = Decimal(Number(sendingData.current.amount)).mul(
          sendingData.current.price,
        ).internal;
        const newTradeFee = Decimal(total).mul(fee).internal;
        setTradeFee(newTradeFee);
        const newYouGet = Decimal(total).sub(fee).internal;
        setYouGet(newYouGet);
      } else {
        setValueToYouGet({
          amount: sendingData.current.amount + '',
        });
      }
    } else {
      setTotalValue('');
      delete sendingData.current.price;
    }
  };
  const handleTotalChange = e => {
    if (!IsNumber(e.target.value)) {
      return;
    }
    setTotalValue(e.target.value);
    if (PriceValue) {
      const newAmount = Decimal(Number(e.target.value)).div(Number(PriceValue))
        .internal;
      sendingData.current.amount = newAmount;
      setAmount(newAmount);
      if (BuySell === BS.sell && LMSValue !== LMS.Market) {
        const fee = CurrencyPairDetails.fee?.makerFee;
        const total = Decimal(Number(sendingData.current.amount)).mul(PriceValue)
          .internal;
        const newTradeFee = Decimal(total).mul(fee).internal;
        setTradeFee(newTradeFee);
        const newYouGet = Decimal(total).sub(fee).internal;
        setYouGet(newYouGet);
      } else {
        setValueToYouGet({ amount: newAmount });
      }
    } else {
      setAmount('');
    }
  };
  const handleStopValueChange = e => {
    sendingData.current.stop_point_price = e.target.value;
    setStopValue(e.target.value);
    requestAnimationFrame(() => {
      setIsStopPriceTooltipOpen(false);
    });
  };
  const marks = [
    {
      value: 0,
      label: '0%',
    },
    {
      value: 25,
      label: '25%',
    },
    {
      value: 50,
      label: '50%',
    },
    {
      value: 75,
      label: '75%',
    },
    {
      value: 100,
      label: '100%',
    },
  ];
  function sliderValuetext (value: number) {
    return `${value}%`;
  }
  const openLoginPopup = () => {
    MessageService.send({ name: MessageNames.OPEN_LOGIN_POPUP });
  };

  const onSliderChange = (event: any, newValue: number) => {
    if (CurrencyPairDetails.pairBalances && SliderValue !== newValue) {
      let amount = Number(
        CurrencyPairDetails.pairBalances[
          BuySell === BS.buy &&
          (LMSValue === LMS.Market || LMSValue === LMS.StopLimit)
            ? 0
            : 1
        ].balance,
      );
      if (
        BuySell === BS.buy &&
        (LMSValue === LMS.Limit || LMSValue === LMS.StopLimit)
      ) {
        setPriceValue(lastPrice.current + '');
        sendingData.current.price = lastPrice.current + '';
        amount =
          Decimal(Number(CurrencyPairDetails.pairBalances[0].balance))
            .mul(newValue)
            .div(100).internal / lastPrice.current;
      } else {
        amount = Decimal(amount)
          .mul(newValue)
          .div(100).internal;
      }
      sendingData.current.amount = amount + '';
      setAmount(amount + '');
      setValueToYouGet({ amount });
      setSliderValue(newValue);
    }
  };
  const handleBuySell = () => {
    sendingData.current.exchange_type = LMSValue;
    sendingData.current.type = BuySell;
    //@ts-ignore
    if (sendingData.current?.amount?.internal) {
      //@ts-ignore
      sendingData.current.amount = sendingData.current?.amount?.internal;
    }

    if (!sendingData.current.pair_name) {
      sendingData.current.pair_name = PairName;
    }
    if (sendingData.current.exchange_type === LMS.Market) {
      delete sendingData.current.price;
    } else {
      sendingData.current.price = PriceValue + '';
    }
    if (sendingData.current.exchange_type !== LMS.StopLimit) {
      delete sendingData.current.stop_point_price;
    }
    dispatch(createNewOrderAction(sendingData.current));
  };
  return (
    <NewOrderWrapper className={`${BuySell == BS.sell ? 'sell' : ''}`}>
      <div
        className={`selectTypeTabsHeader ${
          loggedIn !== true ? 'disabled' : ''
        }`}
      >
        <TypeTabs onTabChange={handleMainTabChange} />
        <SubTabs onTabChange={handleSubTabChange} />
      </div>
      <div className={`content ${loggedIn !== true ? 'disabled' : ''}`}>
        <div className='available'>
          <div className='availableTitle'>
            <FormattedMessage {...translate.Available} />
          </div>
          <div className='availablevalue'>
            {CurrencyPairDetails.pairBalances &&
              CurrencyFormater(
                CurrencyPairDetails.pairBalances[BuySell === BS.buy ? 0 : 1]
                  .balance,
              )}
            {' ' + PairName.split('-')[BuySell === BS.buy ? 1 : 0]}
          </div>
        </div>
        <div className='amountWrapper'>
          <Tooltip
            PopperProps={{
              disablePortal: true,
            }}
            arrow
            TransitionComponent={Zoom}
            placement='top-start'
            open={IsAmountTooltipOpen}
            disableFocusListener
            disableHoverListener
            disableTouchListener
            className='amountTooltip'
            title={amountTooltipText.current}
          >
            <TextField
              variant='outlined'
              margin='dense'
              fullWidth
              value={Amount}
              error={IsAmountTooltipOpen}
              InputProps={{
                inputComponent: NumberFormatCustom as any,
                endAdornment: (
                  <InputAdornment position='end'>
                    {
                      <span className='endSpan'>
                        {
                          PairName.split('-')[
                            BuySell === BS.buy && LMSValue === LMS.Market
                              ? 1
                              : 0
                          ]
                        }
                      </span>
                    }
                  </InputAdornment>
                ),
              }}
              onChange={handleAmountChange}
              label={
                <>
                  <FormattedMessage {...translate.AmountTo} />{' '}
                  <FormattedMessage
                    {...translate[BuySell === BS.buy ? 'Buy' : 'Sell']}
                  />
                </>
              }
            />
          </Tooltip>
          <Slider
            defaultValue={0}
            value={SliderValue}
            getAriaValueText={sliderValuetext}
            aria-labelledby='discrete-slider-always'
            step={25}
            marks={marks}
            onChange={onSliderChange}
          />
        </div>
        <div className='inputs'>
          {LMSValue === LMS.Limit && (
            <div className='limitInputsWrapper'>
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement='top-start'
                open={IsPriceTooltipOpen}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className='amountTooltip'
                title={priceTooltipText.current}
              >
                <TextField
                  variant='outlined'
                  margin='dense'
                  fullWidth
                  value={PriceValue}
                  error={IsPriceTooltipOpen}
                  onChange={e => handlePriceChange(e.target.value)}
                  label={
                    <>
                      {BuySell === BS.buy ? (
                        <FormattedMessage {...translate.BuyAtMostPrice} />
                      ) : (
                        <FormattedMessage {...translate.SellAtLeastPrice} />
                      )}
                    </>
                  }
                  InputProps={{
                    inputComponent: NumberFormatCustom as any,
                    endAdornment: (
                      <InputAdornment position='end'>
                        {
                          <span className='endSpan'>
                            {PairName.split('-')[1]}
                          </span>
                        }
                      </InputAdornment>
                    ),
                  }}
                />
              </Tooltip>
              <TextField
                variant='outlined'
                margin='dense'
                value={TotalValue}
                onChange={handleTotalChange}
                fullWidth
                InputProps={{
                  inputComponent: NumberFormatCustom as any,
                  endAdornment: (
                    <InputAdornment position='end'>
                      {
                        <span className='endSpan'>
                          {PairName.split('-')[1]}
                        </span>
                      }
                    </InputAdornment>
                  ),
                }}
                label={
                  <>
                    <FormattedMessage {...translate.Total} />{' '}
                  </>
                }
              />
            </div>
          )}
          {LMSValue === LMS.Market && (
            <div className='marketInputsWrapper'>
              <TextField
                disabled
                className='disabledInput'
                fullWidth
                variant='outlined'
                margin='dense'
                label={<FormattedMessage {...translate.AtBestMarketPrice} />}
              ></TextField>
            </div>
          )}
          {LMSValue === LMS.StopLimit && (
            <div className='stopLimitInputsWrapper'>
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement='top-start'
                open={IsStopPriceTooltipOpen}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className='amountTooltip'
                title={stopPriceTooltipText.current}
              >
                <TextField
                  variant='outlined'
                  margin='dense'
                  fullWidth
                  error={IsStopPriceTooltipOpen}
                  value={StopValue}
                  onChange={handleStopValueChange}
                  InputProps={{
                    inputComponent: NumberFormatCustom as any,
                    endAdornment: (
                      <InputAdornment position='end'>
                        {
                          <span className='endSpan'>
                            {PairName.split('-')[1]}
                          </span>
                        }
                      </InputAdornment>
                    ),
                  }}
                  label={<FormattedMessage {...translate.IfPriceRisesTo} />}
                />
              </Tooltip>
              <Tooltip
                PopperProps={{
                  disablePortal: true,
                }}
                arrow
                TransitionComponent={Zoom}
                placement='top-start'
                open={IsPriceTooltipOpen}
                disableFocusListener
                disableHoverListener
                disableTouchListener
                className='amountTooltip'
                title={priceTooltipText.current}
              >
                <TextField
                  variant='outlined'
                  margin='dense'
                  fullWidth
                  error={IsPriceTooltipOpen}
                  value={PriceValue}
                  onChange={e => handlePriceChange(e.target.value)}
                  InputProps={{
                    inputComponent: NumberFormatCustom as any,
                    endAdornment: (
                      <InputAdornment position='end'>
                        {
                          <span className='endSpan'>
                            {PairName.split('-')[1]}
                          </span>
                        }
                      </InputAdornment>
                    ),
                  }}
                  label={
                    <>
                      {BuySell === BS.buy ? (
                        <FormattedMessage {...translate.BuyAtMostPrice} />
                      ) : (
                        <FormattedMessage {...translate.SellAtLeastPrice} />
                      )}
                    </>
                  }
                />
              </Tooltip>
              <TextField
                variant='outlined'
                margin='dense'
                fullWidth
                value={TotalValue}
                onChange={handleTotalChange}
                InputProps={{
                  inputComponent: NumberFormatCustom as any,
                  endAdornment: (
                    <InputAdornment position='end'>
                      {
                        <span className='endSpan'>
                          {PairName.split('-')[BuySell === BS.buy ? 1 : 0]}
                        </span>
                      }
                    </InputAdornment>
                  ),
                }}
                label={
                  <>
                    <FormattedMessage {...translate.Total} />
                  </>
                }
              />
            </div>
          )}
        </div>
      </div>
      {CurrencyPairDetails.fee && CurrencyPairDetails.fee.makerFee && (
        <div className='feeWrapper'>
          <div className='tradeFee row'>
            <div className='tit'>
              <FormattedMessage {...translate.TradeFee} />
            </div>
            <div className='val'>
              {CurrencyFormater(
                Number(TradeFee).toFixed(
                  BuySell === BS.buy
                    ? secoundTofixLength.current
                    : firstTofixLength.current,
                ) + '',
              )}{' '}
              {PairName.split('-')[BuySell === BS.buy ? 0 : 1]}
            </div>
          </div>

          <div
            className='youGet row'
            //style={{ opacity: YouGet.internal ? 1 : 0 }}
          >
            <div className='tit'>
              <FormattedMessage {...translate.YouGet} />
            </div>

            <div className='val'>
              {LMSValue === LMS.Market && YouGet && YouGet > 0 && ' ≈ '}
              {CurrencyFormater(
                Number(YouGet && YouGet > 0 ? YouGet : 0).toFixed(
                  BuySell === BS.buy
                    ? secoundTofixLength.current
                    : firstTofixLength.current,
                ) + '',
              )}{' '}
              {PairName.split('-')[BuySell === BS.buy ? 0 : 1]}
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
              : () => {};
          }}
          style={{
            cursor: IsLoadingBuySell === true ? 'not-allowed' : 'pointer',
          }}
          color={BuySell === BS.buy ? 'primary' : 'secondary'}
          className={`${
            loggedIn === true && BuySell === BS.buy
              ? 'green'
              : loggedIn === true
              ? 'red'
              : ''
          }`}
          fullWidth
          variant='contained'
        >
          {loggedIn == undefined || loggedIn === false ? (
            <FormattedMessage {...translate.LoginToSubmitNewOrder} />
          ) : (
            <IsLoadingWithText
              isLoading={IsLoadingBuySell}
              text={
                <FormattedMessage
                  {...translate[BuySell === BS.buy ? 'Buy' : 'Sell']}
                />
              }
            />
          )}
          {loggedIn == undefined || loggedIn === false ? (
            <div></div>
          ) : (
            IsLoadingBuySell === false && (
              <div style={{ margin: '0 5px' }}> {PairName.split('-')[0]}</div>
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
