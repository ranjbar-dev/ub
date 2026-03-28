import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';

export const SegmentSelectionWrapper = styled(Card)`
  height: 32px;
  background: var(--white);
  width: fit-content;
  margin: auto;
  margin-top: 2vh;
  border-radius: 10px !important;
  box-shadow: none !important;
  padding: 0;
  display: flex;
  align-items: center;

  .option {
    padding: 9px 24px !important;
  }
  .selected {
    color: var(--textBlue);
  }
  .Mui-selected {
    span {
      color: #396de0 !important;
      min-width: max-content;
    }
  }
  .MuiTabs-indicator {
    height: 32px;
    top: 8px;
    border-radius: 10px;
    background: rgba(57, 109, 224, 0.04);
    border: 1px solid #396de0;
    transition: all 80ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;
  }
  .MuiTab-root {
    font-weight: 600;
  }
  &.funds {
    margin-bottom: 2vh;
  }

  .MuiButtonBase-root {
    /* min-width: 130px !important; */
    span {
      padding: 0;
    }
  }

  .MuiTab-root {
    transition: margin 0.1s;
    min-width: unset !important;
    max-width: fit-content;
    padding: 0;
    margin: 0 0px;
    padding: 9px 24px;
    /* &.first {
      margin-left: 48px;
    } */
  }
  .MuiTab-textColorPrimary.Mui-selected {
    /* padding: 9px 24px;
    margin: 0 16px; */
    &.first {
      /* margin-left: 0px; */
    }
    &.last {
      /* margin-right: 0px; */
    }
  }
  .MuiTabs-root {
    min-height: 32px;
    min-width: 408px;
  }
  &.funds {
    .MuiTabs-root {
      min-height: 32px;
      min-width: 514px;
    }
  }
`;
