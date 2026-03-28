import styled from 'styles/styled-components';
const TradeHistoryDetailWrapper = styled.div`
  width: calc(100% - 24px);
  height: calc(100% - 40px);
  position: fixed;
  left: 12px;
  display: flex;
  bottom: 0px;
  overflow: hidden;
  /* pointer-events: none; */
  justify-content: center;
  align-items: center;
  &.multiLine {
    flex-direction: column;
    align-items: start;
  }

  .dataWrapper {
    display: flex;
    position: absolute;
    left: 40.9%;
  }
  .txid {
    color: var(--textGrey);
  }
  .txidValue {
    color: var(--blackText);
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    user-select: text;
  }
`;
export default TradeHistoryDetailWrapper;
