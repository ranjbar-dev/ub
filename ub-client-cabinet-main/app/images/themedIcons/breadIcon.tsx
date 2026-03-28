import React from 'react';

export default function BreadIcon() {
  const fill = 'var(--textGrey)';
  return (
    <div style={{ transform: 'rotate(-90deg)' }}>
      <svg width="16" height="16" viewBox="0 0 16 16">
        <defs>
          <clipPath id="clipPathExpandMore">
            <rect
              id="Rectangle_4145"
              data-name="Rectangle 4145"
              width="16"
              height="16"
              fill={fill}
            />
          </clipPath>
        </defs>
        <g
          id="Group_7797"
          data-name="Group 7797"
          transform="translate(16 16) rotate(180)"
        >
          <g
            id="Group_7796"
            data-name="Group 7796"
            clipPath="url(#clipPathExpandMore)"
          >
            <g
              id="Group_7795"
              data-name="Group 7795"
              transform="translate(3.501 5.5)"
            >
              <path
                id="Path_7288"
                data-name="Path 7288"
                d="M1167.26,587.832a.5.5,0,0,1,.354.146l4,4a.5.5,0,0,1-.707.707l-3.647-3.646-3.646,3.646a.5.5,0,1,1-.707-.707l4-4A.5.5,0,0,1,1167.26,587.832Z"
                transform="translate(-1162.761 -587.832)"
                fill={fill}
              />
            </g>
          </g>
        </g>
      </svg>
    </div>
  );
}
