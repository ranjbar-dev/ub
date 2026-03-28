/**
 * Shared test utilities for ub-admin rendering and mocking.
 * Import from 'utils/testUtils' in any test file.
 */
import React, { ReactElement } from 'react';
import { render, RenderOptions, RenderResult } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { ThemeProvider } from 'styled-components';
import { HelmetProvider } from 'react-helmet-async';
import { Router } from 'react-router-dom';
import { createMemoryHistory, MemoryHistory } from 'history';
import { themes } from 'styles/theme/themes';
import { configureAppStore } from 'store/configureStore';

// Re-export everything from @testing-library/react
export * from '@testing-library/react';

interface RenderWithProvidersOptions extends Omit<RenderOptions, 'wrapper'> {
  preloadedState?: Record<string, unknown>;
  history?: MemoryHistory;
}

/**
 * Renders a component wrapped with all required providers:
 * Redux Provider, ThemeProvider, HelmetProvider, Router.
 */
export function renderWithProviders(
  ui: ReactElement,
  options: RenderWithProvidersOptions = {},
) {
  const { preloadedState = {}, history = createMemoryHistory(), ...renderOptions } = options;

  // Use the same store configuration as the app
  const store = configureAppStore(preloadedState, history);

  function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <Provider store={store}>
        <ThemeProvider theme={themes.dark}>
          <HelmetProvider>
            {/* cast needed: history v4 MemoryHistory vs react-router-dom v5 History types */}
            {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
            <Router history={history as any}>{children}</Router>
          </HelmetProvider>
        </ThemeProvider>
      </Provider>
    );
  }

  return {
    ...render(ui, { wrapper: Wrapper, ...renderOptions }),
    store,
    history,
  };
}

/** Creates a mock fetch that resolves with the given response body and status. */
export function createMockFetch(body: unknown, status = 200): jest.Mock {
  const mockFn = jest.fn(() =>
    Promise.resolve(
      new Response(JSON.stringify(body), {
        status,
        headers: { 'Content-type': 'application/json' },
      }),
    ),
  );
  return mockFn;
}

/** Creates a StandardResponse-shaped object for testing. */
export function createMockApiResponse<T>(data: T, status = true, message = '') {
  return { status, message, data };
}

/** Mock localStorage with jest spies. Call in beforeEach, returns cleanup fn. */
export function mockLocalStorage(initialData: Record<string, string> = {}) {
  const store: Record<string, string> = { ...initialData };
  const getItemSpy = jest.spyOn(Storage.prototype, 'getItem').mockImplementation((key: string) => store[key] ?? null);
  const setItemSpy = jest.spyOn(Storage.prototype, 'setItem').mockImplementation((key: string, value: string) => { store[key] = value; });
  const removeItemSpy = jest.spyOn(Storage.prototype, 'removeItem').mockImplementation((key: string) => { delete store[key]; });
  const clearSpy = jest.spyOn(Storage.prototype, 'clear').mockImplementation(() => { Object.keys(store).forEach(k => delete store[k]); });

  return {
    store,
    getItemSpy,
    setItemSpy,
    removeItemSpy,
    clearSpy,
    restore: () => {
      getItemSpy.mockRestore();
      setItemSpy.mockRestore();
      removeItemSpy.mockRestore();
      clearSpy.mockRestore();
    },
  };
}
