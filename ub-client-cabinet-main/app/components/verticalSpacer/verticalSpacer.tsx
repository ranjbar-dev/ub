import styled from 'styles/styled-components';
import React from 'react';
interface _props {
  height: number;
  ofVeiwportHeight: boolean;
}
function _height (ofVeiwportHeight: boolean, height: number): string {
  return ofVeiwportHeight === true ? height + 'vh' : height + 'px';
}
const VerticalSpacer = (props: _props) => {
  const Div = styled.div`
    width: 100%;
    min-height: ${_height(props.ofVeiwportHeight, props.height)};
  `;

  return <Div />;
};
export default VerticalSpacer;
