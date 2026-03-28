import React from 'react';
import { iconStyle } from './iconStyle';

interface Props {}

function MenuWithdrawIcon (props: Props) {
  const {} = props;

  return (
    <svg style={iconStyle} width='24' height='24' viewBox='0 0 24 24'>
      <defs>
        <clipPath id='MenuPcIcon'>
          <rect
            id='Rectangle_5057'
            data-name='Rectangle 5057'
            width='24'
            height='24'
            transform='translate(60.664 425.255)'
            fill='#fff'
          />
        </clipPath>
      </defs>
      <g id='Withdraw-Icon' transform='translate(-60.664 -425.255)'>
        <g id='Group_15312' data-name='Group 15312' clipPath='url(#MenuPcIcon)'>
          <g id='Group_15311' data-name='Group 15311'>
            <path
              id='Path_14948'
              data-name='Path 14948'
              d='M80.664,428.755a.5.5,0,0,0-.5.5v1.995a.5.5,0,0,1-.505.505H65.669a.5.5,0,0,1-.5-.505v-1.995a.5.5,0,0,0-1,0v1.995a1.507,1.507,0,0,0,1.5,1.505h1.5v11a2,2,0,0,0,2,2h7.006a2,2,0,0,0,2-2v-11h1.495a1.507,1.507,0,0,0,1.505-1.505v-1.995A.5.5,0,0,0,80.664,428.755Zm-3.5,15a1,1,0,0,1-1,1H69.161a1,1,0,0,1-1-1v-11h9Z'
              fill='#fff'
            />
            <path
              id='Path_14949'
              data-name='Path 14949'
              d='M72.132,442.064v.516a.5.5,0,0,0,1,0v-.54a1.728,1.728,0,0,0,1.43-1.7v-.231a1.73,1.73,0,0,0-1.728-1.728h-.643a.729.729,0,0,1-.728-.728v-.231a.728.728,0,0,1,.728-.728h.643a.728.728,0,0,1,.728.728.5.5,0,0,0,1,0,1.727,1.727,0,0,0-1.43-1.7v-.471a.5.5,0,0,0-1,0v.447a1.726,1.726,0,0,0-1.669,1.722v.231a1.73,1.73,0,0,0,1.728,1.728h.643a.728.728,0,0,1,.728.728v.231a.729.729,0,0,1-.728.728h-.643a.729.729,0,0,1-.728-.728v-.351a.5.5,0,0,0-1,0v.351A1.727,1.727,0,0,0,72.132,442.064Z'
              fill='#fff'
            />
          </g>
        </g>
      </g>
    </svg>
  );
}

export default MenuWithdrawIcon;
