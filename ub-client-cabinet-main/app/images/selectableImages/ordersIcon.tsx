import React from 'react';

export default function OrdersIcon(props: { isSelected: boolean }) {
  const isSelected = props.isSelected;
  const fillColor =
    isSelected === false ? 'var(--blackText)' : 'var(--appHeaderSelectedColor)';
  return (
    <>
      <svg width="24" height="24" viewBox="0 0 24 24">
        <defs>
          <clipPath id="clip-path3">
            <rect
              id="Rectangle_4001"
              data-name="Rectangle 4001"
              width="24"
              height="24"
              transform="translate(1300.533 539.042)"
              fill={fillColor}
            />
          </clipPath>
        </defs>
        <g
          id="Group_7440"
          data-name="Group 7440"
          transform="translate(-1300.532 -539.042)"
        >
          <g id="Group_7439" data-name="Group 7439" clipPath="url(#clip-path3)">
            <g id="Group_7438" data-name="Group 7438">
              <g id="Group_7434" data-name="Group 7434">
                <path
                  id="Path_7064"
                  data-name="Path 7064"
                  d="M1318.532,559.542h-12a2.5,2.5,0,0,1-2.5-2.5v-12a2.5,2.5,0,0,1,2.5-2.5h12a2.5,2.5,0,0,1,2.5,2.5v12A2.5,2.5,0,0,1,1318.532,559.542Zm-12-16a1.5,1.5,0,0,0-1.5,1.5v12a1.5,1.5,0,0,0,1.5,1.5h12a1.5,1.5,0,0,0,1.5-1.5v-12a1.5,1.5,0,0,0-1.5-1.5Z"
                  fill={fillColor}
                />
              </g>
              <g id="Group_7435" data-name="Group 7435">
                <path
                  id="Path_7065"
                  data-name="Path 7065"
                  d="M1316.532,547.57h-8a.5.5,0,0,1,0-1h8a.5.5,0,0,1,0,1Z"
                  fill={fillColor}
                />
              </g>
              <g id="Group_7436" data-name="Group 7436">
                <path
                  id="Path_7066"
                  data-name="Path 7066"
                  d="M1316.532,551.57h-8a.5.5,0,0,1,0-1h8a.5.5,0,0,1,0,1Z"
                  fill={fillColor}
                />
              </g>
              <g id="Group_7437" data-name="Group 7437">
                <path
                  id="Path_7067"
                  data-name="Path 7067"
                  d="M1316.532,555.627h-8a.5.5,0,0,1,0-1h8a.5.5,0,0,1,0,1Z"
                  fill={fillColor}
                />
              </g>
            </g>
          </g>
        </g>
      </svg>
    </>
  );
}
