import React, { useState, useEffect } from 'react';
import { MenuItem, Select, ListItemText } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/OrdersPage/messages';
import { Subscriber, MessageNames } from 'services/message_service';
import ExpandMore from 'images/themedIcons/expandMore';
export default function SelectType (props: {
  onTypeSelect: (key: 'type', v: 'buy' | 'sell' | 'all') => void;
}) {
  const [Type, setType] = useState('all');
  const types = [
    { name: <FormattedMessage {...translate.buy} />, value: 'buy' },
    { name: <FormattedMessage {...translate.sell} />, value: 'sell' },
  ];

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
    setType(value);
    props.onTypeSelect('type', value);
  };
  return (
    <>
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
        value={Type}
        onChange={(e: any) => {
          setChange(e.target.value);
        }}
      >
        <MenuItem value={'all'}>
          <FormattedMessage {...translate.all} />
        </MenuItem>
        {types.map((item, index) => {
          return (
            <MenuItem key={'selectTypeInFilter' + index} value={item.value}>
              <ListItemText className='addressCoin' primary={item.name} />
            </MenuItem>
          );
        })}
      </Select>
    </>
  );
}
