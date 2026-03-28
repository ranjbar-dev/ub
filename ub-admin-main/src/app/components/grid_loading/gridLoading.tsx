import React from 'react';

import { Wrapper } from './wrapper';

const GridLoading = (props: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <Wrapper {...props}>
      <div className="lds-ellipsis">
        <div></div>
        <div></div>
        <div></div>
        <div></div>
      </div>
    </Wrapper>
  );
};
export { GridLoading };
