import React from 'react';

export default function TradeIcon(props: { isSelected: boolean }) {
  const isSelected = props.isSelected;
  const fillColor =
    isSelected === false ? 'var(--blackText)' : 'var(--appHeaderSelectedColor)';
  return (
    <>
      <svg width="24" height="24" viewBox="0 0 24 24">
        <defs>
          <clipPath id="clip-path">
            <rect
              id="Rectangle_4008"
              data-name="Rectangle 4008"
              width="24"
              height="24"
              transform="translate(1192.029 539.042)"
              fill={fillColor}
            />
          </clipPath>
        </defs>
        <g
          id="Group_7473"
          data-name="Group 7473"
          transform="translate(-1192.029 -539.042)"
        >
          <g id="Group_7472" data-name="Group 7472" clipPath="url(#clip-path)">
            <g id="Group_7471" data-name="Group 7471">
              <path
                id="Path_7085"
                data-name="Path 7085"
                d="M1212.028,558.536h-3.5v-7.5a.5.5,0,0,0-1,0v7.5h-3v-5.489a.5.5,0,0,0-1,0v5.489h-3v-3.5a.5.5,0,0,0-1,0v3.5h-3.5a.5.5,0,0,0,0,1h7.944a.445.445,0,0,0,.053.011.512.512,0,0,0,.054-.011h7.949a.5.5,0,0,0,0-1Z"
                fill={fillColor}
              />
              <path
                id="Path_7086"
                data-name="Path 7086"
                d="M1200.024,551.536a.5.5,0,0,0,.354-.146l7.149-7.147v2.793a.5.5,0,0,0,1,0v-4a.5.5,0,0,0-.038-.191.5.5,0,0,0-.271-.27.487.487,0,0,0-.191-.039h-4a.5.5,0,0,0,0,1h2.793l-7.149,7.147a.5.5,0,0,0,.353.853Z"
                fill={fillColor}
              />
            </g>
          </g>
        </g>
      </svg>
    </>
  );
}
