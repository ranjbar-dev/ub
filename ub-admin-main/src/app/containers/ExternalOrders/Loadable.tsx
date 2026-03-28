/**
 *
 * Asynchronously loads the component for ExternalOrders
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const ExternalOrders = lazyLoad(
  () => import('./index'),
  module => module.ExternalOrders,
  { fallback: <GridLoading /> },
);
