/**
 * Shimmer
 */
import React from 'react';
// Style
import styled, { css, keyframes } from 'styled-components';
import { shadeColor } from './shadeColor';

// Utils

export const Shimmer = ({ duration, height, hexColour, width }) => {
  return (
    <StyledShimmer
      duration={duration}
      height={height}
      hexColor={hexColour}
      width={width}
    />
  );
};

const placeholderShimmer = keyframes`
    0% {
      background-position: -468px 0;
    }

    100% {
      background-position: 468px 0;
    }
`;
const defaultColour = '#f6f7f8';

const StyledShimmer = styled.div<{
  duration: number;
  height: string;
  hexColor: string;
  width: string;
}>`
  animation-duration: ${({ duration }) => (duration ? duration + 'ms' : '1s')};
  animation-fill-mode: forwards;
  animation-iteration-count: infinite;
  animation-name: ${placeholderShimmer};
  animation-timing-function: linear;

  ${({ hexColor }) => {
    return css`
      background: ${hexColor ?? defaultColour};
      background-image: linear-gradient(
        to right,
        ${hexColor ?? defaultColour} 0%,
        ${shadeColor(hexColor ?? defaultColour, -8)} 20%,
        ${hexColor ?? defaultColour} 40%,
        ${hexColor ?? defaultColour} 100%
      );
    `;
  }}

  background-repeat: no-repeat;
  background-size: 800px 104px;
  display: inline-block;
  height: ${({ height }) => height ?? '100%'};
  position: relative;
  width: ${({ width }) => width ?? '100%'};
`;
