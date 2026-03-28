/**
 * Asynchronously loads the component for LoginPage
 */

import { GridLoading } from 'app/components/grid_loading/gridLoading';
import * as React from 'react';
import { lazyLoad } from 'utils/loadable';

export const LoginPage = lazyLoad(
  () => import('./index'),
  module => module.LoginPage,
  {
    fallback: <GridLoading />,
  },
);
