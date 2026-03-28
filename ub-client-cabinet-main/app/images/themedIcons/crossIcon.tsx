import React from 'react';

export default function CrossIcon (props: { color?: string; style?: any }) {
  const fill = props.color ?? 'var(--textGrey)';
  return (
    <svg width='16' height='16' style={props.style} viewBox='0 0 16 16'>
      <defs>
        <clipPath id='clipPathCross'>
          <rect
            id='Rectangle_4437'
            data-name='Rectangle 4437'
            width='16'
            height='16'
            transform='translate(366.103 79.007)'
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id='Group_8544'
        data-name='Group 8544'
        transform='translate(-366.103 -79.007)'
      >
        <g
          id='Group_8543'
          data-name='Group 8543'
          clipPath='url(#clipPathCross)'
        >
          <path
            id='Path_7616'
            data-name='Path 7616'
            d='M374.81,87.007l3.646-3.647a.5.5,0,0,0-.707-.707L374.1,86.3l-3.647-3.647a.5.5,0,1,0-.707.707l3.647,3.647-3.647,3.646a.5.5,0,0,0,.354.854.5.5,0,0,0,.353-.147l3.647-3.646,3.646,3.646a.5.5,0,0,0,.707-.707Z'
            fill={fill}
          />
        </g>
      </g>
    </svg>
  );
}
