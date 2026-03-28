import React from 'react';
import * as ReactDOM from 'react-dom';

export const CellRenderer = (children: React.ReactNode) => {
  let wrapperDiv = document.createElement('div');
  wrapperDiv.setAttribute('style', 'overflow:visible;');
  ReactDOM.render(<>{children}</>, wrapperDiv);
  return wrapperDiv;
};
