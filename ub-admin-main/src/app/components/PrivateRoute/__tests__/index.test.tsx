import React from 'react';
import { render, screen } from '@testing-library/react';
import { Router, Switch, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import PrivateRoute from '../index';
import { LocalStorageKeys } from 'services/constants';
import { AppPages } from 'app/constants';

const MockComponent: React.FC = () => <div>Protected Content</div>;
const LoginPage: React.FC = () => <div>Login Page</div>;

function renderWithRouter(token: string | null, path = '/test') {
  const history = createMemoryHistory({ initialEntries: [path] });

  if (token) {
    localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, token);
  } else {
    localStorage.removeItem(LocalStorageKeys.ACCESS_TOKEN);
  }

  return render(
    <Router history={history}>
      <Switch>
        <Route exact path={AppPages.RootPage} component={LoginPage} />
        <PrivateRoute path="/test" component={MockComponent} />
      </Switch>
    </Router>,
  );
}

afterEach(() => {
  localStorage.clear();
});

describe('PrivateRoute', () => {
  describe('with a valid token', () => {
    it('renders the protected component', () => {
      renderWithRouter('mock-jwt');
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });

    it('does not redirect to login', () => {
      renderWithRouter('mock-jwt');
      expect(screen.queryByText('Login Page')).not.toBeInTheDocument();
    });
  });

  describe('without a token', () => {
    it('redirects to the login page', () => {
      renderWithRouter(null);
      expect(screen.getByText('Login Page')).toBeInTheDocument();
    });

    it('does not render the protected component', () => {
      renderWithRouter(null);
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });
  });

  describe('token removal mid-session', () => {
    it('renders protected content with token then redirects after token removal', () => {
      const history = createMemoryHistory({ initialEntries: ['/test'] });
      localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, 'mock-jwt');

      const { rerender } = render(
        <Router history={history}>
          <Switch>
            <Route exact path={AppPages.RootPage} component={LoginPage} />
            <PrivateRoute path="/test" component={MockComponent} />
          </Switch>
        </Router>,
      );

      expect(screen.getByText('Protected Content')).toBeInTheDocument();

      localStorage.removeItem(LocalStorageKeys.ACCESS_TOKEN);

      rerender(
        <Router history={history}>
          <Switch>
            <Route exact path={AppPages.RootPage} component={LoginPage} />
            <PrivateRoute path="/test" component={MockComponent} />
          </Switch>
        </Router>,
      );

      expect(screen.getByText('Login Page')).toBeInTheDocument();
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });
  });
});
