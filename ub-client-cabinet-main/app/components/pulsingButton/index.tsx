import React, { useEffect, useState } from 'react';
import styled from 'styles/styled-components';

import { Button } from '@material-ui/core';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { Subscriber, MessageNames } from 'services/message_service';
function PulsingButton (props: { title: any; onClick: Function }) {
  const [IsLoading, setIsLoading] = useState(false);
  useEffect(() => {
    const subscription = Subscriber.subscribe((message: any) => {
      if (message.name == MessageNames.SETLOADING) {
        if (message.element == 'pulsingButton') {
          setIsLoading(message.payload);
        }
      }
    });
    return () => {
      subscription.unsubscribe();
    };
  }, []);

  return (
    <Wrapper>
      <Button
        onClick={() => props.onClick()}
        color='primary'
        variant='contained'
        className='button'
      >
        <IsLoadingWithText
          isLoading={IsLoading === true ? IsLoading : false}
          text={props.title}
        />
      </Button>
      <div className='blob'>
        <div className='ring'></div>
        <div className='pulse pulse--1'></div>
        <div className='pulse pulse--2'></div>
        <div className='pulse pulse--3'></div>
        <div className='pulse pulse--4'></div>
      </div>
    </Wrapper>
  );
}

export default PulsingButton;
const $blob_diameter = '0px';
const $ring_diameter = '60px';
const $animation_duration = '2s';
const $x_scale = '1.2';
const $y_scale = '1.5';
const Wrapper = styled.div`
  .button {
    margin-top: 8px;
    min-width: 92px;
    z-index: 1;
    padding: 4px 24px !important;
  }
  .blob {
    width: 0px;
    height: 0px;
    background: transparent;
    border-radius: 7px;
    float: left;
    margin-top: 40px;
  }

  .ring,
  .pulse {
    width: 92px;
    height: 32px;
    border-radius: 7px;
    margin-top: -32px;
    border: solid #396de0;
  }

  .ring {
    border-width: 2px;
  }

  .pulse {
    transform: scale(1, 1);
    border-width: 1px;
    animation: PULSE infinite;
    animation-duration: ${$animation_duration};
    &.pulse--1 {
      animation-delay: 0.1s;
    }
    &.pulse--2 {
      animation-delay: 0.3s;
    }
    &.pulse--3 {
      animation-delay: 0.5s;
    }
    &.pulse--4 {
      animation-delay: 0.8s;
    }
  }

  @keyframes PULSE {
    100% {
      transform: scale(${$x_scale}, ${$y_scale});
      opacity: 0;
    }
  }
  span.loadingCircle {
    top: 6px !important;
  }
`;
