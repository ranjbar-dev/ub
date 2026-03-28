/**
 *
 * Asynchronously loads the component for Billing
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Billing = lazyLoad(
  () => import('./index'),
  module => module.Billing,
  { fallback: <GridLoading /> },
);
