import React from 'react';
import { Route, Redirect, RouteProps } from 'react-router-dom';
import { LocalStorageKeys } from 'services/constants';
import { AppPages } from 'app/constants';

interface PrivateRouteProps extends RouteProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  component: React.ComponentType<any>;
}

/** Returns true only when a non-expired JWT is present in localStorage. */
function isAuthenticated(): boolean {
  const token = localStorage.getItem(LocalStorageKeys.ACCESS_TOKEN);
  if (!token) return false;

  try {
    // JWT structure: header.payload.signature (all base64url-encoded)
    const payloadBase64 = token.split('.')[1];
    if (!payloadBase64) return false;

    // base64url → base64 → JSON
    const padded = payloadBase64.replace(/-/g, '+').replace(/_/g, '/');
    const json = atob(padded);
    const payload: { exp?: number } = JSON.parse(json);

    if (typeof payload.exp !== 'number') {
      // Token has no exp claim — treat as valid (non-standard tokens)
      return true;
    }

    // exp is in seconds; Date.now() is in milliseconds
    return payload.exp * 1000 > Date.now();
  } catch {
    // Malformed token — treat as unauthenticated
    return false;
  }
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ component: Component, ...rest }) => (
  <Route
    {...rest}
    render={(props) =>
      isAuthenticated() ? (
        <Component {...props} />
      ) : (
        <Redirect to={AppPages.RootPage} />
      )
    }
  />
);

export default PrivateRoute;
