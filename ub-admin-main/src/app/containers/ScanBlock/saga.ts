import { toast } from 'app/components/Customized/react-toastify';
import { takeLatest } from 'redux-saga/effects';
import { ScanBlockAPI } from 'services/userManagementService';
import { LoadingIds, SetLoading } from 'utils/loading';
import { safeApiCall } from 'utils/sagaUtils';

import { ScanBlockActions } from "./slice";

export function* scan(action: { type: string; payload: { network: string; block_number: number } }) {
  SetLoading({ id: LoadingIds.ScanButton, loading: true });
  try {
    const response = yield* safeApiCall(ScanBlockAPI, action.payload, { toastOnError: false });
    if (response) {
      toast.success('Scan Process Started');
    } else {
      toast.warn('Scan Process failed');
    }
  } finally {
    SetLoading({ id: LoadingIds.ScanButton, loading: false });
  }
}

export function* scanBlockSaga() {
  yield takeLatest(ScanBlockActions.Scan.type, scan);
}