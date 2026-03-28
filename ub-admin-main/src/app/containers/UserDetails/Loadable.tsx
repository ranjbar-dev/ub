/**
 *
 * Asynchronously loads the component for UserDetails
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const UserDetails = lazyLoad(
  () => import('./index'),
  module => module.UserDetails,
  { fallback: <GridLoading /> },
);
