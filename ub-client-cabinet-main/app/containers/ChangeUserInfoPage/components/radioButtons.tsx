import React, { useState } from 'react';
import styled from 'styles/styled-components';

import {
  FormControl,
  RadioGroup,
  FormControlLabel,
  Radio,
} from '@material-ui/core';

import { FormattedMessage } from 'react-intl';
import translate from '../messages';

export default function RadioButtons (props: {
  initialValue?: string;
  onChange: Function;
}) {
  const [Value, setValue] = useState(
    props.initialValue ? props.initialValue : '',
  );
  const handleChange = e => {
    setValue(e.target.value);
    props.onChange(e.target.value);
  };
  return (
    <Wrapper>
      <FormControl component='fieldset'>
        <RadioGroup
          aria-label='position'
          name='position'
          value={Value || ''}
          onChange={handleChange}
          row
        >
          <FormControlLabel
            value='male'
            control={<Radio size='small' color='primary' />}
            label={<FormattedMessage {...translate.male} />}
            labelPlacement='end'
          />
          <FormControlLabel
            value='female'
            control={<Radio size='small' color='primary' />}
            label={<FormattedMessage {...translate.female} />}
            labelPlacement='end'
          />
        </RadioGroup>
      </FormControl>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  margin: 6px 0;
  span.MuiFormControlLabel-label {
    color: var(--textGrey);
    margin-right: 48px;
  }
`;
