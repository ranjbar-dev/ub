import React from 'react';

export default function FundsIcon(props: { isSelected: boolean }) {
  const isSelected = props.isSelected;
  const fillColor =
    isSelected === false ? 'var(--blackText)' : 'var(--appHeaderSelectedColor)';
  return (
    <>
      <svg width="24" height="24" viewBox="0 0 24 24">
        <defs>
          <clipPath id="clip-path4">
            <rect
              id="Rectangle_4009"
              data-name="Rectangle 4009"
              width="24"
              height="24"
              transform="translate(1219.036 539.036)"
              fill={fillColor}
            />
          </clipPath>
        </defs>
        <g
          id="Group_7476"
          data-name="Group 7476"
          transform="translate(-1219.036 -539.036)"
        >
          <g id="Group_7475" data-name="Group 7475" clipPath="url(#clip-path4)">
            <g id="Group_7474" data-name="Group 7474">
              <path
                id="Path_7087"
                data-name="Path 7087"
                d="M1237.036,544.6h-1.5v-.632a1.5,1.5,0,0,0-1.5-1.5h-9a2.5,2.5,0,0,0-2.5,2.5V557.1a2.5,2.5,0,0,0,2.5,2.5h12a2.5,2.5,0,0,0,2.5-2.5v-10A2.5,2.5,0,0,0,1237.036,544.6Zm-12-1.132h9a.5.5,0,0,1,.5.5v.632h-9.5a2.476,2.476,0,0,0-1.5.512v-.144A1.5,1.5,0,0,1,1225.036,543.47Zm12,15.132h-12a1.5,1.5,0,0,1-1.5-1.5v-10a1.5,1.5,0,0,1,1.5-1.5h12a1.5,1.5,0,0,1,1.5,1.5v5.364h-2.5a1.5,1.5,0,0,0-1.5,1.5v.123a1.5,1.5,0,0,0,1.5,1.5h2.5V557.1A1.5,1.5,0,0,1,1237.036,558.6Zm1.5-4.013h-2.5a.5.5,0,0,1-.5-.5v-.123a.5.5,0,0,1,.5-.5h2.5Z"
                fill={fillColor}
              />
            </g>
          </g>
        </g>
      </svg>
    </>
  );
}
