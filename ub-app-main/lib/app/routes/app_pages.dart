import 'package:get/get.dart';
import 'package:unitedbit/app/modules/funds/pages/autoExchange/views/auto_exchange_view.dart';
import '../modules/exchange/views/widgets/exchange_search.dart';

import '../modules/account/bindings/account_binding.dart';
import '../modules/account/views/account_view.dart';
import '../modules/addNewAddress/bindings/add_new_address_binding.dart';
import '../modules/addNewAddress/views/add_new_address_view.dart';
import '../modules/afterSplash/bindings/after_splash_binding.dart';
import '../modules/afterSplash/views/after_splash_view.dart';
import '../modules/changePassword/bindings/change_password_binding.dart';
import '../modules/changePassword/views/change_password_view.dart';
import '../modules/checkYourEmail/bindings/check_your_email_binding.dart';
import '../modules/checkYourEmail/views/check_your_email_view.dart';
import '../modules/exchange/bindings/exchange_binding.dart';
import '../modules/exchange/views/exchange_view.dart';
import '../modules/forgot/bindings/forgot_binding.dart';
import '../modules/forgot/views/forgot_view.dart';
import '../modules/funds/bindings/funds_binding.dart';
import '../modules/funds/pages/autoExchange/bindings/auto_exchange_binding.dart';
import '../modules/funds/pages/balance/bindings/balance_binding.dart';
import '../modules/funds/pages/balance/views/balance_view.dart';
import '../modules/funds/pages/deposits/bindings/deposits_binding.dart';
import '../modules/funds/pages/deposits/views/depositDetails.dart';
import '../modules/funds/pages/deposits/views/deposits_view.dart';
import '../modules/funds/pages/transactionHistory/bindings/transaction_history_binding.dart';
import '../modules/funds/pages/transactionHistory/views/transaction_history_view.dart';
import '../modules/funds/pages/withdrawals/bindings/withdrawals_binding.dart';
import '../modules/funds/pages/withdrawals/views/withdrawals_view.dart';
import '../modules/funds/views/funds_view.dart';
import '../modules/home/bindings/home_binding.dart';
import '../modules/home/views/home_view.dart';
import '../modules/identityDocuments/bindings/identity_documents_binding.dart';
import '../modules/identityDocuments/views/identity_documents_view.dart';
import '../modules/identityInfo/bindings/identity_info_binding.dart';
import '../modules/identityInfo/views/identity_info_view.dart';
import '../modules/landing/bindings/landing_binding.dart';
import '../modules/landing/views/landing_view.dart';
import '../modules/login/bindings/login_binding.dart';
import '../modules/login/views/login_view.dart';
import '../modules/market/bindings/market_binding.dart';
import '../modules/market/views/edit_favorites_view.dart';
import '../modules/market/views/market_view.dart';
import '../modules/orders/bindings/orders_binding.dart';
import '../modules/orders/pages/openOrders/bindings/open_orders_binding.dart';
import '../modules/orders/pages/openOrders/views/open_orders_view.dart';
import '../modules/orders/pages/orderHistory/bindings/order_history_binding.dart';
import '../modules/orders/pages/orderHistory/views/order_history_view.dart';
import '../modules/orders/views/orders_view.dart';
import '../modules/phoneVerification/bindings/phone_verification_binding.dart';
import '../modules/phoneVerification/views/phone_verification_view.dart';
import '../modules/qrScan/bindings/qr_scan_binding.dart';
import '../modules/qrScan/views/qr_scan_view.dart';
import '../modules/separateMessagePage/bindings/separate_message_page_binding.dart';
import '../modules/separateMessagePage/views/separate_message_page_view.dart';
import '../modules/signup/bindings/signup_binding.dart';
import '../modules/signup/views/signup_view.dart';
import '../modules/trade/bindings/trade_binding.dart';
import '../modules/trade/views/trade_view.dart';
import '../modules/twoFactorAuthentication/bindings/two_factor_authentication_binding.dart';
import '../modules/twoFactorAuthentication/views/two_factor_authentication_view.dart';
import '../modules/webViewPage/bindings/web_view_page_binding.dart';
import '../modules/webViewPage/views/web_view_page_view.dart';
import '../modules/withdrawAddressManagement/bindings/withdraw_address_management_binding.dart';
import '../modules/withdrawAddressManagement/views/withdraw_address_management_view.dart';

part 'app_routes.dart';

class AppPages {
  static const INITIAL = Routes.AFTER_SPLASH;
  static const LOGIN = Routes.LOGIN;
  static const HOME = Routes.HOME;
  static const SIGNUP = Routes.SIGNUP;
  static const FORGOT = Routes.FORGOT;
  static const FUNDS = Routes.FUNDS;
  static const EXCHANGE = Routes.EXCHANGE;
  static const DEPOSITDETAILS = Routes.DEPOST_DETAILS;
  static const ACCOUNT = Routes.ACCOUNT;
  static const WITHDRAWADDRESSMANAGEMENT = Routes.WITHDRAW_ADDRESS_MANAGEMENT;
  static const TRADE = Routes.TRADE;
  static const IDENTITYDOCUMENTS = Routes.IDENTITY_DOCUMENTS;
  static const IDENTITYVERIFICATION = Routes.IDENTITY_VERIFICATION;
  static const TWOFACTORAUTHENTICATION = Routes.TWO_FACTOR_AUTHENTICATION;
  static const CHANGEPASSWORD = Routes.CHANGE_PASSWORD;
  static const MARKET = Routes.MARKET;
  static const LANDING = Routes.LANDING;
  static const ADD_NEW_ADDRESS = Routes.ADD_NEW_ADDRESS;
  static const QR_SCAN = Routes.QR_SCAN;
  static const PHONE_VERIFICATION = Routes.PHONE_VERIFICATION;
  static const ORDERS = Routes.ORDERS;
  static const WITHDRAWALS = Routes.WITHDRAWALS;
  static const DEPOSITS = Routes.DEPOSITS;
  static const EDIT_FAVORITES = Routes.EDIT_FAVORITES;
  static const EXCHANGE_SEARCH = Routes.EXCHANGE_SEARCH;
  static const AUTO_EXCGANGE = Routes.AUTO_EXCHANGE;

  static final routes = [
    GetPage(
      name: _Paths.HOME,
      page: () => HomeView(),
      binding: HomeBinding(),
    ),
    GetPage(
      name: _Paths.LOGIN,
      page: () => LoginView(),
      binding: LoginBinding(),
    ),
    GetPage(
      name: _Paths.LANDING,
      page: () => LandingView(),
      binding: LandingBinding(),
    ),
    GetPage(
      name: _Paths.SIGNUP,
      page: () => SignupView(),
      binding: SignupBinding(),
    ),
    GetPage(
      name: _Paths.FORGOT,
      page: () => ForgotView(),
      binding: ForgotBinding(),
    ),
    GetPage(
      name: _Paths.ACCOUNT,
      page: () => AccountView(),
      binding: AccountBinding(),
    ),
    GetPage(
      name: _Paths.TRADE,
      page: () => TradeView(),
      binding: TradeBinding(),
    ),
    GetPage(
      name: _Paths.OPEN_ORDERS,
      page: () => OpenOrdersView(),
      binding: OpenOrdersBinding(),
    ),
    GetPage(
      name: _Paths.ORDERS,
      page: () => OrdersView(),
      binding: OrdersBinding(),
    ),
    GetPage(
      name: _Paths.ORDER_HISTORY,
      page: () => OrderHistoryView(),
      binding: OrderHistoryBinding(),
    ),
    GetPage(
        name: _Paths.FUNDS,
        page: () => FundsView(),
        bindings: [FundsBinding(), BalanceBinding()]),
    GetPage(
      name: _Paths.BALANCE,
      page: () => BalanceView(),
      binding: BalanceBinding(),
    ),
    GetPage(
      name: _Paths.DEPOSITS,
      page: () => DepositsView(),
      binding: DepositsBinding(),
    ),
    GetPage(
      name: _Paths.DEPOST_DETAILS,
      page: () => DepostDetailsView(),
      binding: DepositsBinding(),
    ),
    GetPage(
      name: _Paths.WITHDRAWALS,
      page: () => WithdrawalsView(),
      binding: WithdrawalsBinding(),
    ),
    GetPage(
      name: _Paths.TRANSACTION_HISTORY,
      page: () => TransactionHistoryView(),
      binding: TransactionHistoryBinding(),
    ),
    GetPage(
      name: _Paths.MARKET,
      page: () => MarketView(),
      binding: MarketBinding(),
    ),
    GetPage(
      name: _Paths.CHANGE_PASSWORD,
      page: () => ChangePasswordView(),
      binding: ChangePasswordBinding(),
    ),
    GetPage(
      name: _Paths.WITHDRAW_ADDRESS_MANAGEMENT,
      page: () => WithdrawAddressManagementView(),
      binding: WithdrawAddressManagementBinding(),
    ),
    GetPage(
      name: _Paths.TWO_FACTOR_AUTHENTICATION,
      page: () => TwoFactorAuthenticationView(),
      binding: TwoFactorAuthenticationBinding(),
    ),
    GetPage(
      name: _Paths.IDENTITY_VERIFICATION,
      page: () => IdentityInfoView(),
      binding: IdentityInfoBinding(),
    ),
    GetPage(
      name: _Paths.IDENTITY_DOCUMENTS,
      page: () => IdentityDocumentsView(),
      binding: IdentityDocumentsBinding(),
    ),
    GetPage(
      name: _Paths.PHONE_VERIFICATION,
      page: () => PhoneVerificationView(),
      binding: PhoneVerificationBinding(),
    ),
    GetPage(
      name: _Paths.ADD_NEW_ADDRESS,
      page: () => AddNewAddressView(),
      binding: AddNewAddressBinding(),
    ),
    GetPage(
      name: _Paths.QR_SCAN,
      page: () => QrScanView(),
      binding: QrScanBinding(),
    ),
    GetPage(
      name: _Paths.EDIT_FAVORITES,
      page: () => EditFavoritesView(),
      binding: MarketBinding(),
    ),
    GetPage(
      name: _Paths.AFTER_SPLASH,
      page: () => AfterSplashView(),
      binding: AfterSplashBinding(),
    ),
    GetPage(
      name: _Paths.WEB_VIEW_PAGE,
      page: () => WebViewPageView(),
      binding: WebViewPageBinding(),
    ),
    GetPage(
      name: _Paths.CHECK_YOUR_EMAIL,
      page: () => CheckYourEmailView(),
      binding: CheckYourEmailBinding(),
    ),
    GetPage(
      name: _Paths.SEPARATE_MESSAGE_PAGE,
      page: () => SeparateMessagePageView(image: ''),
      binding: SeparateMessagePageBinding(),
    ),
    GetPage(
      name: _Paths.EXCHANGE,
      page: () => ExchangeView(),
      binding: ExchangeBinding(),
    ),
    GetPage(
      name: _Paths.EXCHANGE_SEARCH,
      page: () => ExchangeSearch(),
      binding: ExchangeBinding(),
    ),
    GetPage(
      name: _Paths.AUTO_EXCHANGE,
      page: () => AutoExchangeView(),
      binding: AutoExchangeBinding(),
    ),
  ];
}
