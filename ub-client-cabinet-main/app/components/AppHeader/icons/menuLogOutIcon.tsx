import React from 'react';
import { iconStyle } from './iconStyle';

interface Props {}

function MenuLogOutIcon (props: Props) {
  const {} = props;

  return (
    <svg style={iconStyle} width='24' height='24' viewBox='0 0 24 24'>
      <defs>
        <clipPath id='MenuLogOutIcon'>
          <rect
            id='Rectangle_4478'
            data-name='Rectangle 4478'
            width='24'
            height='24'
            transform='translate(273.893 227.522)'
            fill='#fff'
          />
        </clipPath>
      </defs>
      <g
        id='Group_8659'
        data-name='Group 8659'
        transform='translate(-273.893 -227.522)'
      >
        <g
          id='Group_8658'
          data-name='Group 8658'
          clipPath='url(#MenuLogOutIcon)'
        >
          <g id='Group_8657' data-name='Group 8657'>
            <g id='Group_8655' data-name='Group 8655'>
              <g id='Group_8654' data-name='Group 8654'>
                <path
                  id='Path_7651'
                  data-name='Path 7651'
                  d='M281.576,243.024l-2.961-2.961h10.158a.471.471,0,1,0,0-.941H278.615l2.961-2.962a.47.47,0,0,0-.665-.665l-3.1,3.1a1.393,1.393,0,0,0-.4.9.237.237,0,0,0,0,.186,1.389,1.389,0,0,0,.4.9l3.1,3.1a.471.471,0,0,0,.665-.666Z'
                  fill='#fff'
                />
              </g>
            </g>
            <g id='Group_8656' data-name='Group 8656'>
              <path
                id='Path_7652'
                data-name='Path 7652'
                d='M291.888,248.022h-6a2.5,2.5,0,0,1-2.5-2.5v-2a.5.5,0,1,1,1,0v2a1.5,1.5,0,0,0,1.5,1.5h6a1.5,1.5,0,0,0,1.5-1.5v-12a1.5,1.5,0,0,0-1.5-1.5h-6a1.5,1.5,0,0,0-1.5,1.5v2a.5.5,0,0,1-1,0v-2a2.5,2.5,0,0,1,2.5-2.5h6a2.5,2.5,0,0,1,2.5,2.5v12A2.5,2.5,0,0,1,291.888,248.022Z'
                fill='#fff'
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
}

export default MenuLogOutIcon;
