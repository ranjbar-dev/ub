import React from 'react';
import { screen } from '@testing-library/react';
import { renderWithProviders } from 'utils/testUtils';

// Mock the saga
jest.mock('../saga', () => ({
  depositsSaga: function* () {},
}));

// Mock SimpleGrid to avoid AG Grid complexity
jest.mock('app/components/SimpleGrid/SimpleGrid', () => ({
  SimpleGrid: ({ containerId }: { containerId?: string }) => (
    <div data-testid={`simple-grid-${containerId || 'default'}`}>SimpleGrid</div>
  ),
}));

jest.mock('app/components/titledContainer/TitledContainer', () => ({
  __esModule: true,
  default: ({ children, title }: { children: React.ReactNode; title: string }) => (
    <div data-testid="titled-container">
      <div data-testid="container-title">{title}</div>
      {children}
    </div>
  ),
}));

jest.mock('app/components/wrappers/FullWidthWrapper', () => ({
  FullWidthWrapper: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="full-width-wrapper">{children}</div>
  ),
}));

jest.mock('app/components/materialModal/modal', () => ({
  __esModule: true,
  default: ({ children, isOpen }: { children: React.ReactNode; isOpen: boolean }) =>
    isOpen ? <div data-testid="popup-modal">{children}</div> : null,
}));

jest.mock('../../Billing/components/DepositModal', () => ({
  __esModule: true,
  default: () => <div data-testid="deposit-modal">DepositModal</div>,
}));

// Properly mock localStorage with bracket-access support
beforeAll(() => {
  localStorage.setItem('currencies', JSON.stringify({
    currencies: [
      { name: 'Bitcoin', code: 'BTC' },
      { name: 'Ethereum', code: 'ETH' },
    ],
  }));
});

afterAll(() => {
  localStorage.clear();
});

describe('<Deposits />', () => {
  // Lazy import to ensure mocks are in place
  let Deposits: React.ComponentType;

  beforeAll(async () => {
    const mod = await import('../index');
    Deposits = mod.Deposits;
  });

  it('should render without crashing', () => {
    renderWithProviders(<Deposits />);
    expect(screen.getByTestId('full-width-wrapper')).toBeInTheDocument();
  });

  it('should display titled container', () => {
    renderWithProviders(<Deposits />);
    expect(screen.getByTestId('titled-container')).toBeInTheDocument();
  });
});