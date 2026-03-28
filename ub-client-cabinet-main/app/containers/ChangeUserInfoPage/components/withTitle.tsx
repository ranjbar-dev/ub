import React from 'react';
import styled from 'styles/styled-components';

export default function WithTitle (props) {
  return (
    <Wrapper style={{ flex: props.flex }}>
      <div className='title'>{props.title}</div>
      {props.children}
    </Wrapper>
  );
}
const Wrapper = styled.div`
  padding: 1vh 0 0 0;
  width: 280px;
  .title {
    span {
      font-size: 12px;
      font-weight: 600;
      color: var(--textBlue);
    }
  }
`;
