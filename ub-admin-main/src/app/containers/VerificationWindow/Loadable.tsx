/**
 *
 * Asynchronously loads the component for VerificationWindow
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const VerificationWindow = lazyLoad(
  () => import('./index'),
  module => module.VerificationWindow,
  { fallback: <GridLoading /> },
);
