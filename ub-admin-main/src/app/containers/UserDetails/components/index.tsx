import { Button } from '@material-ui/core';
import React, { memo, useState, useCallback } from 'react';
import styled from 'styled-components/macro';

interface Props {
  options: { title: string; component: React.ReactNode; icon?: React.ReactNode }[];
}

function UserDetailsTabs(props: Props) {
  const [SelectedIndex, setSelectedIndex] = useState(0);

  const handleSegmentClick = useCallback(index => {
    setSelectedIndex(index);
  }, []);

  return (
    <Wrapper>
      <div className="subTabButtons">
        {props.options.map((item, index) => {
          return (
            <Button
              disableRipple
              startIcon={item.icon ?? null}
              className={`tabButton ${
                index === SelectedIndex ? 'active' : 'inactive'
              }`}
              key={'segment' + index}
              onClick={() => handleSegmentClick(index)}
            >
              {item.title}
            </Button>
          );
        })}
      </div>
      <div className="pageContainer NWindow">
        {props.options[SelectedIndex].component}
      </div>
    </Wrapper>
  );
}

export default memo(UserDetailsTabs);
const Wrapper = styled.div`
  .tabButton {
    min-width: 70px !important;
    margin: 5px 10px !important;

    color: ${p => p.theme.blackText} !important;
  }
  .active {
    background: ${p => p.theme.lightBlue} !important;
    color: #3e5ca1 !important;
    background: transparent !important;
    font-weight: 600;
  }
  .inactive {
    background: transparent !important;
    svg {
      fill: #888888 !important;
    }
    color: rgb(95 95 95) !important;
  }
  .pageContainer {
    display: flex;
    justify-content: start;
    height: calc(100% - 140px);
    padding-left: 24px;
    &.NWindow {
      padding-left: 0px;
    }
  }
  .subTabButtons {
    background: #e4eef6;
    margin-bottom: 33px;
    border-bottom: 1px solid #dfe1e3;
  }
  .MuiButton-label {
    font-size: 12px !important;
  }
`;
