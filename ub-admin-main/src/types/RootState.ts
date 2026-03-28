import { AdminsState } from 'app/containers/Admins/types';
import { BalancesState } from 'app/containers/Balances/types';
import { BillingState } from 'app/containers/Billing/types';
import { CurrencyPairsState } from 'app/containers/CurrencyPairs/types';
import { DepositsState } from 'app/containers/Deposits/types';
import { ExternalExchangeState } from 'app/containers/ExternalExchange/types';
import { ExternalOrdersState } from 'app/containers/ExternalOrders/types';
import { FilledOrdersState } from 'app/containers/FilledOrders/types';
import { FinanceMethodsState } from 'app/containers/FinanceMethods/types';
import { HomePageState } from 'app/containers/HomePage/types';
import { LiquidityOrdersState } from 'app/containers/LiquidityOrders/types';
import { LoginHistoryState } from 'app/containers/LoginHistory/types';
import { LoginPageState } from 'app/containers/LoginPage/types';
import { MarketTicksState } from 'app/containers/MarketTicks/types';
import { OpenOrdersState } from 'app/containers/OpenOrders/types';
import { OrdersState } from 'app/containers/Orders/types';
import { ReportsState } from 'app/containers/Reports/types';
import { ScanBlockState } from 'app/containers/ScanBlock/types';
import { UserAccountsState } from 'app/containers/UserAccounts/types';
import { UserDetailsState } from 'app/containers/UserDetails/types';
import { VerificationWindowState } from 'app/containers/VerificationWindow/types';
import { WithdrawalsState } from 'app/containers/Withdrawals/types';
import { RouterState } from 'connected-react-router';
import { GlobalState } from 'store/slice';
import { ThemeState } from 'styles/theme/types';
// [IMPORT NEW CONTAINERSTATE ABOVE] < Needed for generating containers seamlessly

/* 
  Because the redux-injectors injects your reducers asynchronously somewhere in your code
  You have to declare them here manually
  Properties are optional because they are injected when the components are mounted sometime in your application's life. 
  So, not available always
*/
export interface RootState {
  theme?: ThemeState;
  loginPage?: LoginPageState;
  global?: GlobalState;
  router?: RouterState;
  userAccounts?: UserAccountsState;
  userDetails?: UserDetailsState;
  billing?: BillingState;
  reports?: ReportsState;
  orders?: OrdersState;
  verificationWindow?: VerificationWindowState;
  loginHistory?: LoginHistoryState;
  openOrders?: OpenOrdersState;
  filledOrders?: FilledOrdersState;
  externalOrders?: ExternalOrdersState;
  deposits?: DepositsState;
  withdrawals?: WithdrawalsState;
  financeMethods?: FinanceMethodsState;
  currencyPairs?: CurrencyPairsState;
  externalExchange?: ExternalExchangeState;
  marketTicks?: MarketTicksState;
  admins?: AdminsState;
  homePage?: HomePageState;
  balances?: BalancesState;
  scanBlock?: ScanBlockState;
  liquidityOrders?: LiquidityOrdersState;
  // [INSERT NEW REDUCER KEY ABOVE] < Needed for generating containers seamlessly
}
