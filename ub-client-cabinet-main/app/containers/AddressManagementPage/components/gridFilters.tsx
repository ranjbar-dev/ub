import React, { useState } from 'react';

import { FormattedMessage } from 'react-intl';
import { Select, MenuItem, TextField, ListItemText } from '@material-ui/core';

import { Currency } from '../types';
import translate from '../messages';
import styled from 'styles/styled-components';
import { Subscriber, MessageNames } from 'services/message_service';
import { NarrowInputs } from 'global-styles';
import ExpandMore from 'images/themedIcons/expandMore';

export default function GridFilters (props: { currencyList: Currency[] }) {
  const [SelectedIndex, setSelectedIndex] = useState(-1);
  const handleSelectChange = e => {
    Subscriber.next({
      name: MessageNames.SET_GRID_FILTER,
      value:
        e.target.value != -1 ? props.currencyList[e.target.value].code : '',
      filterField: 'code',
    });
    setSelectedIndex(e.target.value);
  };
  const handleSearchChange = e => {
    Subscriber.next({
      name: MessageNames.SET_GRID_FILTER,
      value: e.target.value,
      filterField: 'label',
    });
  };
  return (
    <Wrapper className='with-ag-header'>
      <div className='select '>
        <Select
          className='select'
          margin='dense'
          IconComponent={ExpandMore}
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          fullWidth
          variant='outlined'
          value={SelectedIndex}
          onChange={handleSelectChange}
        >
          <MenuItem value={-1}>
            <FormattedMessage {...translate.all} />
          </MenuItem>
          {props.currencyList.map((item: Currency, index: number) => {
            return (
              <MenuItem key={'currencyList' + index} value={index}>
                <ListItemText className='addressCoin' primary={item.name} />
              </MenuItem>
            );
          })}
        </Select>
      </div>
      <div className='search'>
        <TextField
          margin='dense'
          onChange={handleSearchChange}
          fullWidth
          variant='outlined'
          label={<FormattedMessage {...translate.label} />}
        />
      </div>
      <div className='space'></div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  .select {
    flex: 3;
  }
  .search {
    flex: 3;
    margin: 0 10px;
    input {
      font-size: 12px;
    }
  }
  .space {
    flex: 9;
  }
  ${NarrowInputs}
`;
