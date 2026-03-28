import React, { useEffect, ClipboardEvent } from 'react';
import PinInput from 'react-pin-input';
import styled from 'styles/styled-components';
let v = '';
export default function UBPinInput (props: {
  onComplete: Function;
  onChange?: Function;
  onEnter?: Function;
}) {
  useEffect(() => {
    return () => {
      v = '';
    };
  }, []);
  return (
    <Wrapper
      onPaste={(e: ClipboardEvent) => {
        const f = e.clipboardData.getData('text');
        if (props.onEnter) {
          setTimeout(() => {
            //@ts-ignore
            props.onEnter(f);
            return;
          }, 200);
        }
      }}
      onKeyDown={e => {
        if (props.onEnter && v.length === 6 && e.keyCode === 13) {
          props.onEnter(v);
        }
      }}
    >
      <PinInput
        length={6}
        initialValue=''
        onChange={(value, index) => {
          v = value;

          if (props.onChange) props.onChange(value, index);
        }}
        type='numeric'
        focus
        style={{
          padding: '10px',
          display: 'flex',
          fontSize: '40px',
          justifyContent: 'space-around',
          width: '465px',
        }}
        inputStyle={{
          borderTop: 'none',
          borderRight: 'none',
          borderLeft: 'none',
          borderBottom: '2px solid #D8D8D8',
        }}
        inputFocusStyle={{ borderBottom: '2px solid var(--textBlue)' }}
        onComplete={(value, index) => {
          props.onComplete(value);
        }}
      />
    </Wrapper>
  );
}
const Wrapper = styled.div`
  align-self: center;
  .pincode-input-text {
    font-size: 40px !important;
    font-weight: 600;
    color: var(--blackText);
  }
`;
