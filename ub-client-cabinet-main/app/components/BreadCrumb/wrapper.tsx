import styled from 'styles/styled-components';

export const Wrapper = styled.div`
  padding: 7px 0 10px 0;
  .MuiBreadcrumbs-separator {
    margin-top: 2px;
    path {
      fill: var(--textGrey);
    }
  }
  .MuiTypography-colorPrimary {
    color: var(--textBlue);
  }
  .costumeLink {
    cursor: pointer;
  }
  .lastBread {
    span {
      color: var(--textGrey);
    }
  }
  span {
    font-weight: 600;
    font-size: 13px;
  }
  .MuiBreadcrumbs-li {
    display: flex;
  }
`;
