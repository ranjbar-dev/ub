import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the tradeChart state domain
 */

const selectTradeChartDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by TradeChart
 */

const makeSelectTradeChartConfig = () =>
  createSelector(selectTradeChartDomain, (substate) => {
    return substate.tradeChart ? substate.tradeChart.chartConfig : {};
  });

export { makeSelectTradeChartConfig };
