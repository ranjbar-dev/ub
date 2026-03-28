import { InputAdornment, TextField } from '@material-ui/core';
import React from 'react';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/TradePage/messages';
import { toFraction } from 'utils/formatters';

interface Props {
  value: string;
  onChange: (v: string) => void;
  endLabel: string;
  maxFraction?: number;
  translateKey: string;
}

function TotalInput (props: Props) {
  const { value, onChange, endLabel, maxFraction, translateKey } = props;
  const handleChange = (v: string) => {
    if (
      v &&
      maxFraction &&
      v.includes('.') &&
      v.split('.')[1] &&
      v.split('.')[1].length > maxFraction
    ) {
      return;
    }
    if (v === '.') {
      v = '0.';
    }
    onChange(v);
  };
  return (
    <TextField
      variant='outlined'
      margin='dense'
      value={maxFraction ? toFraction(value, 2) : value}
      onChange={e => handleChange(e.target.value)}
      fullWidth
      InputProps={{
        endAdornment: (
          <InputAdornment position='end'>
            {<span className='endSpan'>{endLabel}</span>}
          </InputAdornment>
        ),
      }}
      label={
        <>
          <FormattedMessage {...translate[translateKey]} />{' '}
        </>
      }
    />
  );
}

export default TotalInput;
