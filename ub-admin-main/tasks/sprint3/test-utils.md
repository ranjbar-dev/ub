# Task: Create Shared Test Utilities

## Goal
Create `src/utils/testUtils.tsx` with reusable test helpers so every test file doesn't reinvent the wheel.

## Context
- 23 existing tests each create their own helpers (e.g., `renderWithTheme()`)
- No centralized mock for fetch, localStorage, or Redux store wrapping
- @testing-library/react, jest-styled-components, react-test-renderer all available
- setupTests.ts already imports @testing-library/jest-dom/extend-expect

## File to Create: `src/utils/testUtils.tsx`

```tsx
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
import { globalReducer } from 'store/slice';

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
): RenderResult & { store: ReturnType<typeof configureStore>; history: MemoryHistory } {
  const { preloadedState = {}, history = createMemoryHistory(), ...renderOptions } = options;

  const store = configureStore({
    reducer: {
      global: globalReducer,
    },
    preloadedState,
  });

  function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <Provider store={store}>
        <ThemeProvider theme={themes.dark}>
          <HelmetProvider>
            <Router history={history}>{children}</Router>
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
```

## Validation
- File compiles: `npx tsc --noEmit src/utils/testUtils.tsx` (or just build)
- Imports resolve correctly (themes, globalReducer, etc.)

## IMPORTANT
- Check the actual exports from `store/slice.ts` — the reducer might be exported as `globalSlice.reducer` or `globalReducer` — use whatever the file actually exports
- Check `styles/theme/themes` for the actual theme object shape
- Adapt imports to match actual file exports
