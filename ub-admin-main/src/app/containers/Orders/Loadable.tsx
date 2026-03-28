/**
 *
 * Asynchronously loads the component for Orders
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Orders = lazyLoad(
  () => import('./index'),
  module => module.Orders,
  { fallback: <GridLoading /> },
);
