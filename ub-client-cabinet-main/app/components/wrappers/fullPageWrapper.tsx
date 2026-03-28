import styled from 'styles/styled-components';

const FullPageWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--white);
  .head {
    /* flex:1; */
    display: flex;
    align-items: center;
    flex-direction: row-reverse;
    background: var(--white);
  }
  .body {
    flex: 15;
    display: flex;
    justify-content: center;
  }
`;
export { FullPageWrapper };
