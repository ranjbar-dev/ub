import React from 'react';

export default function DragIcon() {
  const fill = 'var(--blackText)';
  return (
    <svg viewBox="0 0 17.41 17.41">
      <defs>
        <clipPath id="clip-path11q" transform="translate(-3.29 -3.29)">
          <rect fill="none" width="24" height="24" />
        </clipPath>
      </defs>
      <title>Drag Chart</title>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Layer_2-2" data-name="Layer 2">
          <g clipPath="url(#clip-path11q)">
            <rect x="8.21" y="3.71" width="1" height="10" />
            <rect x="3.71" y="8.21" width="10" height="1" />
            <polygon points="10.47 3.18 8.71 1.41 6.94 3.18 6.23 2.48 8.71 0 11.18 2.48 10.47 3.18" />
            <polygon points="8.71 17.41 6.23 14.94 6.94 14.23 8.71 16 10.47 14.23 11.18 14.94 8.71 17.41" />
            <polygon points="2.48 11.18 0 8.71 2.48 6.23 3.18 6.94 1.41 8.71 3.18 10.47 2.48 11.18" />
            <polygon points="14.94 11.18 14.23 10.47 16 8.71 14.23 6.94 14.94 6.23 17.41 8.71 14.94 11.18" />
          </g>
        </g>
      </g>
    </svg>
  );
}
