import React, { useEffect } from 'react';
import { Select, MenuItem } from '@material-ui/core';

import LayOutIcon from 'images/themedIcons/layoutIcon';
import styled from 'styles/styled-components';
import ChartLayOutIcon from 'images/themedIcons/chartLayoutIcon';
import {
  MessageService,
  MessageNames,
  Subscriber,
} from 'services/message_service';
import { LocalStorageKeys } from 'services/constants';
import { LayoutNames } from './constants';

const layouts = [
  { icon: <LayOutIcon />, name: 'Default Deck', id: LayoutNames.DEFAULT },
  { icon: <ChartLayOutIcon />, name: 'Chart Deck', id: LayoutNames.CHART },
  { icon: <LayOutIcon />, name: 'Custom Deck', id: LayoutNames.CUSTOM },
];

export default function LayOutSelect () {
  const [Layout, setLayout] = React.useState(
    localStorage[LocalStorageKeys.LAYOUT_NAME] === LayoutNames.DEFAULT
      ? 0
      : localStorage[LocalStorageKeys.LAYOUT_NAME] === LayoutNames.CHART
      ? 1
      : localStorage[LocalStorageKeys.LAYOUT_NAME] === LayoutNames.CUSTOM
      ? 2
      : 0,
  );
  if (!localStorage[LocalStorageKeys.LAYOUT_NAME]) {
    localStorage[LocalStorageKeys.LAYOUT_NAME] = LayoutNames.DEFAULT;
  }
  const handleChange = event => {
    setLayout(event.target.value);
    MessageService.send({
      name: MessageNames.CHANGE_LAYOUT,
      payload: layouts[event.target.value].id,
    });
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.LAYOUT_CHANGE) {
        setLayout(2);
        MessageService.send({
          name: MessageNames.CHANGE_LAYOUT,
          payload: LayoutNames.CUSTOM,
        });
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Wrapper>
      <Select
        MenuProps={{
          style: { marginLeft: '-33px' },
          getContentAnchorEl: null,
          anchorOrigin: {
            vertical: 'bottom',
            horizontal: 'left',
          },
        }}
        className='layoutToggleSelect'
        value={Layout}
        onChange={handleChange}
      >
        {layouts.map((item, index) => {
          return (
            <MenuItem
              className='layoutItem'
              key={'layout' + item.name}
              value={index}
            >
              {item.icon}
              <span className='text'>{item.name}</span>
            </MenuItem>
          );
        })}
      </Select>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  margin: 8px 0px;

  .MuiInput-input {
    .text {
      display: none;
    }
  }
  .MuiInput-underline:after {
    display: none;
  }
  .MuiInput-underline:before {
    display: none;
  }
`;
