import React from 'react';

export default function AnimatedNoRows (props: {
  image?: any;
  texts: any[];
  icon?: any;
  isMini?: boolean;
}) {
  const style = {
    transition: 'transform 0.4s',
    transform: 'scale(0)',
    transitionTimingFunction: 'cubic-bezier(0.175, 0.885, 0.32, 1.275)',
  };
  return (
    <div className={`NoRowsWrapper ${props.isMini ? 'miniNoRows' : ''}`}>
      {setNoAddressStyle()}
      <div className='NoRowsPlaceHolder noRowsImage' style={style}>
        {props.icon}
      </div>
      <img
        src={props.image}
        style={style}
        className='NoRowsPlaceHolder noRowsImage'
      />
      {props.texts.map((item, index: number) => {
        return (
          <div
            key={'noRowsText' + index}
            className='NoRowsPlaceHolder noRowsText'
            style={{
              ...style,
              margin: '10px',
              transitionDelay: `0.${index + 1}s`,
            }}
          >
            {item}
          </div>
        );
      })}
    </div>
  );
}
const setNoAddressStyle = () => {
  setTimeout(() => {
    const ele = document.getElementsByClassName(
      'NoRowsPlaceHolder',
    ) as HTMLCollectionOf<HTMLElement>;
    for (let i = 0; i < ele.length; i++) {
      ele[i].style.transform = 'scale(1)';
    }
  }, 300);
};
