import ActionTypes from './constants';
import { takeLatest, put, all, call, takeEvery } from 'redux-saga/effects';
import {
  setIsLoadingAction,
  setCurrenciesAction,
  setWithdrawAddressesAction,
  addOneToWithdrawAddressesAction,
  applyDeleteAddressAction,
  applyfavoriteAddressAction,
} from './actions';
import {
  getCurrenciesAPI,
  getWithDrawAddressesAPI,
  addNewWithDrawAddressAPI,
  deleteWithDrawAddressAPI,
  setFavoriteWithDrawAddressAPI,
} from 'services/address_management_service';
import { StandardResponse } from 'services/constants';
import { toast } from 'components/Customized/react-toastify';
import { Subscriber, MessageNames } from 'services/message_service';
import { ToastMessages } from 'services/toastService';
import { WithdrawAddress } from './types';

function * getInitialData (action: { type: string }) {
  yield put(setIsLoadingAction(true));

  try {
    const [
      currenciesResponse,
      withdrawAddressesResponse,
    ]: StandardResponse[] = yield all([
      call(getCurrenciesAPI),
      call(getWithDrawAddressesAPI),
    ]);
    if (currenciesResponse.status === true) {
      yield put(setCurrenciesAction(currenciesResponse.data.currencies));
    }
    if (withdrawAddressesResponse.status === true) {
      yield put(setWithdrawAddressesAction(withdrawAddressesResponse.data));
    }
  } catch (error) {
    yield put(setIsLoadingAction(false));
  }
}
function * addNewAddress (action: {
  type: string;
  payload: { address: string; code: string; label: string; network?: string };
}) {
  Subscriber.next({
    name: MessageNames.SETLOADING,
    element: 'pulsingButton',
    payload: true,
  });
  const data = action.payload;

  try {
    const response: StandardResponse = yield call(
      addNewWithDrawAddressAPI,
      data,
    );

    if (response.status === false) {
      if (response.message) {
        toast.error(response.message);
      }

      ToastMessages(response.data);
    } else {
      toast.success('Address saved');

      yield put(
        addOneToWithdrawAddressesAction(response.data[0] ?? response.data),
      );

      Subscriber.next({
        name: MessageNames.ADD_DATA_ROW_TO_GRID,
        payload: response.data[0] ?? response.data,
        index: 0,
      });
    }
    Subscriber.next({
      name: MessageNames.SETLOADING,
      element: 'pulsingButton',
      payload: false,
    });
  } catch (error) {
    Subscriber.next({
      name: MessageNames.SETLOADING,
      element: 'pulsingButton',
      payload: false,
    });
  }
}
function * deleteAddress (action: {
  type: string;
  payload: { data: WithdrawAddress; rowIndex: number };
}) {
  Subscriber.next({
    name: MessageNames.DELETE_GRID_ROW,
    payload: action.payload.data,
  });
  try {
    const response: StandardResponse = yield call(deleteWithDrawAddressAPI, {
      ids: [action.payload.data.id],
    });
    if (response.status === true) {
      toast.success('Address removed!');
      yield put(applyDeleteAddressAction(action.payload));
    } else {
      revertDelete(action);
      toast.error('error while deleting address');
    }
  } catch (error) {
    revertDelete(action);
    toast.error('error while deleting address');
  }
}

const revertDelete = (action: {
  type: string;
  payload: { data: WithdrawAddress; rowIndex: number };
}) => {
  Subscriber.next({
    name: MessageNames.ADD_DATA_ROW_TO_GRID,
    payload: action.payload.data,
    index: action.payload.rowIndex,
  });
};
function * setFavoriteAddress (action: {
  type: string;
  payload: {
    data: {
      action: string;
      id: number;
    };
    rowIndex: number;
  };
}) {
  Subscriber.next({
    name: MessageNames.SET_FAVIORITE_ADDRESS,
    payload: action.payload.data,
    isFavorite: action.payload.data.action == 'add' ? true : false,
  });
  try {
    const response: StandardResponse = yield call(
      setFavoriteWithDrawAddressAPI,
      action.payload.data,
    );
    if (response.status === true) {
      yield put(applyfavoriteAddressAction(action.payload));
      toast.success(
        ` ${
          action.payload.data.action == 'add' ? 'added to' : 'removed from'
        } favorite addresses`,
      );
    } else {
    }
  } catch (error) {
    Subscriber.next({
      name: MessageNames.SET_FAVIORITE_ADDRESS,
      payload: action.payload.data,
      isFavorite: false,
    });
  }
}
export default function * addressManagementPageSaga () {
  yield takeLatest(ActionTypes.INITIAL_ACTION, getInitialData);
  yield takeLatest(ActionTypes.ADD_NEW_ADDRESS_ACTION, addNewAddress);
  yield takeEvery(ActionTypes.DELETE_ADDRESS_ACTION, deleteAddress);
  yield takeLatest(ActionTypes.FAVORITE_ADDRESS_ACTION, setFavoriteAddress);
}
