import { Reducer, Store } from 'redux';
import { RouterState } from 'redux-first-history';
import { ContainerState as LanguageProviderState } from 'containers/LanguageProvider/types';
import { ContainerState as AppState } from 'containers/App/types';
import { ContainerState as HomeState } from 'containers/HomePage/types';
import { ContainerState as LoginState } from 'containers/LoginPage/types';
import { ContainerState as ChangePasswordState } from 'containers/ChangePassword/types';
import { ContainerState as PhoneVerificationPageState } from 'containers/PhoneVerificationPage/types';
import { ContainerState as AddressManagementPageState } from 'containers/AddressManagementPage/types';
import { ContainerState as AcountPageState } from 'containers/AcountPage/types';
import { ContainerState as OrdersPageState } from 'containers/OrdersPage/types';
import { ContainerState as FundsPageState } from 'containers/FundsPage/types';
import { ContainerState as DocumentVerificationPageState } from 'containers/DocumentVerificationPage/types';
import { ContainerState as ChangeUserInfoPageState } from 'containers/ChangeUserInfoPage/types';
import { ContainerState as GoogleAuthenticationPageState } from 'containers/GoogleAuthenticationPage/types';
import { ContainerState as SignupPageState } from 'containers/SignupPage/types';
import { ContainerState as RecapchaContainerState } from 'containers/RecapchaContainer/types';
import { ContainerState as EmailAuthenticationState } from 'containers/EmailVerification/types';
import { ContainerState as UpdatePasswordState } from 'containers/UpdatePasswordPage/types';
import { ContainerState as TradePageState } from 'containers/TradePage/types';
import { ContainerState as TradeChartState } from 'containers/TradePage/components/TradeChart';
import { ContainerState as TradeHeaderState } from 'containers/TradePage/components/TradeHeader';
import { ContainerState as ContactUsPageState } from 'containers/ContactUsPage/types';

export interface InjectedStore extends Store {
  injectedReducers: any;
  injectedSagas: any;
  runSaga(
    saga: (() => IterableIterator<any>) | undefined,
    args: any | undefined,
  ): any;
}

export interface InjectReducerParams {
  key: keyof ApplicationRootState;
  reducer: Reducer<any, any>;
}

export interface InjectSagaParams {
  key: keyof ApplicationRootState;
  saga: () => IterableIterator<any>;
  mode?: string | undefined;
}

// Your root reducer type, which is your redux state types also
export interface ApplicationRootState {
  readonly router: RouterState;
  readonly global: AppState;
  readonly language: LanguageProviderState;
  readonly home: HomeState;
  readonly changePassword: ChangePasswordState;
  readonly phoneVerificationPage: PhoneVerificationPageState;
  readonly AddressManagementPage: AddressManagementPageState;
  readonly loginPage: LoginState;
  readonly signupPage: SignupPageState;
  readonly acountPage: AcountPageState;
  readonly ordersPage: OrdersPageState;
  readonly tradeChart: TradeChartState;
  readonly tradePage: TradePageState;
  readonly tradeHeader: TradeHeaderState;
  readonly fundsPage: FundsPageState;
  readonly changeUserInfoPage: ChangeUserInfoPageState;
  readonly recapchaContainer: RecapchaContainerState;
  readonly documentVerificationPage: DocumentVerificationPageState;
  readonly emailAuthentication: EmailAuthenticationState;
  readonly updatePasswordPage: UpdatePasswordState;
  readonly googleAuthenticationPage: GoogleAuthenticationPageState;
  readonly contactUsPage: ContactUsPageState;
  // for testing purposes
  readonly test: any;
}
