import React from 'react';

export default function RefreshIcon(props: { color?: string }) {
  const fillColor = props.color ? props.color : '#396DE0';
  return (
    <svg width="25" height="24" viewBox="0 0 25 24">
      <defs>
        <clipPath id="refreshIconSvg">
          <rect
            id="Rectangle_4045"
            data-name="Rectangle 4045"
            width="24"
            height="25"
            fill={fillColor}
          />
        </clipPath>
      </defs>

      <g
        id="Group_7561"
        data-name="Group 7561"
        transform="translate(0 24) rotate(-90)"
      >
        <g
          id="Group_7560"
          data-name="Group 7560"
          transform="translate(0)"
          clipPath="url(#refreshIconSvg)"
        >
          <g
            id="Group_7559"
            data-name="Group 7559"
            transform="translate(3.5 5.729)"
          >
            <g
              id="Group_7557"
              data-name="Group 7557"
              transform="translate(0 0)"
            >
              <path
                id="Path_7154"
                data-name="Path 7154"
                d="M6.5,13.542A6.647,6.647,0,0,1,0,6.771,6.647,6.647,0,0,1,6.5,0,.51.51,0,0,1,7,.521a.511.511,0,0,1-.5.521A5.625,5.625,0,0,0,1,6.771,5.625,5.625,0,0,0,6.5,12.5,5.625,5.625,0,0,0,12,6.771a.5.5,0,1,1,1,0A6.647,6.647,0,0,1,6.5,13.542Z"
                transform="translate(0 0)"
                fill={fillColor}
              />
            </g>
            <g
              id="Group_7558"
              data-name="Group 7558"
              transform="translate(8 2.362)"
            >
              <path
                id="Path_7155"
                data-name="Path 7155"
                d="M8.5,5.208a.488.488,0,0,1-.354-.152L4.5,1.257.853,5.055a.486.486,0,0,1-.707,0,.537.537,0,0,1,0-.736l4-4.167a.486.486,0,0,1,.707,0l4,4.167a.535.535,0,0,1,0,.736A.484.484,0,0,1,8.5,5.208Z"
                fill={fillColor}
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
}
