import React from 'react';
import styled from 'styles/styled-components';

export default function VerticalAlignedWrapper (props: any) {
  function returnStyle (): string {
    let style = 'display: grid;';
    if (props.maxWidth) {
      style += `max-width:${props.maxWidth}%;`;
    }
    return style;
  }
  const Div = styled.div`
    ${returnStyle()}
  `;
  return <Div>{...props.children}</Div>;
}
