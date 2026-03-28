import React from 'react';

export default function CheckIcon() {
  const fill = 'var(--textGrey)';
  return (
    <svg className="checkIcon" width="24" height="24" viewBox="0 0 24 24">
      <defs>
        <clipPath id="clipPathcheck">
          <rect
            id="Rectangle_4177"
            data-name="Rectangle 4177"
            width="24"
            height="24"
            transform="translate(1089.137 578.94)"
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id="Group_8560"
        data-name="Group 8560"
        transform="translate(-1089.137 -578.94)"
      >
        <g
          id="Group_7838"
          data-name="Group 7838"
          clipPath="url(#clipPathcheck)"
        >
          <g id="Group_7837" data-name="Group 7837">
            <path
              id="Path_7300"
              data-name="Path 7300"
              d="M1107.152,599.44h-12.03a2.5,2.5,0,0,1-2.5-2.5v-12a2.5,2.5,0,0,1,2.5-2.5h12.03a2.5,2.5,0,0,1,2.5,2.5v12A2.5,2.5,0,0,1,1107.152,599.44Zm-12.03-16a1.5,1.5,0,0,0-1.5,1.5v12a1.5,1.5,0,0,0,1.5,1.5h12.03a1.5,1.5,0,0,0,1.5-1.5v-12a1.5,1.5,0,0,0-1.5-1.5Z"
              fill={fill}
            />
          </g>
        </g>
      </g>
    </svg>
  );
}
