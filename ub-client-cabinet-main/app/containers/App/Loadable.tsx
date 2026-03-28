/**
 * Asynchronously loads the component for ChangePassword
 */

import * as React from 'react';
import loadable from 'utils/loadable';
import AppLoader from './appLoader';

export default loadable(() => import('./index'), {
  fallback: <AppLoader />,
});
