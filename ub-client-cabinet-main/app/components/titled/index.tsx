import React from 'react';
import { WrapperWithTitle } from 'components/wrappers/WrapperWithTitle';
import { Divider } from '@material-ui/core';

const TitledComponent = (props: any) => (
  <WrapperWithTitle id={props.id}>
    <p className='title'>{props.title}</p>
    <Divider />
    {props.children}
  </WrapperWithTitle>
);
export default TitledComponent;
