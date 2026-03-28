/**
 *
 * Tests for ContactUsPage
 *
 * @see https://github.com/react-react-tree/master/docs/testing
 *
 */

import React from 'react';
import { render } from '@testing-library/react';
import { IntlProvider } from 'react-intl';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';

import ContactUsPage from '../index';
import { DEFAULT_LOCALE } from '../../../i18n';
import configureStore from '../../../configureStore';
describe('<ContactUsPage />', () => {
  let store;

  beforeEach(() => {
    store = configureStore({}, createMemoryHistory());
  });

  it('Expect to not log errors in console', () => {
    const spy = jest.spyOn(global.console, 'error');
    render(
      <Provider store={store}>
        <IntlProvider locale={DEFAULT_LOCALE}>
          <ContactUsPage />
        </IntlProvider>
      </Provider>,
    );
    const relevantCalls = spy.mock.calls.filter(
      (call) => !String(call[0]).includes('uses the legacy childContextTypes API'),
    );
    expect(relevantCalls).toHaveLength(0);
  });

  it.todo('Expect to have additional unit tests specified');

  /**
   * Unskip this test to use it
   *
   * @see {@link https://jestjs.io/docs/en/api#testskipname-fn}
   */
  it.skip('Should render and match the snapshot', () => {
    const {
      container: { firstChild },
    } = render(
      <Provider store={store}>
        <IntlProvider locale={DEFAULT_LOCALE}>
          <ContactUsPage />
        </IntlProvider>
      </Provider>,
    );
    expect(firstChild).toMatchSnapshot();
  });
});
