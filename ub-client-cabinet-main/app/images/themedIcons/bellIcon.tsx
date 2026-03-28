import React from 'react';

export default function BellIcon() {
  const fill = 'var(--blackText)';
  return (
    <svg width="24" height="24" viewBox="0 0 24 24">
      <defs>
        <clipPath id="clipPathBell">
          <rect
            id="Rectangle_4005"
            data-name="Rectangle 4005"
            width="24"
            height="24"
            transform="translate(1419.154 539.262)"
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id="Group_7491"
        data-name="Group 7491"
        transform="translate(-1419.154 -539.262)"
      >
        <g id="Group_7462" data-name="Group 7462" clipPath="url(#clipPathBell)">
          <path
            id="Path_7081"
            data-name="Path 7081"
            d="M1439.154,556.762h-1.5v-7.5a6.5,6.5,0,0,0-13,0v7.5h-1.5a.5.5,0,0,0,0,1h6a2,2,0,0,0,4,0h6a.5.5,0,0,0,0-1Zm-13.5,0v-7.5a5.5,5.5,0,0,1,11,0v7.5h-11Zm5.5,2a1,1,0,0,1-1-1h2A1,1,0,0,1,1431.154,558.762Z"
            fill={fill}
          />
        </g>
      </g>
    </svg>
  );
}
