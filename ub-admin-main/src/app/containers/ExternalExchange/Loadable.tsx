/**
 *
 * Asynchronously loads the component for ExternalExchange
 *
 */
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const ExternalExchange = lazyLoad(
  () => import('./index'),
  module => module.ExternalExchange,
  { fallback: <GridLoading /> },
);
