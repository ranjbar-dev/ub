import React from 'react';
import styled from 'styles/styled-components';

export default function BlueTitle (props: {
  number?: number;
  title: any;
  className?: string;
}) {
  return (
    <Wrapper className={props.className ?? ''}>
      {props.number && <span>{props.number} . </span>}
      <span>{props.title}</span>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  span {
    color: var(--textBlue);
    font-size: 13px !important;
    font-weight: 600;
  }
  &.mb24 {
    margin-bottom: 24px;
  }
  &.mb12 {
    margin-bottom: 24px;
  }
`;
