import React, { memo, useMemo } from 'react';

import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';

import TradeLayout from './layout/tradeLayout';

const StreamComponentsWrapper = () => {
  return (
    <>
      {useMemo(
        () => (
          <TradeLayout />
        ),
        [],
      )}
    </>
  );
};
export default memo(StreamComponentsWrapper);
