import { toast } from "app/components/Customized/react-toastify";
import { delay, put, takeLatest } from "redux-saga/effects";
import { MessageNames, MessageService } from "services/messageService";
import { GetLiquidityOrdersAPI, UpdateCommissionReportAPI } from "services/orderManagementService";
import { safeApiCall } from "utils/sagaUtils";

import { LiquidityOrdersActions } from "./slice";

export function* GetLiquidityOrders(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetLiquidityOrdersAPI, action.payload);
  if (response) {
    yield put(LiquidityOrdersActions.setLiquidityOrdersData(response.data as Record<string, unknown>));
  }
}

export function* UpdateCommissionReport(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'updateCommisionReportButton',
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateCommissionReportAPI, action.payload);
    if (response) {
      yield delay(2000);
      toast.success('commissionsUpdated');
      yield put(
        LiquidityOrdersActions.GetLiquidityOrdersAction(action.payload),
      );
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'updateCommisionReportButton',
      payload: false,
    });
  }
}

export function* liquidityOrdersSaga() {
  yield takeLatest(LiquidityOrdersActions.GetLiquidityOrdersAction.type, GetLiquidityOrders);
  yield takeLatest(LiquidityOrdersActions.UpdateCommissionReportAction.type, UpdateCommissionReport);
}