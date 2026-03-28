/**
 *
 * Asynchronously loads the component for Deposits
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Deposits = lazyLoad(
  () => import('./index'),
  module => module.Deposits,
  { fallback: <GridLoading /> },
);
