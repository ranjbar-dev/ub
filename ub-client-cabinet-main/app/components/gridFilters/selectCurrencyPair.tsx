import React, { useState, useEffect } from 'react';
import { Select } from '@material-ui/core';
import { MenuItem } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/OrdersPage/messages';
import { createStructuredSelector } from 'reselect';
import { makeSelectCurrencies } from 'containers/OrdersPage/selectors';
import { useSelector } from 'react-redux';
import { Currency } from 'containers/AddressManagementPage/types';
import { Subscriber, MessageNames } from 'services/message_service';
import ExpandMore from 'images/themedIcons/expandMore';
const stateSelector = createStructuredSelector({
  currencies: makeSelectCurrencies(),
});

export default function SelectCurrencyPair(props: {
  onCurrencySelect: Function;
}) {
  let pairString = 'all-all';
  const [Pair0, setPair0] = useState('all');
  const [Pair1, setPair1] = useState('all');
  const { currencies } = useSelector(stateSelector);
  const setCurrencyPair = (index: number, value: string) => {
    if (index === 0) {
      setPair0(value);
    } else {
      setPair1(value);
    }
  };
  useEffect(() => {
    pairString = Pair0 + '-' + Pair1;
    props.onCurrencySelect('pair_currency_name', pairString);
    return () => {
      // cleanup
    };
  }, [Pair1, Pair0]);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_GRID_FILTER) {
        setChange('all');
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const setChange = value => {
    setPair0(value);
    setPair1(value);
    props.onCurrencySelect('pair_currency_name', 'all-all');
  };

  return (
    <>
      <Select
        IconComponent={ExpandMore}
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        className="select pair1"
        margin="dense"
        fullWidth
        variant="outlined"
        value={Pair0}
        onChange={(e: any) => {
          setCurrencyPair(0, e.target.value);
        }}
      >
        <MenuItem value={'all'}>
          <FormattedMessage {...translate.all} />
        </MenuItem>
        {currencies.map((item: Currency, index: number) => {
          return (
            <MenuItem
              disabled={item.code === Pair1}
              key={'pair1' + index}
              value={item.code}
            >
              {item.code}
            </MenuItem>
          );
        })}
      </Select>
      <div className="divider"></div>
      <Select
        IconComponent={ExpandMore}
        className="select  pair2"
        margin="dense"
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        fullWidth
        variant="outlined"
        value={Pair1}
        onChange={(e: any) => {
          setCurrencyPair(1, e.target.value);
        }}
      >
        <MenuItem value={'all'}>
          <FormattedMessage {...translate.all} />
        </MenuItem>
        {currencies.map((item: Currency, index: number) => {
          return (
            <MenuItem
              disabled={item.code === Pair0}
              key={'pair2' + index}
              value={item.code}
            >
              {item.code}
            </MenuItem>
          );
        })}
      </Select>
    </>
  );
}
