import React, { memo, useState, useEffect } from 'react';

interface Props {
  className: string;
  onChange: (value: string) => void;
  initialValue: string;
  'aria-label'?: string;
  placeholder?: string;
}

/**
 * Minimal controlled plain HTML text input.
 * Syncs its value from `initialValue` via useEffect, useful for external resets.
 *
 * @example
 * ```tsx
 * <RawInput className="priceField" initialValue="0.00" onChange={(v) => setPrice(v)} />
 * ```
 */
function RawInput(props: Props) {
  const { className, onChange, initialValue } = props;
  const [Value, setValue] = useState(initialValue);
  useEffect(() => {
    setValue(initialValue);
    return () => {};
  }, [initialValue]);
  return (
    <input
      className={className}
      type="text"
      value={Value}
      aria-label={props['aria-label']}
      placeholder={props.placeholder}
      onChange={e => {
        setValue(e.target.value);
        onChange(e.target.value);
      }}
    />
  );
}

export default memo(RawInput);
