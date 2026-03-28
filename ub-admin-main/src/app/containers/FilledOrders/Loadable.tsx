/**
 *
 * Asynchronously loads the component for FilledOrders
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const FilledOrders = lazyLoad(
  () => import('./index'),
  module => module.FilledOrders,
  { fallback: <GridLoading /> },
);
