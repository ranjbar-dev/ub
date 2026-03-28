import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';
import { Subscriber, MessageNames } from 'services/message_service';
import { LocalStorageKeys } from 'services/constants';
import { OrderPage } from 'containers/TradePage/constants';

export default function StreamNotificationOverLay (props: {
  pageName: OrderPage;
  blinkIfStatusEquals: string;
}) {
  const [Active, setActive] = useState(false);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.NEW_ORDER_NOTIFICATION &&
        localStorage[LocalStorageKeys.VISIBLE_ORDER_SECTION] !==
          props.pageName &&
        message.payload &&
        message.payload.status === props.blinkIfStatusEquals
      ) {
        setActive(true);
      }
      if (
        message.name === MessageNames.HIDE_ORDER_NOTIFICATION &&
        message.payload == props.pageName
      ) {
        setActive(false);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [Active]);
  return <Wrapper className={Active === true ? 'active' : ''}></Wrapper>;
}
const Wrapper = styled.div`
  position: absolute;
  width: 5px;
  height: 5px;
  background: var(--orange);
  z-index: 1;
  transition: opacity 1s;
  opacity: 0;
  right: 3px;
  border-radius: 10px;
  top: 10px;
  &.active {
    opacity: 1;
  }
  /*@keyframes blink-animation {
    to {
      opacity: 1;
    }
  }
  @-webkit-keyframes blink-animation {
    to {
      opacity: 1;
    }
  }*/
`;
