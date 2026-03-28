import React from 'react';

export default function CheckedIcon() {
  const fill = 'var(--textBlue)';
  return (
    <svg width="24" height="24" viewBox="0 0 24 24">
      <defs>
        <clipPath id="clipPathchecked">
          <rect
            id="Rectangle_4426"
            data-name="Rectangle 4426"
            width="24"
            height="24"
            transform="translate(25.126 228.272)"
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id="Group_8335"
        data-name="Group 8335"
        transform="translate(-25.126 -228.272)"
      >
        <g
          id="Group_8334"
          data-name="Group 8334"
          clipPath="url(#clipPathchecked)"
        >
          <g id="Group_8332" data-name="Group 8332">
            <path
              id="Path_7600"
              data-name="Path 7600"
              d="M43.141,248.772H31.111a2.5,2.5,0,0,1-2.5-2.5v-12a2.5,2.5,0,0,1,2.5-2.5h12.03a2.5,2.5,0,0,1,2.5,2.5v12A2.5,2.5,0,0,1,43.141,248.772Zm-12.03-16a1.5,1.5,0,0,0-1.5,1.5v12a1.5,1.5,0,0,0,1.5,1.5h12.03a1.5,1.5,0,0,0,1.5-1.5v-12a1.5,1.5,0,0,0-1.5-1.5Z"
              fill={fill}
            />
          </g>
          <g id="Group_8333" data-name="Group 8333">
            <path
              id="Path_7601"
              data-name="Path 7601"
              d="M35.945,244.6a1.5,1.5,0,0,1-1.1-.477l-3.11-3.34a.5.5,0,0,1,.733-.682l3.109,3.34a.5.5,0,0,0,.394.159.5.5,0,0,0,.373-.2l5.4-7.256a.5.5,0,0,1,.8.6L37.15,244a1.5,1.5,0,0,1-1.122.6C36,244.6,35.973,244.6,35.945,244.6Z"
              fill={fill}
            />
          </g>
        </g>
      </g>
    </svg>
  );
}
