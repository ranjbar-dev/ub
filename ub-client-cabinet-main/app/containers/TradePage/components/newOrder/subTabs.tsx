import React from 'react';
import styled from 'styles/styled-components';
import { Tabs, Tab } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import translate from 'containers/TradePage/messages';
import { subType } from './types';
export default function SubTabs (props: { onTabChange: (e: subType) => void }) {
  const [activeIndex, setactiveIndex] = React.useState(0);
  const handleChange = (event, newactiveIndex) => {
    setactiveIndex(newactiveIndex);
    switch (newactiveIndex) {
      case 0:
        props.onTabChange('limit');
        break;
      case 1:
        props.onTabChange('market');
        break;
      case 2:
        props.onTabChange('stop_limit');
        break;

      default:
        break;
    }
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
              <span>
                <FormattedMessage {...translate.Limit} />
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
                <FormattedMessage {...translate.Market} />
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
                <FormattedMessage {...translate.StopLimit} />
              </span>
            </>
          }
        />
      </Tabs>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  --tabWidth: 95px;
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
    margin-top: 8px;
  }
`;
