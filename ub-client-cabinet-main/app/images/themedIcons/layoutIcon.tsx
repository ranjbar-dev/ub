import React from 'react';

export default function LayOutIcon(props: { size?: string }) {
  return (
    <svg
      className="layoutIcon"
      focusable="false"
      viewBox="0 0 24 24"
      aria-hidden="true"
      role="presentation"
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
    >
      <defs>
        <clipPath id="a">
          <rect
            width="24"
            height="24"
            transform="translate(147.494 4.63)"
          ></rect>
        </clipPath>
      </defs>
      <g transform="translate(-147.494 -4.63)">
        <g>
          <path d="M165.494,8.13h-12a2.5,2.5,0,0,0-2.5,2.5v12a2.5,2.5,0,0,0,2.5,2.5h12a2.5,2.5,0,0,0,2.5-2.5v-12A2.5,2.5,0,0,0,165.494,8.13Zm-13.5,6.334h8.95V18.8h-8.95Zm1.5-5.334h7.45v4.334h-8.95V10.63A1.5,1.5,0,0,1,153.494,9.13Zm-1.5,13.5V19.8h8.95V24.13h-7.45A1.5,1.5,0,0,1,151.994,22.63Zm15,0a1.5,1.5,0,0,1-1.5,1.5h-3.55v-15h3.55a1.5,1.5,0,0,1,1.5,1.5Z"></path>
        </g>
      </g>
    </svg>
  );
}
