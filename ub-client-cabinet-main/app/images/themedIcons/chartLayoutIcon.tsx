import React from 'react';

export default function ChartLayOutIcon(props: { size?: string }) {
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
      // style="margin-right: 10px;"
    >
      <defs>
        <clipPath id="abb">
          <rect
            width="24"
            height="24"
            transform="translate(52.082 1.751)"
          ></rect>
        </clipPath>
      </defs>
      <g transform="translate(-52.082 -1.751)">
        <g>
          <path d="M70.082,5.251h-12a2.5,2.5,0,0,0-2.5,2.5v12a2.5,2.5,0,0,0,2.5,2.5h12a2.5,2.5,0,0,0,2.5-2.5v-12A2.5,2.5,0,0,0,70.082,5.251Zm-12,1h7.45v7h-8.95v-5.5A1.5,1.5,0,0,1,58.082,6.251Zm-1.5,13.5v-5.5h8.95v7h-7.45A1.5,1.5,0,0,1,56.582,19.751Zm15,0a1.5,1.5,0,0,1-1.5,1.5h-3.55v-15h3.55a1.5,1.5,0,0,1,1.5,1.5Z"></path>
        </g>
      </g>
    </svg>
  );
}
