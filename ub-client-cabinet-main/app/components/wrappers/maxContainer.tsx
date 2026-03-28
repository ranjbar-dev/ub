import styled from 'styles/styled-components';

const MaxContainer = styled.div`
  width: 100%;
  display: flex;
  max-height: var(--maxHeight);
  flex-direction: column;
  align-items: center;
  height: 100%;
  &.max740 {
    max-height: 740px;
  }
  &.mh1 {
    max-height: 590px;
  }
`;
export { MaxContainer };
