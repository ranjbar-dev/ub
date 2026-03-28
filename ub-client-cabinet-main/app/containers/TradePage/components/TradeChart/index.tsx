/*
 *
 * TradeChart
 *
 */

import React from 'react';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import reducer from './reducer';
import saga from './saga';
import styled from 'styles/styled-components';
import UBChartContainer from './mainChart';
import DragIcon from 'images/themedIcons/dragIcon';

function TradeChart (props: { enabled: boolean }) {
  useInjectReducer({ key: 'tradeChart', reducer: reducer });
  useInjectSaga({ key: 'tradeChart', saga: saga });

  return (
    <Wrapper>
      <div className='dragHandle'>
        <DragIcon />
      </div>
      {props.enabled && <UBChartContainer />}
    </Wrapper>
  );
}

export default TradeChart;
const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  background: var(--oddRows);
  border-radius: var(--cardBorderRadius);
  .dragHandle {
    position: absolute;
    z-index: 1;
    right: 40px;
    width: 35px;
    height: 35px;
    overflow: hidden;
    top: -5px;
    padding: 8px;
    * {
      fill: var(--dragIconColor);
    }
  }
`;
