import { InputAdornment, TextField } from '@material-ui/core';
import React from 'react';
import styled from 'styles/styled-components';
import { formatCurrencyWithMaxFraction } from 'utils/formatters';

interface Props {
  error: string;
  value: string;
  label: string;
  placeholder: string;
  maxFraction?: number;
  onChange: (v: string) => void;
}

export const InputWithFormatter = (props: Props) => {
  const {
    onChange,
    value = '',
    error = '',
    label,
    placeholder = '',
    maxFraction = 12,
  } = props;

  const formatted = formatCurrencyWithMaxFraction(value, maxFraction);

  const handleChange = (
    e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>,
  ): void => {
    let v = e.target.value.replace(/ /g, '');
    if (v.includes('..')) {
      return;
    }
    if (v.length > 0) {
      const count = (v.match(/\./g) || []).length;
      if (count > 1) {
        return;
      }
    }
    if (v === '.') {
      v = '0.';
    }
    if (isNaN(Number(v.replace(/,/g, '')))) {
      return;
    }
    if (maxFraction && v.includes('.')) {
      const splitted = v.split('.');
      if (splitted[1] && splitted[1].length > maxFraction) {
        //prevent input if fraction length is much
        return;
      }
    }
    const noComma = v.replace(/,/g, '');
    if (Number(noComma) > 10000000000) {
      return;
    }
    onChange(v.replace(/,/g, ''));
  };
  return (
    <TextField
      variant='outlined'
      margin='dense'
      fullWidth
      value={formatted}
      error={error !== ''}
      InputProps={{
        endAdornment: (
          <InputAdornment position='end'>
            {<span className='endSpan'>{label}</span>}
          </InputAdornment>
        ),
      }}
      onChange={e => handleChange(e)}
      label={<span>{placeholder}</span>}
    />
  );
};
const Wrapper = styled.div`
  .endSpan {
    line-height: 1;
    font-size: 13px;
    font-weight: 500;
    margin-top: -1px;
    color: var(--textGrey);
  }
`;
