/**
 * Testing the NotFoundPage
 */

import React from 'react';
import { render } from '@testing-library/react';
import { IntlProvider } from 'react-intl';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';

import NotFound from '../index';
import messages from '../messages';
import configureStore from '../../../configureStore';

describe('<NotFound />', () => {
  it('should render the Page Not Found text', () => {
    const store = configureStore({}, createMemoryHistory());
    const { queryByText } = render(
      // tslint:disable-next-line: jsx-wrap-multiline
      <Provider store={store}>
        <IntlProvider locale="en">
          <NotFound />
        </IntlProvider>
      </Provider>,
    );
    expect(queryByText(messages.Pagenotfound.defaultMessage)).toBeInTheDocument();
  });
});
