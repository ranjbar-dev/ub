import React, { useState, memo, useEffect } from 'react';
// import 'components/Customized/vaadin/vaadin-date-picker/vaadin-date-picker.js';
import 'ub-date-picker/vaadin-date-picker.js';
import styled from 'styles/styled-components';
import { Button } from '@material-ui/core';

import calendarIcon from 'images/calendarIcon.svg';
import { Subscriber, MessageNames } from 'services/message_service';

let startDatePicker, endDatePicker: any;
const OrderHistoryDatePickers = (props: {
  startDate: string;
  endDate: string;
  onDateSelect: (k: 'start_date' | 'end_date', v: string) => void;
}) => {
  const { startDate, endDate } = props;
  const [StartDate, setStartDate] = useState('Start Time');
  const [EndDate, setEndDate] = useState('End Time');

  useEffect(() => {
    setStartDate(startDate);
    return () => {};
  }, [startDate]);

  useEffect(() => {
    setEndDate(endDate);
    return () => {};
  }, [endDate]);

  customElements.whenDefined('vaadin-date-picker').then(function () {
    startDatePicker = document.querySelector('#start');
    endDatePicker = document.querySelector('#end');

    if (startDatePicker) {
      startDatePicker.removeEventListener('change', null);
      startDatePicker.addEventListener('change', function (event: any) {
        setStartDate(event.target.value);
        props.onDateSelect('start_date', event.target.value);
      });
    }
    if (endDatePicker) {
      endDatePicker.removeEventListener('change', null);
      endDatePicker.addEventListener('change', function (event: any) {
        setEndDate(event.target.value);
        props.onDateSelect('end_date', event.target.value);
      });
    }
  });
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_GRID_FILTER) {
        resetChange();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const resetChange = () => {
    setStartDate('Start Time');
    setEndDate('End Time');
    props.onDateSelect('start_date', '');
    props.onDateSelect('end_date', '');
  };
  return (
    <Wrapper>
      <HiddenWrapper>
        <vaadin-date-picker
          id='start'
          placeholder='Start Time'
        ></vaadin-date-picker>
        <div className='abs'>
          <vaadin-date-picker
            id='end'
            placeholder='End Time'
          ></vaadin-date-picker>
        </div>
      </HiddenWrapper>
      <div className='inputs'>
        <Button
          // value={EndDate}
          onClick={() => {
            startDatePicker.click();
          }}
          className='dateInput'
          variant='outlined'
        >
          <span className='date'>{StartDate}</span>
          <img src={calendarIcon} />
        </Button>
        <Button
          // value={EndDate}
          onClick={() => {
            endDatePicker.click();
          }}
          className='dateInput'
          variant='outlined'
        >
          <span className='date'>{EndDate}</span>
          <img src={calendarIcon} />
        </Button>
      </div>
    </Wrapper>
  );
};

export default memo(OrderHistoryDatePickers);

const Wrapper = styled.div`
  display: flex;
  min-width: 275px;
  .inputs {
    position: absolute;
    margin-top: 8px;
    display: flex;
    height: 100%;
    place-items: center;
    margin-top: 2px;
    justify-content: center;
    .dateInput {
      width: 129px;
      margin: 0 6px;
      min-height: 32px;
      max-height: 32px;
      padding: 0;
      border: 1px solid var(--inputBorderColor);
      background: transparent !important;
      &:hover {
        border: 1px solid var(--textBlue);
      }
      img {
        position: absolute;
        right: 12px;
      }
      span.date {
        position: absolute;
        left: 12px;
        font-weight: 600;
        font-size: 13px;
        color: var(--textGrey);
      }
    }
  }
`;
const HiddenWrapper = styled.div`
  opacity: 0;
  pointer-events: none;
  .abs {
    position: absolute;
    top: 0;
    left: 130px;
  }
`;
