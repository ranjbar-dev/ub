import { render } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import React from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import { configureAppStore } from 'store/configureStore';
import { ThemeProvider } from 'styles/theme/ThemeProvider';

import { NavBar } from '../index';

const history = createMemoryHistory();
const store = configureAppStore({}, history);

describe('<NavBar />', () => {
  it('should render without crashing', () => {
    const { container } = render(
      <Provider store={store}>
        <ThemeProvider>
          <MemoryRouter>
            <NavBar />
          </MemoryRouter>
        </ThemeProvider>
      </Provider>,
    );
    expect(container.firstChild).toBeTruthy();
  });
});
