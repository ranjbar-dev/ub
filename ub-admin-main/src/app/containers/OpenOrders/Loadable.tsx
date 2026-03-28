/**
 *
 * Asynchronously loads the component for OpenOrders
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const OpenOrders = lazyLoad(
  () => import('./index'),
  module => module.OpenOrders,
  { fallback: <GridLoading /> },
);
