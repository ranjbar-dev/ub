import React from 'react';
import { Route, Redirect, RouteProps } from 'react-router-dom';
import { LocalStorageKeys } from 'services/constants';
import { AppPages } from 'app/constants';

interface PrivateRouteProps extends RouteProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  component: React.ComponentType<any>;
}

/**
 * Route wrapper that redirects to login if user has no access token.
 * Checks localStorage for ACCESS_TOKEN — mirrors the auth check used
 * throughout the app (ShowSideNav, apiService headers).
 */
const PrivateRoute: React.FC<PrivateRouteProps> = ({ component: Component, ...rest }) => (
  <Route
    {...rest}
    render={(props) =>
      localStorage[LocalStorageKeys.ACCESS_TOKEN] ? (
        <Component {...props} />
      ) : (
        <Redirect to={AppPages.RootPage} />
      )
    }
  />
);

export default PrivateRoute;
