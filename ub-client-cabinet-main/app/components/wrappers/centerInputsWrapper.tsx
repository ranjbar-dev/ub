import styled from 'styles/styled-components';

const CenterInputsWrapper = styled.div`
  min-width: 17vw;
  flex: 4;
  display: flex;
  flex-direction: column;
  place-content: space-around;
  padding-bottom: 4vh;
  align-items: center;
  color: var(--textGrey);
  min-width: 280px;
  /* align-items: center; */
  .MuiInputBase-root {
    height: 40px;
  }
  .MuiSelect-selectMenu:focus {
    background-color: transparent;
  }
  .phoneNumber {
    color: var(--textBlue);
    font-weight: 600;
  }
  .boldGrey {
    font-weight: 600;
    color: var(--textGrey);
  }
  &.noPadd {
    flex: 6;
    place-content: center;
    padding: 0;
  }
  &.pt3 {
    padding-top: 3px;
  }
  &.even {
    place-content: space-evenly;
  }
  &.start {
    place-content: flex-start;
  }
  .inputWrapper {
    max-width: 278px;
    min-width: 278px;
  }
`;
export { CenterInputsWrapper };
