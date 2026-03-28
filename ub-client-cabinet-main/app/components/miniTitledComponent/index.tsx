import React, { memo } from 'react';
import styled from 'styles/styled-components';

const MiniTitledComponent = (props: { title: any; children: any }) => {
  return (
    <Wrapper>
      <Title className='dragHandle'>{props.title}</Title>
      <Content>{props.children}</Content>
    </Wrapper>
  );
};
export default memo(MiniTitledComponent, () => true);

const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: var(--cardBorderRadius);
  background: var(--white);
`;

const Content = styled.div`
  padding: 5px;
  .ag-header {
    background: transparent !important;
  }
`;

const Title = styled.div`
  padding: 6px 12px;
  font-size: 9px;
  font-weight: 700;
  height: 24px;
  border-top-left-radius: var(--cardBorderRadius);
  border-top-right-radius: var(--cardBorderRadius);
  background: var(--oddRows);
  text-transform: uppercase;
  color: var(--blackText);

  span {
    font-size: 9px;
    text-transform: uppercase;
    color: var(--blackText);
    ::selection {
      background: transparent;
    }
    ::-moz-selection {
      background: transparent;
    }
  }
`;
