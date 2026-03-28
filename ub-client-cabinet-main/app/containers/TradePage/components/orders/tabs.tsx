import React from 'react';
import styled from 'styles/styled-components';
import { Tabs, Tab } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/TradePage/messages';
import { OrderPage } from 'containers/TradePage/constants';
import StreamNotificationOverLay from './streamNotificationOverLay';
import { MessageService, MessageNames } from 'services/message_service';
export default function OrderTabs (props: { onTabChange: Function }) {
  const [activeIndex, setactiveIndex] = React.useState(0);
  const handleChange = (event, newactiveIndex) => {
    setactiveIndex(newactiveIndex);
    if (newactiveIndex !== activeIndex)
      setTimeout(() => {
        switch (newactiveIndex) {
          case 0:
            props.onTabChange(OrderPage.OpenOrders);
            MessageService.send({
              name: MessageNames.HIDE_ORDER_NOTIFICATION,
              payload: OrderPage.OpenOrders,
            });
            break;
          case 1:
            props.onTabChange(OrderPage.OrderHistory);
            MessageService.send({
              name: MessageNames.HIDE_ORDER_NOTIFICATION,
              payload: OrderPage.OrderHistory,
            });
            break;
          case 2:
            props.onTabChange(OrderPage.TradeHistory);
            break;

          default:
            break;
        }
      }, 110);
  };
  return (
    <Wrapper>
      <Tabs
        value={activeIndex}
        onChange={handleChange}
        indicatorColor='primary'
        textColor='primary'
      >
        <Tab
          disableRipple={true}
          className='typeTab'
          label={
            <>
              <StreamNotificationOverLay
                pageName={OrderPage.OpenOrders}
                blinkIfStatusEquals='open'
              />
              <span>
                <FormattedMessage {...translate.OpenOrders} />
              </span>
            </>
          }
        />
        <Tab
          disableRipple={true}
          className='typeTab'
          label={
            <>
              <StreamNotificationOverLay
                pageName={OrderPage.OrderHistory}
                blinkIfStatusEquals='filled'
              />
              <span>
                <FormattedMessage {...translate.OrderHistory} />
              </span>
            </>
          }
        />
        <Tab
          disableRipple={true}
          className='typeTab'
          label={
            <>
              <span>
                <FormattedMessage {...translate.TradeHistory} />
              </span>
            </>
          }
        />
      </Tabs>
      <div className='dragHandle'></div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  --tabWidth: 100px;
  .MuiTabs-indicator {
    min-width: var(--tabWidth) !important;
    background-color: var(--textBlue) !important;
    height: 2px;
  }
  .typeTab {
    max-width: var(--tabWidth);
    min-width: var(--tabWidth);
    .MuiTab-wrapper {
      max-width: var(--tabWidth);
    }
    span {
      color: var(--blackText);
      font-weight: 600;
      font-size: 12px;
    }
  }
  .Mui-selected {
    span {
      color: var(--textBlue) !important;
    }
  }
  border-bottom: 1px solid var(--lightGrey);
  .MuiTabs-root,
  .typeTab {
    min-height: 37px;
    max-height: 37px;
    padding: 0;
  }
  .MuiTabs-indicator {
    transition: all 60ms linear 0ms;
  }
  .dragHandle {
    width: calc(100% - 325px);
    height: 37px;
    float: right;
    position: absolute;
    top: 0;
    right: 0;
  }
`;
