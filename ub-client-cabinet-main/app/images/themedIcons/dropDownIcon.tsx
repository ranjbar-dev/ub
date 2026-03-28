import React from 'react';

export default function DropDownIcon(props: { color?: string }) {
  const fill = props.color ?? 'var(--dropDownIconBackGround)';
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="11"
      height="11"
      viewBox="0 0 11 11"
    >
      <g
        id="Group_7378"
        data-name="Group 7378"
        transform="translate(-1897 -420.667)"
      >
        <circle
          id="Ellipse_200"
          data-name="Ellipse 200"
          cx="5.5"
          cy="5.5"
          r="5.5"
          transform="translate(1897 420.667)"
          fill={fill}
        />
        <g id="sort" transform="translate(1905.296 427.933) rotate(180)">
          <path
            id="path"
            d="M0,.618.618,0,3.266,2.648.618,5.3,0,4.678l2.03-2.03Z"
            transform="translate(0 3.266) rotate(-90)"
            fill="#959595"
            fillRule="evenodd"
          />
        </g>
      </g>
    </svg>
  );
}
