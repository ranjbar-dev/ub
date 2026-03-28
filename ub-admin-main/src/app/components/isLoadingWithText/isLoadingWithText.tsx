import React from 'react';

import LoadingInButton from '../loadingInButton/loadingInButton';

/**
 * Spinner-inside-button loading indicator.
 * Shows a spinner when loading, fades text out, and shows text when idle.
 *
 * @example
 * ```tsx
 * <IsLoadingWithText isLoading={isFetching} text="Save" />
 * ```
 */
export default function IsLoadingWithText(props: {
  isLoading: boolean;
  text: React.ReactNode;
}) {
  let loadingOpacity = props.isLoading === true ? 1 : 0;
  let textOpacity = props.isLoading === false ? 1 : 0;
  return (
    <div>
      <span
        className="loadingCircle"
        style={{
          opacity: loadingOpacity,
          position: 'absolute',
          left: 'calc(50% - 8px)',
          top: '8px',
        }}
      >
        <LoadingInButton />
      </span>
      <span style={{ opacity: textOpacity }}>{props.text}</span>
    </div>
  );
}
