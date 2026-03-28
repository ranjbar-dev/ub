/**
 *
 * Asynchronously loads the component for Reports
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Reports = lazyLoad(
  () => import('./index'),
  module => module.Reports,
  { fallback: <GridLoading /> },
);
