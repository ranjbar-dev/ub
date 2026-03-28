/**
*
* Asynchronously loads the component for LiquidityOrders
*
*/
import { GridLoading } from 'app/components/grid_loading/gridLoading';
import React from 'react';
import { lazyLoad } from 'utils/loadable';

export const LiquidityOrders = lazyLoad(() => import('./index'), module => module.LiquidityOrders, {fallback: <GridLoading />,},);