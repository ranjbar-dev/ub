import React from 'react';

export default function ExchangeIcon(props: { isSelected: boolean }) {
  const isSelected = props.isSelected;
  const fillColor = isSelected === false ? 'var(--blackText)' : 'var(--textBlue)';
  return (
    <>
      <svg width="24" height="24" viewBox="0 0 24 24">
        <defs>
          <clipPath id="clip-path2">
            <rect
              id="Rectangle_4011"
              data-name="Rectangle 4011"
              width="24"
              height="24"
              transform="translate(1273.04 539.028)"
              fill={fillColor}
            />
          </clipPath>
        </defs>
        <g
          id="Group_7484"
          data-name="Group 7484"
          transform="translate(-1273.04 -539.028)"
        >
          <g id="Group_7483" data-name="Group 7483" clipPath="url(#clip-path2)">
            <g id="Group_7482" data-name="Group 7482">
              <path
                id="Path_7091"
                data-name="Path 7091"
                d="M1291.029,545.376h-12.752l1.994-1.995a.5.5,0,0,0-.707-.707l-2.8,2.8a.492.492,0,0,0-.21.4l0,.01a.494.494,0,0,0,.144.364l2.828,2.828a.5.5,0,0,0,.707-.707l-1.995-1.995h12.793a1.5,1.5,0,0,1,1.5,1.5v2a.5.5,0,0,0,1,0v-2A2.5,2.5,0,0,0,1291.029,545.376Z"
                fill={fillColor}
              />
              <path
                id="Path_7092"
                data-name="Path 7092"
                d="M1293.383,555.806l-2.828-2.829a.5.5,0,0,0-.707.707l1.994,1.995H1279.05a1.5,1.5,0,0,1-1.5-1.5v-2a.5.5,0,0,0-1,0v2a2.5,2.5,0,0,0,2.5,2.5H1291.8l-1.995,1.995a.5.5,0,0,0,.707.707l2.807-2.808a.489.489,0,0,0,.207-.394l0-.01A.494.494,0,0,0,1293.383,555.806Z"
                fill={fillColor}
              />
            </g>
          </g>
        </g>
      </svg>
    </>
  );
}
