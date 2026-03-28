import React from 'react';
import { screen } from '@testing-library/react';
import { renderWithProviders } from 'utils/testUtils';
import { Billing } from '../index';
import { InitialUserDetails } from '../../UserAccounts/types';

jest.mock('../saga', () => ({
  billingSaga: function* () {},
}));

jest.mock('app/containers/UserDetails/components', () => ({
  __esModule: true,
  default: ({ options }: { options: Array<{ title: string; component: React.ReactNode }> }) => (
    <div data-testid="user-details-tabs">
      {options.map((option, index) => (
        <div key={index} data-testid={`tab-${index}`}>
          <span data-testid={`tab-title-${index}`}>{option.title}</span>
          <div data-testid={`tab-content-${index}`}>{option.component}</div>
        </div>
      ))}
    </div>
  ),
}));

jest.mock('../components/BillingDepositsDataGrid', () => ({
  __esModule: true,
  default: () => <div data-testid="billing-deposits-grid">BillingDepositsDataGrid</div>,
}));

jest.mock('../components/BillingWithdrawalssDataGrid', () => ({
  __esModule: true,
  default: () => <div data-testid="billing-withdrawals-grid">BillingWithdrawalssDataGrid</div>,
}));

jest.mock('../components/BillingAllTransactionsDataGrid', () => ({
  __esModule: true,
  default: () => <div data-testid="billing-all-transactions-grid">BillingAllTransactionsDataGrid</div>,
}));

// Properly mock localStorage
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

const mockInitialData: InitialUserDetails = {
  id: 1,
  email: 'test@example.com',
  fullName: 'John Doe',
  accountStatus: 'active',
  address: '123 Main St',
  addressConfirmationStatus: 'confirmed',
  birthDate: '1990-01-01',
  city: 'New York',
  country: 'USA',
  countryId: 1,
  gender: 'male',
  groupId: '1',
  groupName: 'Standard',
  identityConfirmationStatus: 'confirmed',
  managerId: '1',
  managerName: 'Admin',
  metaData: {
    userLevels: [],
    userGroups: [],
    userStatuses: [],
    userProfileStatuses: [],
    countries: [],
  },
  mobile: '+1234567890',
  phoneConfirmationStatus: 'confirmed',
  postalCode: '10001',
  profileStatus: 'completed',
  referKey: 'ABC123',
  referralId: '456',
  registeredIp: '192.168.1.1',
  registrationDate: '2023-01-01',
  status: 'active',
  systemId: 1,
  totalBalance: '1000.00',
  totalCommissions: '50.00',
  totalDeposit: '2000.00',
  totalOnTrade: '500.00',
  totalWithdraw: '1000.00',
  trustLevel: 3,
  userLevelId: 1,
  userLevelName: 'Basic',
};

describe('<Billing />', () => {
  it('should render without crashing', () => {
    renderWithProviders(<Billing initialData={mockInitialData} />);
    expect(screen.getByTestId('user-details-tabs')).toBeInTheDocument();
  });

  it('should render tab layout with all three tabs', () => {
    renderWithProviders(<Billing initialData={mockInitialData} />);
    expect(screen.getByTestId('tab-0')).toBeInTheDocument();
    expect(screen.getByTestId('tab-1')).toBeInTheDocument();
    expect(screen.getByTestId('tab-2')).toBeInTheDocument();
  });
});