import styled from 'styles/styled-components';

const MaxWidthWrapper = styled.div`
  width: 95vw;
  @media screen and (min-width: 1200px) {
    width: 1200px;
  }

  margin: auto;
  input,
  textarea {
    font-weight: 600;
  }
`;
export { MaxWidthWrapper };
