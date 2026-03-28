import React from 'react';
import styled from 'styles/styled-components';
import Anime from 'react-anime';
import fileHoverIcon from 'images/fileHoverIcon.svg';

export default function Hover () {
  return (
    <Wrapper>
      <Anime
        duration={200}
        easing='easeOutCirc'
        scale={[0.5, 1]}
        opacity={[0, 1]}
      >
        <img src={fileHoverIcon} alt='' />
      </Anime>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
`;
