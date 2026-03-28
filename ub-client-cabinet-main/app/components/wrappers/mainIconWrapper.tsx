import styled from 'styles/styled-components';

const MainIconWrapper = styled.div`
  padding: 10vh 0 1vh 0;
  flex: 5;

  &.pt6vh {
    padding-top: 6vh;
  }
  &.noPadding {
    padding: 0;
  }
  &.fl {
    display: flex;
    align-items: flex-end;
    padding: 0;
  }
  &.fl {
    display: flex;
    align-items: flex-end;
    padding: 0;
  }
  &.fl12 {
    flex: 12;
  }
  &.fl9 {
    flex: 9;
  }
  &.fl10 {
    flex: 10;
  }
  &.minimized {
    svg {
      max-width: 276px;
      height: unset;
    }
  }
`;
export { MainIconWrapper };
