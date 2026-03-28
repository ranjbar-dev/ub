import styled from 'styles/styled-components';

export const AnimatedAppear = styled.div`
  opacity: 0;
  transition: opacity 1s;
  &.visible {
    opacity: 1;
  }
  .scaled {
    transition: transform 0.8s;
    &.In {
      opacity: 1;
      transform: scale(1);
    }
    &.Out {
      opacity: 0;
      transform: scale(0.8);
    }
    &.delayed1 {
      transition-delay: 0s;
    }
    &.delayed2 {
      transition-delay: 0.2s;
    }
    &.delayed3 {
      transition-delay: 0.3s;
    }
    &.delayed4 {
      transition-delay: 0s;
    }
    &.delayed5 {
      transition-delay: 0s;
    }
  }
`;
