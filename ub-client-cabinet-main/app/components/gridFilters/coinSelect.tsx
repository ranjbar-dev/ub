import React, { useState, useEffect } from 'react';
import { MenuItem, Select, FormControl, InputLabel } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/OrdersPage/messages';
import { Subscriber, MessageNames } from 'services/message_service';
import ExpandMore from 'images/themedIcons/expandMore';
import { LocalStorageKeys } from 'services/constants';
import { Currency } from 'containers/App/types';
import styled from 'styles/styled-components';
export default function SelectCoin (props: { onCoinSelect: Function }) {
  const [Coin, setCoin] = useState('all');
  const coins: Currency[] = JSON.parse(
    localStorage[LocalStorageKeys.CURRENCIES],
  );

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
    setCoin(value);
    props.onCoinSelect(value);
  };
  return (
    <Wrapper>
      <FormControl className='formControl'>
        <InputLabel>
          <FormattedMessage {...translate.coin} />
        </InputLabel>
        <Select
          IconComponent={ExpandMore}
          className='select'
          margin='dense'
          fullWidth
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          variant='outlined'
          value={Coin}
          onChange={(e: any) => {
            setChange(e.target.value);
          }}
        >
          <MenuItem value={'all'}>
            <FormattedMessage {...translate.all} />
          </MenuItem>
          {coins.map((item, index) => {
            return (
              <MenuItem key={'selectCoinInFilter' + index} value={item.code}>
                {item.code}
              </MenuItem>
            );
          })}
        </Select>
      </FormControl>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  .formControl {
    width: 100%;
    legend {
      width: 30px !important;
    }
    label {
      left: 10px !important;
    }
  }
`;
