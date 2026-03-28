import React, { useState, useEffect } from 'react';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/OrdersPage/messages';
import { MenuItem, ListItemText } from '@material-ui/core';
import { Select } from '@material-ui/core';
import { Subscriber, MessageNames } from 'services/message_service';
import ExpandMore from 'images/themedIcons/expandMore';
const timePeriods = [
  {
    name: <FormattedMessage {...translate.week1} />,
    value: '1week',
  },
  {
    name: <FormattedMessage {...translate.month1} />,
    value: '1month',
  },
  {
    name: <FormattedMessage {...translate.month3} />,
    value: '3month',
  },
];
export default function SelectTimePeriod (props: {
  onPeriodSelect: Function;
  period: string;
}) {
  const { period } = props;
  const [SelectedTimePeriod, setSelectedTimePeriod] = useState('all');

  useEffect(() => {
    setSelectedTimePeriod(period.replace(/ /g, ''));
    return () => {};
  }, [period]);

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
    setSelectedTimePeriod(value);
    props.onPeriodSelect(value);
  };

  return (
    <>
      <Select
        IconComponent={ExpandMore}
        className='select'
        margin='dense'
        fullWidth
        variant='outlined'
        MenuProps={{
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        value={SelectedTimePeriod}
        onChange={(e: any) => {
          setChange(e.target.value);
        }}
      >
        <MenuItem value={'all'}>
          <FormattedMessage {...translate.timePeriod} />
        </MenuItem>
        {timePeriods.map((item: { name: any; value: string }, index) => {
          return (
            <MenuItem key={'timePeriod' + index} value={item.value}>
              <ListItemText className='addressCoin' primary={item.name} />
            </MenuItem>
          );
        })}
      </Select>
    </>
  );
}
