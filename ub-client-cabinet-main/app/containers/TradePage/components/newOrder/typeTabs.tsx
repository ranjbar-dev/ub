import React, { useEffect } from 'react';
import styled from 'styles/styled-components';
import { Tabs, Tab } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/TradePage/messages';
import { Subscriber, MessageNames } from 'services/message_service';
import { mainType } from './types';

export default function TypeTabs (props: {
  onTabChange: (e: mainType) => void;
}) {
  const [activeIndex, setactiveIndex] = React.useState(0);
  const handleChange = (event, newactiveIndex: number) => {
    setactiveIndex(newactiveIndex);
    switch (newactiveIndex) {
      case 0:
        props.onTabChange('buy');
        break;

      case 1:
        props.onTabChange('sell');
        break;

      default:
        break;
    }
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SELECT_ORDERBOOK_ROW) {
        handleChange('', message.payload.type === 'buy' ? 0 : 1);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

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
              <span>
                <FormattedMessage {...translate.Buy} />
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
                <FormattedMessage {...translate.Sell} />
              </span>
            </>
          }
        />
      </Tabs>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  --tabWidth: 55px;
  .MuiTabs-indicator {
    /* min-width: var(--tabWidth) !important; */
    background-color: var(--tabColor) !important;
  }
  .typeTab {
    max-width: var(--tabWidth);
    min-width: var(--tabWidth);
    .MuiTab-wrapper {
      max-width: var(--tabWidth);
    }
    span {
      transition: color 0.3s;
      color: var(--textGrey) !important;
      font-weight: 500;
      font-size: 14px;
    }
    &.Mui-selected {
      span {
        color: var(--tabColor) !important;
      }
    }
  }
  border-bottom: 1px solid var(--lightGrey);
  .MuiTab-root {
    min-width: unset !important;
    max-width: fit-content;
    padding: 0 5px;
    margin-right: 12px;
    min-height: 35px;
  }
  .MuiTabs-root {
    min-height: 37px;
  }
`;
