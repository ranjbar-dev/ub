import React from 'react';

export default function EyeIcon(props: { color?: string }) {
  const fill = props.color ?? 'var(--textBlue)';
  return (
    <svg width="24" height="24" viewBox="0 0 24 24">
      <defs>
        <clipPath id="clipPathEye">
          <rect
            id="Rectangle_4365"
            data-name="Rectangle 4365"
            width="24"
            height="24"
            transform="translate(93.657 31.346)"
            fill={fill}
          />
        </clipPath>
      </defs>
      <g
        id="Group_8225"
        data-name="Group 8225"
        transform="translate(-93.657 -31.346)"
      >
        <g id="Group_8224" data-name="Group 8224" clipPath="url(#clipPathEye)">
          <g id="Group_8223" data-name="Group 8223">
            <g id="Group_8221" data-name="Group 8221">
              <path
                id="Path_7562"
                data-name="Path 7562"
                d="M105.832,47.728c-4.656,0-6.993-4.034-7.09-4.206a.5.5,0,0,1,0-.5,8.641,8.641,0,0,1,7.086-4.06,7.753,7.753,0,0,1,6.755,4.164.5.5,0,0,1,0,.438A7.753,7.753,0,0,1,105.832,47.728Zm-6.058-4.447a7.628,7.628,0,0,0,6.058,3.447,6.7,6.7,0,0,0,5.736-3.382,6.7,6.7,0,0,0-5.736-3.383A7.669,7.669,0,0,0,99.774,43.281Z"
                fill={fill}
              />
            </g>
            <g id="Group_8222" data-name="Group 8222">
              <path
                id="Path_7563"
                data-name="Path 7563"
                d="M105.832,45.846a2.5,2.5,0,1,1,2.5-2.5A2.5,2.5,0,0,1,105.832,45.846Zm0-4a1.5,1.5,0,1,0,1.5,1.5A1.5,1.5,0,0,0,105.832,41.846Z"
                fill={fill}
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
}
