import React from 'react';
import ContentLoader from 'react-content-loader';

const TextLoader = (props: { width?: number; height?: number }) => {
  return (
    <ContentLoader
      speed={1}
      width={props.width ? props.width : 400}
      height={20}
      viewBox="0 0 400 20"
      backgroundColor="var(--white)"
      foregroundColor="var(--lightGrey)"
    >
      <rect
        x="2"
        y="0"
        rx="3"
        ry="3"
        width={props.width ? props.width + '' : '337'}
        height={props.height ? props.height + '' : '10'}
      />
    </ContentLoader>
  );
};

export default TextLoader;
