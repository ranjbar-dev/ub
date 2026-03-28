/**
 *
 * Asynchronously loads the component for MarketTicks
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const MarketTicks = lazyLoad(
  () => import('./index'),
  module => module.MarketTicks,
  { fallback: <GridLoading /> },
);
