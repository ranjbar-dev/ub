import styled from 'styles/styled-components';
import { Card } from '@material-ui/core';

export const StepsWrapper = styled(Card)`
  background: var(--white);
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  .mainIcon {
    margin-bottom: 36px;
    svg {
      max-width: 336px;
    }
  }
  .inputsWrapper {
    min-width: 17vw;
    display: flex;
    flex-direction: column;
    place-content: space-around;
    align-items: center;
    .inputWrapper {
      margin-bottom: 4px;
      &.last {
        margin-bottom: 36px;
      }
    }
    .inputWithValidator {
      max-width: 278px;
      min-width: 278px;
    }
  }
  .buttonsWrapper {
    min-width: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    .MuiButton-root {
    }
    .cancelButton {
      margin-top: 20px;
    }
  }
`;
