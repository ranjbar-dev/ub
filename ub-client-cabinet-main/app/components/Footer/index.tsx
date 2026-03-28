import React from 'react';

import Wrapper from './Wrapper';
//import copy from 'images/copyright.svg';
const Footer = props => {
  return (
    <Wrapper {...props}>
      {/*<img src={copy} alt='' />
       */}
      © 2021 UNITEDBIT EXCHANGE. All Rights Reserved. V: {process.env.VERSION}
    </Wrapper>
  );
};

export default Footer;
