import ActionTypes from './constants';
import { takeLatest, put, call, takeEvery, all } from 'redux-saga/effects';
import {
  setIsLoadingBalancePageDataAction,
  setBalancePageDataAction,
  setIsLoadingdepositAndWithDrawDataAction,
  setDepositAndWithDrawDataAction,
  setIsLoadingTransactionHistoryPageDataAction,
  setTransactionHistoryPageDataAction,
  setFormerWithdrawAddressesAction,
  addWithdrawDataAction,
  getInfiniteDWAction,
  addFormerWithdrawAddressesAction,
  withdrawAction,
  setUserDataAction,
} from './actions';
import { toast } from 'components/Customized/react-toastify';
import { StandardResponse } from 'services/constants';
import {
  getBalancesAPI,
  getDepositAndWithdrawAPI,
  getTransactionHistoryAPI,
  getOrderDetailAPI,
  getFormerWithdrawAddressesAPI,
  withdrawAPI,
  preWithdrawAPI,
} from 'services/funds_services';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { WithdrawModel, InfiniteDwModel } from './types';
import { ToastMessages } from 'services/toastService';
import { addNewWithDrawAddressAPI } from 'services/address_management_service';
import { FilterModel } from 'containers/OrdersPage/types';
import { addOneToWithdrawAddressesAction } from 'containers/AddressManagementPage/actions';
import { BalanceArrayFormatter } from './utils';
import { getUserDataAPI } from 'services/security_service';

export const infinitePageSize = window.innerHeight > 1100 ? 60 : 30;

function* getBalances(action: {
  type: string;
  payload: { isSilent?: boolean };
}) {
  if (!action.payload.isSilent) {
    yield put(setIsLoadingBalancePageDataAction(true));
  }
  try {
    const response: StandardResponse = yield call(getBalancesAPI);
    if (response.status === false) {
      yield put(setIsLoadingBalancePageDataAction(false));
      toast.error('error getting balances');
      return;
    }
    const BalanceArray = BalanceArrayFormatter(response.data.balances);
    const data = { balances: BalanceArray, ...response.data };
    yield put(setBalancePageDataAction(data));

    if (action.payload.isSilent) {
      MessageService.send({
        name: MessageNames.SET_BALANCE_PAGE_DATA,
        payload: BalanceArray,
      });
    }
  } catch (error) {
    yield put(setIsLoadingBalancePageDataAction(false));
    toast.error('Error while getting balances');
  }
}
function* getDepositAndWithdraws(action: {
  type: string;
  payload: { code: string; type: string };
}) {
  yield put(setIsLoadingdepositAndWithDrawDataAction(true));
  try {
    const response: StandardResponse = yield call(
      getDepositAndWithdrawAPI,
      action.payload.code,
    );
    if (response.status === false) {
      yield put(setIsLoadingdepositAndWithDrawDataAction(false));
      toast.error('error getting deposit and withdraw data');
      return;
    } else if (response.status === true) {
      yield put(setDepositAndWithDrawDataAction(response.data));
      yield put(
        getInfiniteDWAction({
          // code: action.payload.code,
          type: action.payload.type,
          page: 0,
          page_size: infinitePageSize,
        }),
      );
    }
  } catch (error) {
    yield put(setIsLoadingdepositAndWithDrawDataAction(false));
    toast.error('error while getting data');
  }
}
function* getRawDepositAndWithdraws(action: {
  type: string;
  payload: {
    code: string;
    type: string;
    fromCoinChange?: boolean;
    silent?: boolean;
  };
}) {
  const { payload } = action;
  const { fromCoinChange, silent, code, type } = payload;
  Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: true });
  if (fromCoinChange !== true) {
    Subscriber.next({ name: MessageNames.RESET_INFINITE_SCROLL });
    if (!silent) {
      MessageService.send({
        name: MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING,
        payload: true,
      });
    }
  }

  try {
    const response: StandardResponse = yield call(
      getDepositAndWithdrawAPI,
      code,
    );
    if (response.status === false) {
      Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });
      toast.error('error getting deposit and withdraw data');
      return;
    } else if (response.status === true) {
      Subscriber.next({
        name: MessageNames.SET_DEPOSIT_PAGE_DATA,
        payload: response.data,
      });
      yield all([
        setDepositAndWithDrawDataAction(response.data),
        ...(fromCoinChange !== true
          ? [
              getInfiniteDWAction({
                // code: code,
                type: type,
                page: 0,
                page_size: infinitePageSize,
                silent: silent,
              }),
            ]
          : []),
      ]);

      Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });
    }
  } catch (error) {
    Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });

    toast.error('error while getting data');
  }
}
function* getInfiniteDw(action: { type: string; payload: InfiniteDwModel }) {
  if (action.payload.page === 0 && action.payload.silent !== true) {
    MessageService.send({
      name: MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING,
      payload: true,
    });
  }
  const response: StandardResponse = yield call(
    getTransactionHistoryAPI,
    action.payload,
  );
  if (action.payload.page === 0) {
    MessageService.send({
      name: MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING,
      payload: false,
    });
    MessageService.send({
      name: MessageNames.SET_INITIAL_INFINITE_DW_PAGE_DATA,
      payload: response.data.payments,
    });
  } else {
    MessageService.send({
      name: MessageNames.SET_DATA_TO_INFINITE_BOTTOM,
      payload: response.data.payments,
    });
  }
}

function* getFormerWithdrawAddresses(action: {
  type: string;
  payload: { code: string };
}) {
  Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: true });
  try {
    const response: StandardResponse = yield call(
      getFormerWithdrawAddressesAPI,
      action.payload.code,
    );
    if (response.status === false) {
      Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });
      toast.error('error getting deposit and withdraw data');
      return;
    }
    yield put(setFormerWithdrawAddressesAction(response.data));
    Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });
    Subscriber.next({
      name: MessageNames.SET_FORMER_WITHDRAW_ADDRESSES,
      payload: response.data,
    });
  } catch (error) {
    Subscriber.next({ name: MessageNames.SET_PAGE_LOADING, payload: false });
    toast.error('error while getting data');
  }
}

function* getTransactionHistory(action: {
  type: string;
  payload?: FilterModel;
}) {
  if (!action.payload) {
    yield put(setIsLoadingTransactionHistoryPageDataAction(true));
  } else {
    if (!action.payload.silent) {
      MessageService.send({ name: MessageNames.SETLOADING, payload: true });
    }
  }
  try {
    const sendingData: any = action.payload ? { ...action.payload } : {};
    if (sendingData.silent) {
      delete sendingData.silent;
    }
    if (sendingData.dwType) {
      sendingData.type = sendingData.dwType;
      delete sendingData.dwType;
    }
    if (sendingData.address) {
      delete sendingData.address;
    }
    if (action.payload && action.payload.start_date) {
      sendingData.page = 0;
      sendingData.page_size = 300;
    }
    const response: StandardResponse = yield call(
      getTransactionHistoryAPI,
      sendingData,
    );
    if (response.status === false) {
      if (!action.payload) {
        yield put(setIsLoadingTransactionHistoryPageDataAction(false));
      } else {
        MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      }
      toast.error('error getting transaction history data');
      return;
    }
    const currectedData = response.data.payments;
    for (let i = 0; i < currectedData.length; i++) {
      currectedData[i].createdAtToFilter = Number(
        currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
      );
    }
    if (!action.payload || action.payload?.silent) {
      yield put(
        setTransactionHistoryPageDataAction(
          response.data.payments ? currectedData : [],
        ),
      );
    }
    if (!action.payload?.silent) {
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    }
    if (action.payload) {
      MessageService.send({
        name: MessageNames.SET_GRID_DATA,
        payload: response.data.payments ? currectedData : [],
      });
    }
  } catch (error) {
    // yield put(setIsLoadingTransactionHistoryPageDataAction(false));
    if (!action.payload) {
      yield put(setIsLoadingTransactionHistoryPageDataAction(false));
    } else {
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    }
    toast.error('error while getting data');
  }
}

function* getOrderDetail(action: {
  type: string;
  payload: { id: number; rowId: string };
}) {
  try {
    const response: StandardResponse = yield call(getOrderDetailAPI, {
      id: action.payload.id,
    });
    if (response.status === false) {
      toast.error('error getting details');
      return;
    }
    Subscriber.next({
      name: MessageNames.SET_ORDER_DETAIL,
      orderId: action.payload.id,
      rowId: action.payload.rowId,
      payload: response.data,
    });
  } catch (error) {
    toast.error('error while getting data');
  }
}
function* preWithdraw(action: { type: string; payload: WithdrawModel }) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(
      preWithdrawAPI,
      action.payload,
    );
    if (response.status === false) {
      if (response.message && response.message.length > 0) {
        toast.warn(response.message);
      }
      ToastMessages(response.data);
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      return;
    } else if (
      response.data &&
      (response.data.need2fa === true || response.data.needEmailCode === true)
    ) {
      MessageService.send({
        name: MessageNames.OPEN_WITHDRAW_VERIFICATION_POPUP,
        payload: response.data,
      });
    } else {
      yield put(withdrawAction(action.payload));
    }
    //else {
    //	MessageService.send({
    //		name: MessageNames.OPEN_WITHDRAW_VERIFICATION_POPUP,
    //		payload: {need2fa: true,needEmailCode: true}
    //	});
    //}

    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  } catch (error) {
    toast.error('withdraw error');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
}
function* withdraw(action: { type: string; payload: WithdrawModel }) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(withdrawAPI, action.payload);

    if (response.status === false) {
      if (response.message && response.message.length > 0) {
        toast.warn(response.message);
      }

      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      return;
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    yield put(
      addWithdrawDataAction({
        ...response.data.payments[0],
        address: action.payload.address,
      }),
    );
    MessageService.send({
      name: MessageNames.ADD_DATA_ROW_TO_WITHDRAWS,
      payload: response.data.payments[0],
    });
    toast.info('Started withdraw process');
  } catch (error) {
    toast.error('withdraw error');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
}

function* addNewAddress(action: {
  type: string;
  payload: { address: string; code: string; label: string };
}) {
  Subscriber.next({
    name: MessageNames.SET_POPUP_LOADING,
    payload: true,
  });
  const data = action.payload;

  try {
    const response: StandardResponse = yield call(
      addNewWithDrawAddressAPI,
      data,
    );
    if (response.status === false) {
      ToastMessages(response.data);
    } else if (response.status === true) {
      toast.success('address added successfully');
      try {
        Subscriber.next({
          name: MessageNames.ADDITIONAL_ACTION,
          payload: {
            title: 'removeAddButton',
            value: response.data[0],
          },
        });
        yield put(addFormerWithdrawAddressesAction(response.data[0]));
        yield put(addOneToWithdrawAddressesAction(response.data[0]));
      } catch (e) {}
    }
    Subscriber.next({
      name: MessageNames.SET_POPUP_LOADING,
      payload: false,
    });
    Subscriber.next({
      name: MessageNames.CLOSE_MODAL,
    });
  } catch (error) {
    toast.error('address not added');
    Subscriber.next({
      name: MessageNames.CLOSE_MODAL,
    });
  }
}

function* getUserData(action: {
  type: string;
  payload: { id: number; rowId: string };
}) {
  try {
    const response: StandardResponse = yield call(getUserDataAPI);
    if (response.status === false) {
      toast.error('error getting User Data');
      return;
    }
    yield put(setUserDataAction(response.data));
  } catch (error) {
    toast.error('error while getting user data');
  }
}

// Individual exports for testing
export default function* fundsPageSaga() {
  yield takeLatest(ActionTypes.GET_BALANCE_PAGE_DATA_ACTION, getBalances);
  yield takeLatest(
    ActionTypes.GET_DEPOSITE_AND_WITHDRAWS_DATA_ACTION,
    getDepositAndWithdraws,
  );
  yield takeLatest(
    ActionTypes.GET_RAW_DEPOSITE_AND_WITHDRAWS_DATA_ACTION,
    getRawDepositAndWithdraws,
  );
  yield takeLatest(
    ActionTypes.GET_TRANSACTION_HISTORY_PAGE_DATA_ACTION,
    getTransactionHistory,
  );
  yield takeLatest(
    ActionTypes.GET_FORMER_WITHDRAW_ADDRESSES,
    getFormerWithdrawAddresses,
  );
  yield takeEvery(ActionTypes.GET_PAYMENT_DETAIL_ACTION, getOrderDetail);
  yield takeEvery(ActionTypes.GET_USER_DATA_ACTION, getUserData);
  yield takeLatest(ActionTypes.WITHDRAW_ACTION, withdraw);
  yield takeLatest(ActionTypes.PRE_WITHDRAW_ACTION, preWithdraw);
  yield takeLatest(ActionTypes.ADD_NEW_ADDRESS_ACTION, addNewAddress);
  yield takeLatest(ActionTypes.GET_INFINITE_DW_ACTION, getInfiniteDw);
}
