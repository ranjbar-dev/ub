/**
 *
 * Asynchronously loads the component for FinanceMethods
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const FinanceMethods = lazyLoad(
  () => import('./index'),
  module => module.FinanceMethods,
  { fallback: <GridLoading /> },
);
