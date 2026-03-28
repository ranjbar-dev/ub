import { hSmallerThan760, hSmallerThan800 } from 'styles/mediaQueries';
import styled from 'styles/styled-components';

const Wrapper = styled.footer`
  display: flex;
  justify-content: space-between;
  position: fixed;
  bottom: 7px;
  height: 14px;
  padding-left: calc(50vw - 118px);
  color: #818181;
  font-size: 10px;
  width: 100vw;
  pointer-events: none;
  background: transparent;
  /*img {
    background: var(--white);
  }
  &.loggedIn {
    img {
      background: var(--greyBackground);
    }
  }*/
  /* @media screen and (min-width: 1366px) {
    bottom: 0px;
  } */
  ${hSmallerThan800} {
    display: none;
  }
`;

export default Wrapper;
