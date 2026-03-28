/**
 *
 * Asynchronously loads the component for LoginHistory
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const LoginHistory = lazyLoad(
  () => import('./index'),
  module => module.LoginHistory,
  { fallback: <GridLoading /> },
);
