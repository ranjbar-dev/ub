import React, { useState } from 'react';
import styled from 'styles/styled-components';

import { Button } from '@material-ui/core';
import calendarIcon from 'images/calendarIcon.svg';
import 'ub-date-picker/vaadin-date-picker.js';

let startDatePicker, endDatePicker: any;

export default function BirthDaySelector (props: {
  initialValue: string;
  onDateSelect: Function;
}) {
  const [StartDate, setStartDate] = useState(props.initialValue);

  customElements.whenDefined('vaadin-date-picker').then(function () {
    startDatePicker = document.querySelector('#start');
    if (startDatePicker) {
      startDatePicker.addEventListener('change', function (event: any) {
        setStartDate(event.target.value);
        props.onDateSelect(event.target.value);
      });
    }
  });

  return (
    <Wrapper>
      <HiddenWrapper>
        <vaadin-date-picker
          id='start'
          placeholder='Start Time'
        ></vaadin-date-picker>
      </HiddenWrapper>
      <div className='inputs'>
        <Button
          onClick={() => {
            startDatePicker.click();
          }}
          fullWidth
          className='dateInput'
          variant='outlined'
          endIcon={<img src={calendarIcon} />}
        >
          <span className='date'>{StartDate}</span>
        </Button>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  margin-top: 6px;
  .inputs {
    .dateInput {
      min-height: 40px;
      max-height: 40px;
      padding: 0;
      span.date {
        color: var(--textGrey);
        font-weight: 600;
      }
    }
    .dateInput {
      border: 1px solid var(--inputBorderColor);
      background: transparent !important;
      &:hover {
        border: 1px solid var(--textBlue);
      }
      img {
        position: absolute;
        right: 12px;
        top: 10px;
      }
      span.date {
        position: absolute;
        left: 12px;
        font-size: 13px;
        font-weight: 600;
        color: var(--textGrey);
      }
    }
  }
  .MuiButton-label {
    justify-content: space-between;
    padding: 0 12px;
  }
`;
const HiddenWrapper = styled.div`
  opacity: 0;
  pointer-events: none;
  height: 0;
`;
