import React from 'react';
import { screen } from '@testing-library/react';
import { renderWithProviders } from 'utils/testUtils';

jest.mock('../saga', () => ({
  userAccountsSaga: function* () {},
}));

jest.mock('../components/UserAccountsPage', () => ({
  __esModule: true,
  default: () => <div data-testid="user-accounts-page">UserAccountsPage</div>,
}));

jest.mock('../components/VerificationPage', () => ({
  __esModule: true,
  default: () => <div data-testid="verification-page">VerificationPage</div>,
}));

// Mock the selector to avoid connected-react-router dependency
jest.mock('../selectors', () => ({
  selectRouter: () => ({ location: { pathname: '/user-accounts' } }),
  selectUserAccountsData: () => null,
}));

describe('<UserAccounts />', () => {
  let UserAccounts: React.ComponentType;

  beforeAll(async () => {
    const mod = await import('../index');
    UserAccounts = mod.UserAccounts;
  });

  it('should render without crashing', () => {
    renderWithProviders(<UserAccounts />);
    expect(screen.getByTestId('user-accounts-page')).toBeInTheDocument();
  });

  it('should render UserAccountsPage when not on verification route', () => {
    renderWithProviders(<UserAccounts />);
    expect(screen.getByTestId('user-accounts-page')).toBeInTheDocument();
  });
});