import {
  OutlinedInputProps,
  OutlinedInput,
  InputAdornment,
} from '@material-ui/core';
import React, { memo, useState } from 'react';

/**
 * Controlled text input backed by Material-UI OutlinedInput.
 * Supports an optional label prefix and end-adornment unit text.
 *
 * @example
 * ```tsx
 * <UBInput
 *   initialValue="100"
 *   label="Amount"
 *   endText="USD"
 *   onChange={(val) => setAmount(val)}
 * />
 * ```
 */
export interface UBInputProps {
  properties?: OutlinedInputProps;
  onChange: (value: string) => void;
  initialValue: string;
  label?: string;
  endText?: string;
  style?: React.CSSProperties;
  id?: string;
}

function UBInput(props: UBInputProps) {
  const [Value, setValue] = useState(props.initialValue);
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    props.onChange(e.target.value);
  };
  return (
    <>
      {!props.label && (
        <OutlinedInput
          style={{ ...(props.style && props.style) }}
          endAdornment={
            <InputAdornment position="end">
              {props.endText ?? ''}
            </InputAdornment>
          }
          value={Value}
          onChange={handleChange}
          {...props.properties}
        />
      )}
      {props.label && (
        <div className="ubInputContainer">
          <div className="label" id={props.id ? `${props.id}-label` : undefined}>
            {props.label}
            {' : '}
          </div>
          <OutlinedInput
            id={props.id}
            style={{ ...(props.style && props.style) }}
            endAdornment={
              <InputAdornment position="end">
                {props.endText ?? ''}
              </InputAdornment>
            }
            value={Value}
            onChange={handleChange}
            inputProps={{ 'aria-labelledby': props.id ? `${props.id}-label` : undefined }}
            {...props.properties}
          />
        </div>
      )}
    </>
  );
}

export default memo(UBInput);
