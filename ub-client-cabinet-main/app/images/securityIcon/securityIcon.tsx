import React from 'react';
import { SecurityLevel } from 'containers/AcountPage/constants';

export default function SecurityIcon(props: { level: SecurityLevel }) {
  let fillColor;

  switch (props.level) {
    case SecurityLevel.LOW:
      fillColor = 'red';
      break;
    case SecurityLevel.MEDIUM:
      fillColor = '#f69906';
      break;
    case SecurityLevel.High:
      fillColor = 'green';
      break;

    default:
      break;
  }

  return (
    <svg
      id="Group_7754"
      data-name="Group 7754"
      style={{ marginTop: '-4px' }}
      width="24"
      height="24"
      viewBox="0 0 24 24"
    >
      <defs>
        <clipPath id="SecurityIcon">
          <rect
            id="Rectangle_4133"
            data-name="Rectangle 4133"
            width="24"
            height="24"
            fill={fillColor}
          />
        </clipPath>
      </defs>
      <g id="Group_7753" data-name="Group 7753" clipPath="url(#SecurityIcon)">
        <g
          id="Group_7752"
          data-name="Group 7752"
          transform="translate(5.562 3.5)"
        >
          <g id="Group_7750" data-name="Group 7750">
            <path
              id="Path_7266"
              data-name="Path 7266"
              d="M873.255,558.833a.507.507,0,0,1-.224-.052,11.182,11.182,0,0,1-6.214-10.055v-3.258a1.492,1.492,0,0,1,1.49-1.49h.255a2.887,2.887,0,0,0,2.129-1.01,3.465,3.465,0,0,1,5.129,0,2.881,2.881,0,0,0,2.128,1.01h.255a1.491,1.491,0,0,1,1.49,1.49v3.258a11.182,11.182,0,0,1-6.214,10.055A.507.507,0,0,1,873.255,558.833Zm-4.948-13.855a.49.49,0,0,0-.49.49v3.258a10.183,10.183,0,0,0,5.438,9.045,10.186,10.186,0,0,0,5.438-9.045v-3.258a.491.491,0,0,0-.49-.49h-.255a3.876,3.876,0,0,1-2.869-1.337,2.545,2.545,0,0,0-3.648,0,3.877,3.877,0,0,1-2.869,1.337Z"
              transform="translate(-866.817 -541.833)"
              fill={fillColor}
            />
          </g>
          <g
            id="Group_7751"
            data-name="Group 7751"
            transform="translate(3.725 5.433)"
          >
            <path
              id="Path_7267"
              data-name="Path 7267"
              d="M872.918,552.725a1,1,0,0,1-.7-.292l-1.525-1.525a.5.5,0,1,1,.707-.707l1.525,1.525,3.784-4.291a.5.5,0,0,1,.75.662l-3.8,4.3a1,1,0,0,1-.716.332Z"
              transform="translate(-870.542 -547.266)"
              fill={fillColor}
            />
          </g>
        </g>
      </g>
    </svg>
  );
}
