/**
 *
 * Asynchronously loads the component for Withdrawals
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Withdrawals = lazyLoad(
  () => import('./index'),
  module => module.Withdrawals,
  { fallback: <GridLoading /> },
);
