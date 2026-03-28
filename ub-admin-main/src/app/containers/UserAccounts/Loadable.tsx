/**
 *
 * Asynchronously loads the component for UserAccounts
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const UserAccounts = lazyLoad(
  () => import('./index'),
  module => module.UserAccounts,
  { fallback: <GridLoading /> },
);
