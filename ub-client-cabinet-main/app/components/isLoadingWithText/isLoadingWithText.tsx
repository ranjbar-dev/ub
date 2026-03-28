import React from 'react';
import LoadingInButton from 'components/loadingInButton/loadingInButton';

export default function IsLoadingWithText(props: {
  isLoading: boolean;
  text: any;
}) {
  const loadingOpacity = props.isLoading === true ? 1 : 0;
  const textOpacity = props.isLoading === false ? 1 : 0;
  return (
    <div>
      <span
        className="loadingCircle"
        style={{
          opacity: loadingOpacity,
          position: 'absolute',
          left: 'calc(50% - 8px)',
          top: '12px',
        }}
      >
        <LoadingInButton />
      </span>
      <span style={{ opacity: textOpacity }}>{props.text}</span>
    </div>
  );
}
