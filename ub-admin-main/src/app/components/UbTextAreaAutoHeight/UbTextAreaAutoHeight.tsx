import { TextareaAutosize } from '@material-ui/core';
import React, { memo, useState } from 'react';

interface Props {
  onChange: (value: string) => void;
  initialValue: string;
  placeHolder?: string;
  style?: React.CSSProperties;
}

/**
 * Auto-growing textarea backed by Material-UI TextareaAutosize.
 *
 * @example
 * ```tsx
 * <UbTextAreaAutoHeight
 *   initialValue="Notes here..."
 *   placeHolder="Enter notes"
 *   onChange={(val) => setNotes(val)}
 * />
 * ```
 */
function UbTextAreaAutoHeight(props: Props) {
  const { onChange, initialValue, placeHolder, style } = props;
  const [InputValue, setInputValue] = useState(initialValue);
  return (
    <TextareaAutosize
      value={InputValue}
      placeholder={placeHolder ?? ''}
      onChange={e => {
        onChange(e.target.value);
        setInputValue(e.target.value);
      }}
      rows={8}
      style={{
        width: '100%',
        fontFamily: 'Open Sans',
        padding: '12px',
        ...style,
      }}
    />
  );
}

export default memo(UbTextAreaAutoHeight);
