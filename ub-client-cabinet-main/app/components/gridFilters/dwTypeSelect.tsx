import React, { useState, useEffect } from 'react';
import {
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  ListItemText,
} from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/OrdersPage/messages';
import { Subscriber, MessageNames } from 'services/message_service';
import ExpandMore from 'images/themedIcons/expandMore';
import styled from 'styles/styled-components';

export default function SelectDWType (props: { onDWTypeSelect: Function }) {
  const [DWType, setDWType] = useState('all');
  const DWTypes = [{ name: 'deposit' }, { name: 'withdraw' }];

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
    setDWType(value);
    props.onDWTypeSelect(value);
  };
  return (
    <Wrapper>
      <FormControl className='formControl'>
        <InputLabel>
          <FormattedMessage {...translate.type} />
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
          value={DWType}
          onChange={(e: any) => {
            setChange(e.target.value);
          }}
        >
          <MenuItem value={'all'}>
            <FormattedMessage {...translate.all} />
          </MenuItem>
          {DWTypes.map((item, index) => {
            return (
              <MenuItem key={'selectDWTypeInFilter' + index} value={item.name}>
                <ListItemText className='addressCoin' primary={item.name} />
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
