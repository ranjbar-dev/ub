import React, { useEffect, useState } from 'react';

import { FormattedMessage } from 'react-intl';
import { TextField } from '@material-ui/core';
import styled from 'styles/styled-components';
import { NarrowInputs } from 'global-styles';
import translate from './messages';
import { Subscriber, MessageNames } from 'services/message_service';
export default function AddressFilter (props: { onAddressChange: Function }) {
  const [Address, setAddress] = useState('');
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_GRID_FILTER) {
        setAddress('');
        props.onAddressChange('all');
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  const handleSearchChange = e => {
    setAddress(e.target.value);
    if (e.target.value === '') {
      props.onAddressChange('all');
      return;
    }
    props.onAddressChange(e.target.value);
  };
  return (
    <Wrapper>
      <div className='search'>
        <TextField
          margin='dense'
          onChange={handleSearchChange}
          fullWidth
          value={Address}
          variant='outlined'
          label={<FormattedMessage {...translate.address} />}
        />
      </div>
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
    margin-right: 6px;
    margin-top: 2px;
    input {
      font-size: 12px;
    }
    fieldset {
      height: 36px;
    }
  }
  .space {
    flex: 9;
  }
  ${NarrowInputs}
`;
