import React from 'react';

interface Props {}

function InfoIcon (props: Props) {
  const {} = props;
  const color = 'var(--textGrey)';
  return (
    <svg width='16' height='16' viewBox='0 0 16 16'>
      <defs>
        <clipPath id='InfoIcon'>
          <rect
            id='Rectangle_4447'
            data-name='Rectangle 4447'
            width='16'
            height='16'
            transform='translate(390.498 105.6)'
            fill={color}
          />
        </clipPath>
      </defs>
      <g
        id='Group_8565'
        data-name='Group 8565'
        transform='translate(-390.498 -105.6)'
      >
        <g id='Group_8564' data-name='Group 8564' clipPath='url(#InfoIcon)'>
          <g id='Group_8563' data-name='Group 8563'>
            <g id='Group_8560' data-name='Group 8560'>
              <path
                id='Path_7623'
                data-name='Path 7623'
                d='M398.5,118.111a4.511,4.511,0,1,1,4.51-4.511A4.516,4.516,0,0,1,398.5,118.111Zm0-8.021a3.511,3.511,0,1,0,3.51,3.51A3.514,3.514,0,0,0,398.5,110.09Z'
                fill={color}
              />
            </g>
            <g id='Group_8561' data-name='Group 8561'>
              <path
                id='Path_7624'
                data-name='Path 7624'
                d='M398.488,114.615a.5.5,0,0,1-.5-.5v-2.907a.5.5,0,1,1,1,0v2.907A.5.5,0,0,1,398.488,114.615Z'
                fill={color}
              />
            </g>
            <g id='Group_8562' data-name='Group 8562'>
              <path
                id='Path_7625'
                data-name='Path 7625'
                d='M398.488,116.321a.5.5,0,0,1-.5-.5v-.183a.5.5,0,1,1,1,0v.183A.5.5,0,0,1,398.488,116.321Z'
                fill={color}
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
}

export default InfoIcon;
