import React from 'react';

export default function ReplayIcon(props: { disabled?: boolean }) {
  let fill = 'var(--textBlue)';
  if (props.disabled) {
    fill = 'var(--textGrey)';
  }

  return (
    <svg width="24" height="24" viewBox="0 0 24 24">
      <defs>
        <clipPath id="clipPathreplay">
          <rect
            id="Rectangle_4045"
            data-name="Rectangle 4045"
            width="24"
            height="24"
            transform="translate(1124.942 540)"
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id="Group_7561"
        data-name="Group 7561"
        transform="translate(-540 1148.942) rotate(-90)"
      >
        <g
          id="Group_7560"
          data-name="Group 7560"
          clipPath="url(#clipPathreplay)"
        >
          <g id="Group_7559" data-name="Group 7559">
            <g id="Group_7557" data-name="Group 7557">
              <path
                id="Path_7154"
                data-name="Path 7154"
                d="M1134.942,558.5a6.5,6.5,0,0,1,0-13,.5.5,0,0,1,0,1,5.5,5.5,0,1,0,5.5,5.5.5.5,0,0,1,1,0A6.508,6.508,0,0,1,1134.942,558.5Z"
                fill={fill}
              />
            </g>
            <g
              id="Group_7558"
              data-name="Group 7558"
              transform="translate(0 -1)"
            >
              <path
                id="Path_7155"
                data-name="Path 7155"
                d="M1144.942,553.767a.5.5,0,0,1-.354-.146l-3.646-3.647-3.647,3.647a.5.5,0,1,1-.707-.707l4-4a.5.5,0,0,1,.707,0l4,4a.5.5,0,0,1-.353.853Z"
                fill={fill}
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
}
