import { WindowTypes } from 'app/constants';
import { Subject, ReplaySubject, BehaviorSubject } from 'rxjs';

enum MessageNames {
  /**
   * Sets validation errors on one or more form inputs.
   * @payload errorId: string — the input field identifier
   * @payload value: string[] — array of error message strings
   * @sender emailValidator.tsx, passwordValidator.tsx (on form validation failure)
   * @listener FormContainer index.tsx (renders inline error messages)
   */
  SET_INPUT_ERROR = 'SET_INPUT_ERROR',

  /**
   * Shows or hides the global page-level loading overlay.
   * @payload payload: boolean — true = show spinner, false = hide spinner
   * @sender apiService.ts (clears on 401 auth error)
   * @listener Global layout / app root (no explicit subscriber found; used via BehaviorSubject)
   */
  SETLOADING = 'SETLOADING',

  /**
   * Shows or hides a row-level loading indicator inside an AG Grid.
   * @payload value: boolean — true = loading, false = done
   * @payload rowId: number — the row's primary key
   * @payload gridName: GridNames — which grid instance to target
   * @sender Billing saga.ts (during approve/reject row actions)
   * @listener SimpleGrid.tsx
   */
  SET_ROW_LOADING = 'SET_ROW_LOADING',

  /**
   * Shows or hides a button-level loading spinner.
   * @payload value: boolean — true = loading, false = idle
   * @payload loadingId: string — identifies which button to target
   * @sender All action sagas (before/after API calls), loading.ts util
   * @listener isLoadingWithTextAuto.tsx
   */
  SET_BUTTON_LOADING = 'SET_BUTTON_LOADING',

  /**
   * Dismisses the reject confirmation popup in the Billing withdrawal modal.
   * @payload (none)
   * @sender Billing saga.ts (after reject API call completes)
   * @listener Billing/WithdrawModal.tsx
   */
  CLOSE_REJECT_POPUP = 'CLOSE_REJECT_POPUP',

  /**
   * Closes the currently open modal/popup dialog.
   * @payload (none)
   * @sender Balances saga.ts, CurrencyPairs saga.ts, Deposits saga.ts,
   *         FinanceMethods saga.ts, Reports saga.ts, UserDetails saga.ts,
   *         Billing saga.ts (after successful create/edit/delete operations)
   * @listener ConstructiveModal.tsx, Balances index.tsx, Billing/DepositModal.tsx,
   *           Reports/DeleteModal.tsx, Reports/EditModal.tsx, UserDetails/RejectModal.tsx
   */
  CLOSE_POPUP = 'CLOSE_POPUP',

  /**
   * Updates a single row's data in-place inside an AG Grid (avoids full refresh).
   * @payload value: object — the updated row data (must include the row's id)
   * @payload gridName?: GridNames — target grid (optional; grid uses rowId to match)
   * @sender Billing saga.ts, Deposits saga.ts (after row-level API mutations)
   * @listener SimpleGrid.tsx
   */
  UPDATE_GRID_ROW = 'UPDATE_GRID_ROW',

  /**
   * Signals a 401 Unauthorized response from the API layer.
   * @payload (none)
   * @sender apiService.ts (interceptor on every 401 response)
   * @listener App root / auth guard (triggers redirect to login page)
   */
  AUTH_ERROR_EVENT = 'AUTH_ERROR_EVENT',

  /**
   * Delivers the paginated list of user accounts to grid consumers.
   * @payload value: object — API response with user account rows and pagination meta
   * @sender UserAccounts/saga.ts (after fetchUserAccounts API call)
   * @listener HomePage index.tsx, UserAccountsPage.tsx, VerificationPage.tsx
   */
  SET_USER_ACCOUNTS = 'SET_USER_ACCOUNTS',

  /**
   * Requests that a user-details popup be opened in a new browser window.
   * @payload payload: { id: number | string } — the user or entity ID to open
   * @sender FinanceMethods/saga.ts, Reports/saga.ts (on row-action click)
   * @listener NewWindow.tsx (creates a new window and emits OPEN_WINDOW)
   */
  OPEN_NEW_WINDOW = 'OPEN_NEW_WINDOW',

  /**
   * Delivers the list of available trading currency pairs.
   * @payload value: object[] — array of currency pair objects
   * @sender MarketTicks/saga.ts (after fetchCurrencyPairs API call)
   * @listener SyncPage.tsx
   */
  SET_CURRENCY_PAIRS = 'SET_CURRENCY_PAIRS',

  /**
   * Broadcasts the current browser window inner width after a debounced resize.
   * @payload payload: number — window.innerWidth in pixels
   * @sender app/index.tsx (window.onresize handler, 150 ms debounce)
   * @listener Responsive layout components (no explicit subscriber found in current code)
   */
  RESIZE = 'RESIZE',

  /**
   * Notifies grid-adjacent filter panels that an AG Grid has been mounted or resized.
   * @payload value: { id: string; width: number } — grid instance id and its new width
   * @sender SimpleGrid.tsx (on component mount and on AG Grid resize events)
   * @listener GridFilter.tsx (repositions filter panel to align with grid)
   */
  GRID_RESIZE = 'GRID_RESIZE',

  /**
   * Displays a toast/snackbar notification to the user.
   * @payload value: string — the message text to display
   * @payload type: 'default' | 'error' | 'success' | 'warning' | 'info' — severity level
   * @payload userId?: number | string — scopes the toast to a specific user session
   * @sender Billing saga.ts, UserDetails saga.ts, WithdrawModal.tsx, Reports saga.ts
   * @listener SnackBar.tsx
   */
  TOAST = 'TOAST',

  /**
   * Shows the withdrawal confirmation popup for a specific pending withdrawal.
   * @payload userId: string — composite key of `userId + '' + withdrawalId` for targeting
   * @sender Billing saga.ts (after pre-validation), WithdrawModal.tsx (on confirm click)
   * @listener WithdrawModal.tsx (renders the confirm/reject step)
   */
  SHOW_WITHDRAW_CONFIRM_POPUP = 'SHOW_WITHDRAW_CONFIRM_POPUP',

  // ── User Details panel events ────────────────────────────────────────────

  /**
   * Delivers wallet balances for a specific user to the user-details panel.
   * @payload value: object[] — array of wallet objects (currency, balance, address, etc.)
   * @sender UserDetails/saga.ts (after fetchWallets API call)
   * @listener WaletSegment.tsx
   */
  SET_WALLETS_DATA = 'SET_WALLETS_DATA',

  /**
   * Sends raw form data from the verification info form to the document image viewer.
   * @payload value: object — submitted form field values
   * @sender VerificationWindow/LeftInfo.tsx (on form submit)
   * @listener VerificationWindow/ImageWrapper.tsx (updates displayed document data)
   */
  DATASEND = 'DATASEND',

  /**
   * Triggers a file download in the app root.
   * @payload payload: { url: string; filename: string } — download URL and suggested filename
   * @sender (triggered programmatically, e.g. export actions)
   * @listener app/index.tsx (calls downloadFile helper)
   */
  DOWNLOAD_FILE = 'DOWNLOAD_FILE',

  /**
   * Delivers the list of whitelisted withdrawal addresses for a user.
   * @payload value: object[] — array of white-address objects
   * @sender UserDetails/saga.ts (after fetchWhiteAddresses API call)
   * @listener WhiteAddressesSegemnt.tsx
   */
  SET_WHITEADDRESSES_DATA = 'SET_WHITEADDRESSES_DATA',

  /**
   * Delivers the user's asset balances to the Balances container.
   * @payload value: object[] — array of balance records per asset
   * @sender Balances/saga.ts (after fetchBalances API call)
   * @listener Balances index.tsx (renders balance grid via messageName prop)
   */
  SET_BALANCES_DATA = 'SET_BALANCES_DATA',

  /**
   * Delivers paginated balance transfer history records.
   * @payload value: object — API response with history rows and pagination meta
   * @sender Balances/saga.ts (after fetchBalancesHistory API call)
   * @listener Balances/TransferHistory.tsx (via messageName prop)
   */
  SET_BALANCES_HISTORY_DATA = 'SET_BALANCES_HISTORY_DATA',

  /**
   * Delivers the paginated platform-wide open orders for the Open Orders admin page.
   * @payload value: object — API response with open order rows and pagination meta
   * @sender OpenOrders/saga.ts (after fetchOpenOrdersPage API call)
   * @listener OpenOrders index.tsx (via messageName prop)
   */
  SET_OPEN_ORDERS_PAGE_DATA = 'SET_OPEN_ORDERS_PAGE_DATA',

  /**
   * Delivers the filled-orders list for a specific user (user-details context).
   * @payload value: object — API response with filled order rows
   * @sender UserOrders/saga.ts (after fetchFilledOrders API call for a user)
   * @listener FilledOrders index.tsx (via messageName prop)
   */
  SET_FILLED_ORDERS_DATA = 'SET_FILLED_ORDERS_DATA',

  /**
   * Delivers the list of trading currency pairs for the Currency Pairs admin page.
   * @payload value: object — API response with currency pair rows and pagination meta
   * @sender CurrencyPairs/saga.ts (after fetchCurrencyPairs API call)
   * @listener CurrencyPairs index.tsx (via messageName prop)
   */
  SET_CURRENCYPAIRS_DATA = 'SET_CURRENCYPAIRS_DATA',

  /**
   * Delivers the external exchange configuration list.
   * @payload value: object — API response with external exchange records
   * @sender ExternalExchange/saga.ts (after fetchExternalExchange API call)
   * @listener ExternalExchange index.tsx (via messageName prop)
   */
  SET_EXTERNAL_EXCHANGE_DATA = 'SET_EXTERNAL_EXCHANGE_DATA',

  /**
   * Delivers paginated deposit records for the main Deposits admin page.
   * @payload value: object — API response with deposit rows and pagination meta
   * @sender Deposits/saga.ts (after fetchDeposits API call)
   * @listener Deposits index.tsx (via messageName prop)
   */
  SET_DEPOSITS_DATA = 'SET_DEPOSITS_DATA',

  /**
   * Delivers paginated withdrawal records for the main Withdrawals admin page.
   * @payload value: object — API response with withdrawal rows and pagination meta
   * @sender Withdrawals/saga.ts (after fetchWithdrawals API call)
   * @listener Withdrawals index.tsx (via messageName prop, two grid instances)
   */
  SET_WITHDRAWALS_DATA = 'SET_WITHDRAWALS_DATA',

  /**
   * Delivers the list of finance methods (payment providers) for the Finance Methods page.
   * @payload value: object — API response with finance method rows and pagination meta
   * @sender FinanceMethods/saga.ts (after fetchFinanceMethods API call)
   * @listener FinanceMethods index.tsx (via messageName prop)
   */
  SET_FINANCEMETHODS_DATA = 'SET_FINANCEMETHODS_DATA',

  /**
   * Delivers the profile/account data for a specific user (user-details panel).
   * @payload value: object — user profile record
   * @sender UserDetails/saga.ts (after fetchUserData API call)
   * @listener UserDetailsWindow.tsx
   */
  SET_USER_DATA = 'SET_USER_DATA',

  // ── User Verification panel events ──────────────────────────────────────

  /**
   * Instructs NewWindowContainer to open and register a new popup window by ID.
   * @payload payload: { id: string; [key: string]: unknown } — unique window ID and initial data
   * @sender NewWindow.tsx, BillingWithdrawalssDataGrid.tsx, useOpenWithdrawWindow.tsx
   * @listener NewWindowContainer.tsx (creates or focuses the window entry)
   */
  OPEN_WINDOW = 'OPEN_WINDOW',

  /**
   * Instructs NewWindowContainer to close and unregister a popup window by ID.
   * @payload payload: { id: string } — the window ID to close
   * @sender Billing/WithdrawModal.tsx (on modal close action)
   * @listener NewWindowContainer.tsx
   */
  CLOSE_WINDOW = 'CLOSE_WINDOW',

  /**
   * Delivers the detailed record for a withdrawal item to any interested panel.
   * @payload value: object — full withdrawal detail record including userId and id
   * @sender Billing/saga.ts, Withdrawals/saga.ts, FinanceMethods/saga.ts
   * @listener Billing/BillingWithdrawalssDataGrid.tsx (triggers OPEN_WINDOW),
   *           useOpenWithdrawWindow.tsx (opens withdrawal detail window)
   */
  SET_MAIN_WITHDRAWALS_ITEM_DETAILS = 'SET_MAIN_WITHDRAWALS_ITEM_DETAILS',

  /**
   * Triggers a full data reload of the targeted AG Grid.
   * @payload gridName?: GridNames — which grid to refresh (omit to refresh all)
   * @sender Billing/saga.ts, CurrencyPairs/saga.ts, FinanceMethods/saga.ts,
   *         UserDetails/saga.ts (after successful create/edit/delete mutations)
   * @listener SimpleGrid.tsx (calls the grid's datasource refresh)
   */
  REFRESH_GRID = 'REFRESH_GRID',

  /**
   * Applies saved filter/sort parameters to a specific AG Grid instance.
   * @payload value: object — AG Grid column-state or filter-model snapshot
   * @payload gridName?: GridNames — target grid instance
   * @sender Withdrawals/index.tsx (restores persisted grid state)
   * @listener SimpleGrid.tsx (calls grid API to apply the params)
   */
  APPLY_PARAMS_TO_GRID = 'APPLY_PARAMS_TO_GRID',
}
enum GridNames {
  Billing_Withdraw = 'Billing_Withdraw',
  BILLING_DEPOSIT = 'BILLING_DEPOSIT',
  DEPOSITS_PAGE = 'DEPOSITS_PAGE',
  FINANCE_METHODS_PAGE = 'FINANCE_METHODS_PAGE',
  CURRENCY_PAIRS_PAGE = 'CURRENCY_PAIRS_PAGE',
  MAIN_WITHDRAWALS = 'MAIN_WITHDRAWALS',
  USER_VERIFICATION = 'USER_VERIFICATION',
}

/**
 * Standard message shape for all MessageService pub/sub events.
 *
 * @property name - The event type from the MessageNames enum
 * @property value - Primary payload (type varies by event; see MessageNames JSDoc)
 * @property payload - Alternative payload field (some events use this instead of value)
 * @property additional - Extra data for complex events requiring secondary context
 * @property errorId - Identifies which input field carries a validation error
 * @property userId - User context for user-scoped events (e.g. TOAST, SHOW_WITHDRAW_CONFIRM_POPUP)
 * @property loadingId - Identifies which loading indicator (button/spinner) to control
 * @property rowId - Grid row primary key for row-level operations (SET_ROW_LOADING, UPDATE_GRID_ROW)
 * @property gridName - Which AG Grid instance to target (from GridNames enum)
 * @property type - Window type for OPEN_WINDOW, or toast severity for TOAST
 * @property child - Child/nested data for events that carry hierarchical payloads
 */
export interface BroadcastMessage {
  name: MessageNames;
  value?: unknown;
  payload?: unknown;
  additional?: unknown;
  errorId?: string;
  userId?: number | string;
  loadingId?: string;
  rowId?: number;
  gridName?: string;
  type?: WindowTypes | 'default' | 'error' | 'success' | 'warning' | 'info';
  child?: unknown;
}
let value: Record<string, unknown> = {};
const Subscriber = new Subject<BroadcastMessage>();
const MessageService = {
	send: function (msg: BroadcastMessage) {
		Subscriber.next(msg);
	},
};

const RepaySubscriber3 = new ReplaySubject(3);
const ReplayMessageService3 = {
	send: function (msg: BroadcastMessage) {
		RepaySubscriber3.next(msg);
	},
};
const BehaviorSubscriber: BehaviorSubject<BroadcastMessage> = new BehaviorSubject({} as BroadcastMessage);
const BehaviorMessageService = {
	send: function (msg: BroadcastMessage) {
		BehaviorSubscriber.next(msg);
	},
};

export {
	MessageNames,
	MessageService,
	Subscriber,
	ReplayMessageService3,
	RepaySubscriber3,
	BehaviorMessageService,
	BehaviorSubscriber,
	GridNames,
};
