import React from 'react';

export default function IconAndTitle(props: { icon: any; title: any }) {
  return (
    <>
      <div className="iconWrapper">
        <img src={props.icon} alt="" />
      </div>
      <div className="title">{props.title}</div>
    </>
  );
}
