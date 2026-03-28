/**
 *
 * Asynchronously loads the component for Balances
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Balances = lazyLoad(
  () => import('./index'),
  module => module.Balances,
  { fallback: <GridLoading /> },
);
