/*
 *
 * TradeHeader
 *
 */

import React, { memo, useEffect, useRef, useState } from 'react';
import { FormattedMessage } from 'react-intl';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import { makeSelectLoggedIn } from './selectors';
import reducer from './reducer';
import saga from './saga';
import {
  SideSubscriber,
  MessageNames,
  Subscriber,
  MarketWatchSubscriber,
} from 'services/message_service';
import styled from 'styles/styled-components';
import translate from 'containers/TradePage/messages';
import { CurrencyFormater, PairFormat, zeroFixer } from 'utils/formatters';
import Skeleton from '@material-ui/lab/Skeleton';
import { Button } from '@material-ui/core';
import { Buttons, AppPages } from 'containers/App/constants';
import { push } from 'redux-first-history';
import { CurrencyPairDetails } from '../newOrder/types';
import { LocalStorageKeys } from 'services/constants';
import { FundsPages } from 'containers/FundsPage/constants';
import { currencyMap, savedPairName } from 'utils/sharedData';

let lastHeaderPrice = 0;
let priceGrown = true;
export interface CurrencyPairModel {
  price: string;
  percentage: string;
  id: number;
  high: string;
  low: string;
  volume: string;
  name: string;
  equivalent_price: string;
}

const HeaderInfo = (props: {
  title: any;
  value: string;
  valueColor?: string;
  isPercent?: boolean;
}) => {
  return props.value ? (
    <div className="infoWrapper">
      <div className="title">{props.title}</div>
      <div
        className="value"
        style={{
          color: props.valueColor ? props.valueColor : 'var(--blackText)',
        }}
      >
        {+props.value > 0 && props.isPercent ? '+' : ''}
        {props.value}
        {props.isPercent ? ' %' : ''}
      </div>
    </div>
  ) : (
    <div className="infoWrapper">
      <Skeleton
        animation="wave"
        style={{ height: '55px' }}
        className="infoWrapper"
      />
    </div>
  );
};

const stateSelector = createStructuredSelector({
  loggedIn: makeSelectLoggedIn(),
});

function TradeHeader(props: { subject: string }) {
  useInjectReducer({ key: 'tradeHeader', reducer: reducer });
  useInjectSaga({ key: 'tradeHeader', saga: saga });
  const pairName = useRef(savedPairName());

  const { loggedIn } = useSelector(stateSelector);

  const dispatch = useDispatch();

  const [ShowData, setShowData] = useState(false);
  //@ts-ignore
  const [Balance, setBalance]: [CurrencyPairDetails, any] = useState({});
  const [Data, setData]: [CurrencyPairModel, any] = useState({
    price: '',
    percentage: '',
    id: -1,
    name: '',
    high: '',
    low: '',
    volume: '',
    equivalent_price: '',
  });

  useEffect(() => {
    const BSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        setShowData(false);
        pairName.current = message.payload.name;
        setBalance({});
      }
    });
    const MarketWatchSubscription = MarketWatchSubscriber.subscribe(
      (message: any) => {
        if (message.payload.name === pairName.current && !document.hidden) {
          if (Number(message.payload.price) > lastHeaderPrice) {
            priceGrown = true;
          } else if (Number(message.payload.price) < lastHeaderPrice) {
            priceGrown = false;
          }
          setData(message.payload);
          lastHeaderPrice = message.payload.price;
          if (ShowData === false) {
            setShowData(true);
          }
        }
      },
    );
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_CURRENCY_PAIR_DETAILS) {
        setBalance(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
      BSubscription.unsubscribe();
      MarketWatchSubscription.unsubscribe();
      pairName.current = savedPairName();
    };
  }, []);
  const handleDepositClick = () => {
    localStorage[LocalStorageKeys.SELECTED_COIN] = pairName.current.split(
      '-',
    )[0];
    localStorage[LocalStorageKeys.FUND_PAGE] = FundsPages.DEPOSIT;
    dispatch(push(AppPages.Funds));
  };

  const low = zeroFixer(Data.low);
  const high = zeroFixer(Data.high);
  const price = zeroFixer(Data.price);
  const coinDigits = {
    buy: currencyMap({ code: pairName.current.split('-')[1] })?.showDigits ?? 8,
    sell:
      currencyMap({ code: pairName.current.split('-')[0] })?.showDigits ?? 8,
  };

  return (
    <Wrapper className="dragHandle">
      {ShowData === true ? (
        <div className="pairNameAndValueWrapper">
          <div className="pairName">{PairFormat(pairName.current)}</div>
          <div className="parValueWrapper">
            <div
              className={`pairValue ${
                priceGrown === true ? 'greenPrice' : 'redPrice'
              }`}
            >
              {price}
            </div>
            {/*<div className="separator">≃</div>
            <div className="equivalent">
              {CurrencyFormater(Data.equivalent_price)}
            </div>
            <div className="eqName">{' ' + 'USDT'}</div>*/}
          </div>
        </div>
      ) : (
        <div style={{ width: '140px', marginRight: '70px', marginTop: '-3px' }}>
          <Skeleton animation="wave" style={{ height: '55px' }} />
        </div>
      )}

      {high && ShowData && (
        <div className="infosWrapper">
          <HeaderInfo
            title={<FormattedMessage {...translate.Change} />}
            value={Data.percentage}
            isPercent
            valueColor={
              Number(Data.percentage) > 0
                ? 'var(--greenText)'
                : 'var(--redText)'
            }
          />
          <HeaderInfo
            title={<FormattedMessage {...translate.High} />}
            value={high}
          />
          <HeaderInfo
            title={<FormattedMessage {...translate.Low} />}
            value={low}
          />
          <HeaderInfo
            title={<FormattedMessage {...translate.T4Volume} />}
            value={CurrencyFormater(
              Data.volume ? Number(Data.volume).toFixed(2) + '' : '',
            )}
          />
        </div>
      )}
      {loggedIn === true && (
        <div className="balanceData">
          {Balance.pairBalances && (
            <div className="balanceInfo">
              <div className="av">
                <FormattedMessage {...translate.Available} />
                {' : '}
              </div>
              <div className="coin">
                {CurrencyFormater(
                  (+Balance.pairBalances[0].balance).toFixed(
                    coinDigits['buy'],
                  ) + '',
                )}{' '}
                {pairName.current.split('-')[1]}
              </div>
              <div className="equivalent">
                {CurrencyFormater(
                  (+Balance.pairBalances[1].balance).toFixed(
                    coinDigits['sell'],
                  ) + '',
                )}{' '}
                {pairName.current.split('-')[0]}
              </div>
              <div className="drop">{/*<DropDownIcon />*/}</div>
            </div>
          )}
          <Button
            onClick={handleDepositClick}
            color="primary"
            variant="contained"
            className={Buttons.DensePrimary}
          >
            <FormattedMessage {...translate.Deposit} />
          </Button>
        </div>
      )}
    </Wrapper>
  );
}

export default memo(TradeHeader);
const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  border-radius: var(--cardBorderRadius);
  background: var(--white);
  padding: 0 24px;
  display: flex;
  align-items: center;
  .pairNameAndValueWrapper {
    display: flex;
    max-width: 900px;
    min-width: 210px;
    cursor: default;
    .pairName {
      color: var(--blackText);
      font-size: 18px;
      font-weight: 600;
    }
    .pairValue {
      font-size: 18px;
      font-weight: 700;

      margin: 0 5px;
      &.greenPrice {
        color: var(--greenText);
      }
      &.redPrice {
        color: var(--redText);
      }
      /*margin-top: 0px;*/
    }
    .separator {
      margin: 0px 5px 0 0;
      color: var(--blackText);
      margin-top: 2px;
    }
    .equivalent,
    .eqName {
      font-size: 14px;
      color: var(--blackText);
      margin-top: 2px;
    }
    .eqName {
      margin: 0 5px;
      margin-top: 2px;
    }
    .parValueWrapper {
      display: flex;
    }
    @media screen and (max-width: 1520px) {
      /*flex-direction: column;*/
      min-width: 190px;
      .pairName {
        font-size: 15px;
      }
      .pairValue {
        font-size: 18px;
      }
      .equivalent,
      .eqName {
        font-size: 12px;
      }
      .parValueWrapper {
        margin-left: 0px;
        margin-top: -2px;
      }
    }
  }
  .infosWrapper {
    display: flex;
    justify-content: space-around;
  }
  .infoWrapper {
    cursor: default;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 100px;
    align-self: center;
    @media screen and (max-width: 1440px) {
      min-width: 80px;
    }

    .title {
      span {
        font-size: 13px;
        color: var(--textGrey);
      }
    }
    .value {
      color: var(--BlackText);
      font-size: 13px;
      font-weight: 600;
      @media screen and (max-width: 1440px) {
        font-size: 11px;
      }
    }
    .title,
    .value {
      width: 100%;
    }
  }
  .balanceData {
    display: flex;
    justify-content: flex-end;
    position: absolute;
    right: 20px;
  }
  .balanceInfo {
    cursor: default;
    display: flex;
    margin: 0 10px;
    align-items: center;
    margin-top: 2px;
    .av {
      font-size: 12px;
      color: var(--textGrey);
      span {
        font-size: 12px;
        color: var(--textGrey);
      }
    }
    .coin,
    .equivalent {
      font-size: 12px;
      color: var(--blackText);
      font-weight: 600;
    }
    .coin {
      margin: 0 12px 0 8px;
    }
    .drop {
      margin-top: -5px;
      margin-left: 5px;
      cursor: pointer;
    }
  }
  .densePrimary {
    max-height: 24px !important;
    padding: 0;
    min-height: 24px;
    span {
      line-height: normal !important;
      font-size: 12px !important;
      font-family: 'Open Sans' !important;
      font-weight: 600 !important;
      margin-top: -1px;
    }
  }
`;
