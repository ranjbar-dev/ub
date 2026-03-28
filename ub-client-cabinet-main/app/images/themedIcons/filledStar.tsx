import React from 'react';

export default function FilledStar(props: { size?: string }) {
  return (
    <svg
      width={props.size ?? '24'}
      height={props.size ?? '24'}
      viewBox="0 0 24 24"
    >
      <defs>
        <clipPath id="filledStar">
          <rect
            id="Rectangle_4064"
            data-name="Rectangle 4064"
            width="24"
            height="24"
            transform="translate(1037.504 540.378)"
            fill="var(--textBlue)"
          />
        </clipPath>
      </defs>
      <g
        id="Group_7644"
        data-name="Group 7644"
        transform="translate(-1037.504 -540.378)"
      >
        <g id="Group_7643" data-name="Group 7643" clipPath="url(#filledStar)">
          <path
            id="Path_7218"
            data-name="Path 7218"
            d="M1050.564,544.713l1.77,3.586a1.182,1.182,0,0,0,.89.647l3.958.575a1.183,1.183,0,0,1,.656,2.017l-2.865,2.792a1.185,1.185,0,0,0-.34,1.046l.677,3.942a1.183,1.183,0,0,1-1.716,1.246l-3.54-1.861a1.182,1.182,0,0,0-1.1,0l-3.541,1.861a1.182,1.182,0,0,1-1.715-1.246l.676-3.942a1.18,1.18,0,0,0-.34-1.046l-2.864-2.792a1.183,1.183,0,0,1,.655-2.017l3.958-.575a1.185,1.185,0,0,0,.891-.647l1.77-3.586A1.182,1.182,0,0,1,1050.564,544.713Z"
            fill="var(--textBlue)"
          />
        </g>
      </g>
    </svg>
  );
}
