import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';
import { TextField, InputAdornment } from '@material-ui/core';
import errorIcon from 'images/errorIcon.svg';

import { Subscriber, MessageNames } from 'services/message_service';
import EyeIcon from 'images/themedIcons/eyeIcon';
import CrossedEyeIcon from 'images/themedIcons/crossedEyeIcon';
let timeOut;
export default function InputWithValidator (props: {
  label: any;
  uniqueName: string;
  onChange: Function;
  initialValue?: string;
  rows?: number;
  throttleTime?: number;
  inputType?: string;
  onEnter?: Function;
  isPickable?: boolean;
  className?: string;
  startComponent?: any;
  autoComplete?: string;
  autoFocus?: boolean;
}) {
  const [Value, setValue] = useState(
    props.initialValue ? props.initialValue : '',
  );
  const [CanPick, setCanPick] = useState(false);
  const [ValidationError, setValidationError] = useState(null);
  const [IsLong, setIsLong] = useState(false);
  const [TxtAreaHeight, setTxtAreaHeight] = useState(40);

  const resetTextAreaHeight = () => {
    if (props.rows) {
      const inp: any = document.querySelector('#' + props.uniqueName);
      if (inp) {
        const height = inp.scrollHeight;
        if (inp.value.length === 0 || inp.value.length < 10) {
          inp.style.minHeight = 40 + 'px';
          return;
        }
        if (height > 40 && height < 150) {
          inp.style.minHeight = height + 'px';
        } else if (height == 40) {
          inp.style.minHeight = 40 + 'px';
        }
      }
    }
  };

  const handleChange = (e: any) => {
    e.persist();
    clearTimeout(timeOut);
    setValue(e.target.value);
    timeOut = setTimeout(
      () => {
        props.onChange(e.target.value);
      },
      props.throttleTime ? props.throttleTime : 300,
    );
    resetTextAreaHeight();
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_INPUT_ERROR) {
        if (message.value === props.uniqueName) {
          if (message.additional) {
            setIsLong(true);
          }
          setValidationError(message.payload);
        }
      }
    });
    setTimeout(() => {
      resetTextAreaHeight();
    }, 0);
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const handleKeyDown = (e: any) => {
    if (props.onEnter) {
      if (e.keyCode === 13) {
        props.onEnter();
      }
    }
  };
  let type = 'text';
  if (props.inputType) {
    type = props.inputType;
  }
  if (CanPick == true) {
    type = 'text';
  } else {
    type = props.inputType ? props.inputType : 'text';
  }
  return (
    <div style={{ width: '100%' }}>
      <Wrapper
        style={{ width: '100%' }}
        className={`inputWithValidator ${
          props.rows ? 'textArea' + TxtAreaHeight : ''
        } ${props.className ?? ''}`}
      >
        <TextField
          rows={props.rows ? props.rows : 1}
          multiline={props.rows ? true : false}
          fullWidth
          variant='outlined'
          label={props.label}
          margin='dense'
          id={props.uniqueName}
          autoFocus={props.autoFocus}
          value={Value}
          type={type}
          autoComplete={props.autoComplete}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          error={ValidationError ? true : false}
          InputProps={{
            autoComplete: type === 'password' ? 'new-password' : 'type',
            endAdornment: (
              <InputAdornment position='end'>
                {props.isPickable === true && (
                  <div
                    onClick={() => {
                      setCanPick(!CanPick);
                    }}
                    className='eyeButton'
                    //size="small"
                  >
                    {CanPick === false ? (
                      <EyeIcon color={'var(--textGrey)'} />
                    ) : (
                      <CrossedEyeIcon color={'var(--textGrey)'} />
                    )}
                  </div>
                )}
              </InputAdornment>
            ),
            startAdornment: props.startComponent && (
              <InputAdornment position='start'>
                {props.startComponent}
              </InputAdornment>
            ),
          }}
        />
        {ValidationError && (
          <div
            style={{ minWidth: IsLong === true ? '590px' : '' }}
            className='errorWrapper'
          >
            <span className='errorIcon'>
              <img src={errorIcon} alt='' />
            </span>
            <span className='errorText'>{ValidationError}</span>
          </div>
        )}
      </Wrapper>
    </div>
  );
}
const Wrapper = styled.div`
  margin: 10px 0 0 0;
  position: relative;
  &.textArea40 {
    textarea {
      transition: min-height 0.3s;
      min-height: 40px;
    }
  }

  .errorWrapper {
    position: absolute;
    bottom: -15px;
    left: 0px;
    color: var(--redText);
    font-size: 11px;
    min-width: 360px;
    display: flex;
    span {
      font-size: 11px;
    }
    img {
      width: 20px;
    }
  }
  .preNumber {
    background: var(--oddRows);
    padding: 3px 5px;
    border-radius: 7px;
    font-size: 13px;
    color: var(--blackText);
    margin-top: 2px;
    font-weight: 600;
  }
  .eyeButton {
    cursor: pointer;
  }
`;
