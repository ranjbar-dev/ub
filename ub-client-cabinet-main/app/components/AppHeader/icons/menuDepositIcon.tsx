import React from 'react';
import { iconStyle } from './iconStyle';

interface Props {}

function MenuDepositIcon (props: Props) {
  const {} = props;

  return (
    <svg style={iconStyle} width='24' height='24' viewBox='0 0 24 24'>
      <defs>
        <clipPath id='MenuDepositIcon'>
          <rect
            id='Rectangle_5056'
            data-name='Rectangle 5056'
            width='24'
            height='24'
            transform='translate(96.202 425.255)'
            fill='#fff'
          />
        </clipPath>
      </defs>
      <g id='Deposit-Icon' transform='translate(-96.202 -425.255)'>
        <g
          id='Group_15309'
          data-name='Group 15309'
          clipPath='url(#MenuDepositIcon)'
        >
          <path
            id='Path_14947'
            data-name='Path 14947'
            d='M114.2,432.755h-1.651l-1.386-3.676a.5.5,0,0,0-.7-.268l-7.627,3.944h-.631a2.506,2.506,0,0,0-2.5,2.5v7.992a2.506,2.506,0,0,0,2.5,2.5H114.2a2.506,2.506,0,0,0,2.5-2.5v-7.992A2.506,2.506,0,0,0,114.2,432.755Zm-3.775-2.8,1.055,2.8h-6.465Zm3.775,14.8H102.206a1.506,1.506,0,0,1-1.5-1.5v-7.992a1.506,1.506,0,0,1,1.5-1.5H114.2a1.505,1.505,0,0,1,1.5,1.5v3.5h-2.5a1.5,1.5,0,0,0,0,3h2.5v1.5A1.505,1.505,0,0,1,114.2,444.755Zm1.5-4h-2.5a.5.5,0,0,1,0-1h2.5Z'
            fill='#fff'
          />
        </g>
      </g>
    </svg>
  );
}

export default MenuDepositIcon;
