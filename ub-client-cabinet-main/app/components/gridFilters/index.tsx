import React, { useRef, useState } from 'react';
import styled from 'styles/styled-components';
import translate from 'containers/OrdersPage/messages';
import { FormattedMessage } from 'react-intl';
import { NarrowInputs } from 'global-styles';
import OrderHistoryDatePickers from './datePicker';
import { Button } from '@material-ui/core';
import { Checkbox, FormControlLabel } from '@material-ui/core';

import SelectTimePeriod from './selectTimePeriod';
import SelectCurrencyPair from './selectCurrencyPair';
import SelectType from './selectType';
import { Subscriber, MessageNames } from 'services/message_service';
import { Buttons } from 'containers/App/constants';
import StreamLoadingButton from 'components/streamLoadingButton';
import SelectCoin from './coinSelect';
import AddressFilter from './addressFilter';
import SelectDWType from './dwTypeSelect';
import CheckIcon from 'images/themedIcons/checkIcon';
import CheckedIcon from 'images/themedIcons/checkedIcon';
type filterKeys =
  | 'period'
  | 'address'
  | 'type'
  | 'code'
  | 'start_date'
  | 'buy'
  | 'sell'
  | 'end_date'
  | 'hideCancelledOrders';
const GridFilters = (props: {
  onSearchClick: Function;
  onCancelClick: Function;
  hideTick?: boolean;
  TimePeriod?: boolean;
  CurrencyPair?: boolean;
  BuySell?: boolean;
  Coin?: boolean;
  DWType?: boolean;
  Address?: boolean;
}) => {
  const [startDate, setStartDate] = useState('Start Time');
  const [endDate, setEndDate] = useState('End Time');
  const [period, setPeriod] = useState('all');
  const dataToSend = useRef<any>({});
  const setSendingData = (key: filterKeys, value: string) => {
    if (dataToSend.current[key] === value) {
      return;
    }
    if (key === 'period') {
      const now = new Date();
      const end = now.toISOString().split('T')[0];
      let from: any;
      const v: '1week' | '3month' | '1month' | 'all' = value as
        | '1week'
        | '3month'
        | '1month'
        | 'all';
      const date = new Date();
      switch (v) {
        case '1week':
          from = new Date(date.setDate(date.getDate() - 7))
            .toISOString()
            .split('T')[0];
          setStartDate(() => from);
          setEndDate(() => end);
          break;
        case '1month':
          from = new Date(date.setDate(date.getDate() - 30))
            .toISOString()
            .split('T')[0];
          setStartDate(() => from);
          setEndDate(() => end);
          break;
        case '3month':
          from = new Date(date.setDate(date.getDate() - 90))
            .toISOString()
            .split('T')[0];
          setStartDate(() => from);
          setEndDate(() => end);
          break;
        case 'all':
          setStartDate('Start Time');
          setEndDate('Start Time');
      }
    }
    if (key === 'start_date' || key === 'end_date') {
      dataToSend.current['period'] = 'all';
      setPeriod('all ');
      if (key === 'start_date') {
        setStartDate(value ? value : 'Start Time');
      } else {
        setEndDate(value ? value : 'End Time');
      }
    }
    dataToSend.current[key] = value;
    Subscriber.next({
      name: MessageNames.SET_GRID_FILTER,
      payload: dataToSend.current,
    });
  };
  const onSearchClick = () => {
    const data = dataToSend.current;

    for (const key in data) {
      if (data[key] == '' || data[key] == 'all') {
        delete data[key];
      }

      if (key === 'start_date') {
        if (endDate && endDate.includes('-')) {
          data.end_date = endDate;
        }
      }
      if (key === 'end_date') {
        if (startDate && startDate.includes('-')) {
          data.start_date = startDate;
        }
      }
    }
    if (data.start_date) {
      data.start_date = data.start_date.split(' ')[0] + ' 00:00:00';
    }
    if (data.end_date) {
      data.end_date = data.end_date.split(' ')[0] + ' 23:59:59';
    }
    if (data.period) {
      delete data.start_date;
      delete data.end_date;
    }

    props.onSearchClick(data);
  };
  const onfilterCancelClick = () => {
    Subscriber.next({
      name: MessageNames.RESET_GRID_FILTER,
    });
    props.onCancelClick();
  };
  return (
    <div style={{ width: '100%' }}>
      <FilterWrapper className='filterWrapper'>
        {props.Coin == true && (
          <div className='select2'>
            <SelectCoin
              onCoinSelect={(e: string) => {
                setSendingData('code', e);
              }}
            />
          </div>
        )}
        {props.DWType == true && (
          <div className='select' style={{ margin: '0 5px' }}>
            <SelectDWType
              onDWTypeSelect={(e: string) => {
                setSendingData('type', e);
              }}
            />
          </div>
        )}
        {props.Address == true && (
          <div className='address'>
            <AddressFilter
              onAddressChange={(e: string) => {
                setSendingData('address', e);
              }}
            />
          </div>
        )}
        {props.TimePeriod == null && (
          <div className='select'>
            <SelectTimePeriod
              period={period}
              onPeriodSelect={(e: string) => {
                setSendingData('period', e);
              }}
            />
          </div>
        )}
        <div className='datePickers'>
          <OrderHistoryDatePickers
            {...{ startDate, endDate }}
            onDateSelect={(key: string, value: string) => {
              //@ts-ignore
              setSendingData(key, value);
            }}
          />
        </div>
        {props.CurrencyPair == null && (
          <div className='currencyPair'>
            <SelectCurrencyPair
              onCurrencySelect={(key: string, value: string) => {
                //@ts-ignore
                setSendingData(key, value);
              }}
            />
          </div>
        )}
        {props.BuySell == null && (
          <div className='select2'>
            <SelectType
              onTypeSelect={(key, value) => {
                setSendingData(key, value);
              }}
            />
          </div>
        )}
        <div className='buttons'>
          <StreamLoadingButton
            onClick={onSearchClick}
            className='button'
            color='primary'
            variant='contained'
            text={<FormattedMessage {...translate.search} />}
          />
          <Button
            onClick={onfilterCancelClick}
            className={`button cancel ${Buttons.CancelButton}`}
          >
            <FormattedMessage {...translate.Reset} />
          </Button>
        </div>
        <div className='spacer'>
          {!props.hideTick && (
            <FormControlLabel
              className='checkIcon'
              control={
                <Checkbox
                  onChange={(e: any) => {
                    setSendingData('hideCancelledOrders', e.target.checked);
                  }}
                  checkedIcon={<CheckedIcon />}
                  icon={<CheckIcon />}
                  value='checkBox'
                  color='primary'
                />
              }
              label={<FormattedMessage {...translate.Hidecancelledorders} />}
            />
          )}
        </div>
      </FilterWrapper>
    </div>
  );
};
export default GridFilters;
const FilterWrapper = styled.div`
  display: flex;
  margin-bottom: 12px;
  margin-top: 10px;
  /* min-width: 1470px; */
  ${NarrowInputs}
  .MuiInputBase-root.select {
    margin-bottom: 0px;
  }
  .select {
    flex: 1;
    color: var(--textGrey);
  }
  .datePickers {
    flex: 2;
    display: flex;
    position: relative;
    min-width: 285px;
    justify-content: center;
  }
  .currencyPair {
    flex: 2;
    display: flex;
    margin-right: 16px;
    margin-left: 12px;
    @media screen and (max-width: 1600px) {
      margin-left: 20px;
      margin-right: 8px;
    }
    position: relative;
    .pair2 {
      min-width: 50px !important;
      max-width: 152px !important;
      border-top-left-radius: 0px !important;
      border-bottom-left-radius: 0px !important;
    }
    .pair1 {
      min-width: 50px !important;
      max-width: 152px !important;
      border-top-right-radius: 0px !important;
      border-bottom-right-radius: 0px !important;
    }
    .divider {
      width: 12px;
      background: var(--darkGrey);
      height: 30px;
      position: absolute;
      margin-left: 48%;
      margin-top: 9px;
      border: 5px solid var(--white);
      border-top: 6px solid var(--white);
      z-index: 2;
    }
  }
  .select2 {
    flex: 1;
  }
  .address {
    flex: 2;
    min-width: 190px;
  }
  .buttons {
    flex: 2;
    display: flex;
    align-items: center;
    justify-content: flex-start;
    .button {
      height: 32px;
      width: 94px;
      padding: 0 !important;
      margin-top: 2px;
      margin: 4px 0 0 24px;
      @media screen and (max-width: 1600px) {
        margin-left: 8px;
      }
      &.cancel {
        color: var(--textGrey);
      }
    }
  }
  .spacer {
    flex: 2;
    display: flex;
    place-content: flex-end;
    color: var(--textGrey);
    min-width: 190px;
    .Mui-checked {
      path {
        fill: var(--textBlue) !important;
      }
    }
  }
  .horizontalSpacer {
    width: 0.5vw;
  }
  .MuiButton-outlined {
    border: 1px solid var(--inputBorderColor);
    padding: 5px 15px;
  }
  .MuiOutlinedInput-root.Mui-error .MuiOutlinedInput-notchedOutline {
    border-color: var(--inputBorderColor) !important;
  }
  /* fieldset {
    border: 1px solid rgba(0, 0, 0, 0.23) !important;
  } */
  .loadingCircle {
    top: 7px !important;
  }
  .MuiPopover-paper {
    margin-top: 41px;
  }
  .formControl label {
    left: 10px !important;
    color: var(--textGrey);
  }
`;
