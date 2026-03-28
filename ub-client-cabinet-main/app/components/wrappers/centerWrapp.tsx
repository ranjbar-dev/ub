import React from 'react';
import styled from 'styles/styled-components';

export default function CenterWrapp (props) {
  return <Wrapper className={props.className}>{props.children}</Wrapper>;
}
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
  justify-content: center;
  &.popup {
    height: 565px;
  }
`;
