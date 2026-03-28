import styled from 'styles/styled-components';

const Wrapper = styled.div`
  padding: 2px;
  .MuiInput-underline:before {
    display: none;
  }
  .MuiSelect-select {
    color: #818181;
  }
  .MuiSelect-icon {
    top: calc(50% - 14px);
    right: 20px;
    left: unset;
  }
  .MuiFormControl-root {
    min-width: 90px !important;
  }
  .MuiInput-underline:after {
    display: none;
  }
`;

export default Wrapper;
