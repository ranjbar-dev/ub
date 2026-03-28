import { Link } from 'app/components/Link';
import { createMemoryHistory } from 'history';
import React from 'react';
import { HelmetProvider } from 'react-helmet-async';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import renderer from 'react-test-renderer';
import { configureAppStore } from 'store/configureStore';
import { ThemeProvider } from 'styles/theme/ThemeProvider';

import { NotFoundPage } from '..';

const history = createMemoryHistory();
const store = configureAppStore({}, history);

const renderPage = () =>
  renderer.create(
    <Provider store={store}>
      <ThemeProvider>
        <MemoryRouter>
          <HelmetProvider>
            <NotFoundPage />
          </HelmetProvider>
        </MemoryRouter>
      </ThemeProvider>
    </Provider>,
  );

describe('<NotFoundPage />', () => {
  it('should match snapshot', () => {
    const notFoundPage = renderPage();
    expect(notFoundPage.toJSON()).toMatchSnapshot();
  });

  it('should should contain Link', () => {
    const notFoundPage = renderPage();
    expect(notFoundPage.root.findByType(Link)).toBeDefined();
  });
});
