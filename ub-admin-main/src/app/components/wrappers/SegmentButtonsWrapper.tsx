import styled from 'styled-components/macro';

export const SegmentSeletionButtonsWrapper = styled.div`
  height: 40px;
  width: 100%;
  margin: auto;
  margin: 2vh auto;
  border-radius: 10px !important;
  box-shadow: none !important;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  &.orders {
    margin-bottom: 0;
  }
  .segmentButton {
    background: transparent !important;
    padding: 4px 12px !important;
    margin: 0 6px;
    span {
      color: ${p => p.theme.textGrey};
      line-height: initial;
      font-size: 13px;
      font-weight: 600 !important;
      font-family: 'Open Sans' !important;
    }
    border: 1px solid ${p => p.theme.lightGrey};
    &.selected {
      background: ${p => p.theme.white} !important;
      border: 1px solid ${p => p.theme.white};
      span {
        color: ${p => p.theme.primary};
      }
    }
  }
`;
