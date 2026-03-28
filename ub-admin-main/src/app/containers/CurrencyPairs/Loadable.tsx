/**
 *
 * Asynchronously loads the component for CurrencyPairs
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const CurrencyPairs = lazyLoad(
  () => import('./index'),
  module => module.CurrencyPairs,
  { fallback: <GridLoading /> },
);
