import { render } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import React from 'react';
import { HelmetProvider } from 'react-helmet-async';
import { Provider } from 'react-redux';
import { configureAppStore } from 'store/configureStore';
import { ThemeProvider } from 'styled-components';
import { themes } from 'styles/theme/themes';

import { App } from '../index';

const history = createMemoryHistory();
const store = configureAppStore({}, history);

describe('<App />', () => {
  it('should render without crashing', () => {
    const { container } = render(
      <Provider store={store}>
        <ThemeProvider theme={themes.light}>
          <HelmetProvider>
            <App />
          </HelmetProvider>
        </ThemeProvider>
      </Provider>,
    );
    expect(container).toBeTruthy();
  });
});
