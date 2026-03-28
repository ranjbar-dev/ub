import { createRoot } from 'react-dom/client';
import React from 'react';
export const CellRenderer = (children, styles?: string) => {
  const wrapperDiv = document.createElement('div');
  wrapperDiv.setAttribute('style', `overflow:visible;${styles}`);
  const root = createRoot(wrapperDiv);
  root.render(<>{children}</>);
  return wrapperDiv;
};
