/**
 *
 * Asynchronously loads the component for Admins
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const Admins = lazyLoad(
  () => import('./index'),
  module => module.Admins,
  { fallback: <GridLoading /> },
);
