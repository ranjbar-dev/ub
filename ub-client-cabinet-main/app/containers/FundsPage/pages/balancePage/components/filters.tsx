import React, { useState } from 'react';
import styled from 'styles/styled-components';
import IosSwitch from './iosSwitch';
import { FormattedMessage, injectIntl } from 'react-intl';
import translate from '../../../messages';

import { Subscriber, MessageNames } from 'services/message_service';
import { TextField, InputAdornment } from '@material-ui/core';
import { Translator } from 'utils/formatters';
import SearchIcon from 'images/themedIcons/searchIcon';

const Filters = (props: { minimumForSwithFilter?: string; intl: any }) => {
  const intl = props.intl;
  const [SwitchShowValue, setSwitchShowValue] = useState(true);
  const toggleShowSmallBalances = (showSmallBalances: boolean) => {
    Subscriber.next({
      name: MessageNames.SET_GRID_FILTER,
      payload: { showSmallBalances, minimum: props.minimumForSwithFilter },
    });
  };
  const onSearch = (searchCoin: string) => {
    Subscriber.next({
      name: MessageNames.SET_GRID_FILTER,
      payload: { searchCoin },
    });
  };
  return (
    <Wrapper>
      <div className='switch'>
        <span className='title'>
          <FormattedMessage {...translate.Smallbalances} />
        </span>
        <IosSwitch
          title={
            SwitchShowValue === true ? (
              <FormattedMessage {...translate.show} />
            ) : (
              <FormattedMessage {...translate.hide} />
            )
          }
          onChange={(value: boolean) => {
            setSwitchShowValue(value);
            toggleShowSmallBalances(value);
          }}
        />
      </div>
      <div className='search'>
        <TextField
          variant='outlined'
          placeholder={Translator({
            containerPrefix: 'containers.FundsPage',
            intl,
            message: 'SearchCoin',
          })}
          onChange={(e: any) => {
            onSearch(e.target.value);
          }}
          InputProps={{
            endAdornment: (
              <InputAdornment position='end'>
                <SearchIcon />
              </InputAdornment>
            ),
          }}
        />
      </div>
    </Wrapper>
  );
};
export default injectIntl(Filters);
const Wrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  /* margin-bottom: 12px; */
  height: 24px;
  margin-top: 24px;
  min-width: 1000px;
  .switch {
    display: flex;
    align-items: center;
    color: var(--textGrey);
    .title {
      margin: -1px 8px 0 0;
    }
    .MuiIconButton-root:hover {
      background-color: transparent !important;
    }
  }
  .search {
    .MuiInputBase-root {
      max-height: 24px;

      background: white;
    }
    .MuiOutlinedInput-notchedOutline {
      border-color: transparent !important;
    }
    .MuiOutlinedInput-input {
      font-size: 13px;
    }
    input {
      &::placeholder {
        font-weight: 500;
        font-size: 13px;
        font-style: italic;
        color: var(--placeHolderColor);
      }
    }
  }
`;
