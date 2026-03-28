import 'package:basic_utils/basic_utils.dart';
import 'package:dio/dio.dart' show DioError, DioErrorType;
import 'package:flutter/material.dart';

import '../../app/common/custom/toaster/utopic_toast.dart';
import '../../app/global/authorized_order_event_model.dart';
import '../extentions/basic.dart';

mixin Toaster {
  static final toastManager = ToastManager();
  void toast({String message, Widget action, ToastType type, int duration}) {
    toastManager.showToast(
      message,
      type: type,
      duration: duration != null
          ? Duration(milliseconds: duration > 1000 ? duration : duration * 1000)
          : null,
      endAction: action,
      action: ToastAction(
        onPressed: (hideToastFn) {
          hideToastFn();
        },
      ),
    );
  }

  /// duration in milliseconds
  toastSuccess(
    String message, {
    int duration,
  }) {
    toast(
      message: message,
      type: ToastType.success,
      duration: duration,
    );
  }

  toastError(String message, {int duration}) {
    toast(
      message: message,
      type: ToastType.error,
      duration: duration,
    );
  }

  toastWarning(String message, {int duration}) {
    toast(
      message: message,
      type: ToastType.warning,
      duration: duration,
    );
  }

  void toastInfo(
    String message, {
    int duration,
  }) {
    toast(
      message: message,
      type: ToastType.info,
      duration: duration,
    );
  }

  void toastAction(
    String message,
    Widget action, {
    int duration,
  }) {
    toast(
      message: message,
      action: action,
      type: ToastType.info,
      duration: duration,
    );
  }

  toastDioError(DioError e, {ToastType type = ToastType.warning}) {
    if (e.type == DioErrorType.connectTimeout) {
      toast(
        message: 'Please check your connection and try again!',
        type: ToastType.error,
      );
    }
    final errorData = e.response.data['data'];
    final errorMessage = e.response.data["message"];
    if (e.response.statusCode == 422) {
      if (errorData == null || (errorData != null && errorData.isEmpty)) {
        if (errorMessage != null && errorMessage.length > 0) {
          toast(
            message: StringUtils.capitalize(errorMessage),
            type: type,
            duration: 6000,
          );
        }
      } else {
        final errorList = [];
        errorData.values.forEach((v) => errorList.add(v));

        String errorToShow = errorMessage;
        if (errorList.length != 0) {
          errorToShow = errorList[0];
        }
        toast(
          message: StringUtils.capitalize(errorToShow),
          type: type,
          duration: 5000,
        );
      }
    }
    if (e.response.statusCode == 403) {
      final errorMessage = e.response.data["message"];
      toast(
        message: StringUtils.capitalize(errorMessage),
        type: type,
        duration: 5000,
      );
    }
    if (e.response.statusCode == 401) {
      if (e.response.data['message'] != null) {
        if (e.response.data["message"].toString().contains("2fa")) {
          toast(
            message: '2fa code is not correct',
            type: type,
          );
          return;
        }
      }
      toast(
        message: 'Invalid email or password',
        type: type,
      );
    } else if (e.response.statusCode == 500) {
      toast(
        message: 'Error Connecting To Server ,please try again!',
        type: ToastType.error,
      );
    }
  }

  toastAuthorizedEvent(AuthorizedOrderEventModel payload) {
    if (payload.status == 'open') {
      payload.status = 'Placed';
    }
    if (payload.status == 'filled') {
      payload.status = 'Filled';
    }
    if (payload.status == 'canceled') {
      payload.status = 'Canceled';
    }
    final coin = payload.pairCurrency.split('-')[0].toString();
    final otherCoin = payload.pairCurrency.split('-')[1].toString();
    final formattedAmount = (payload.amount != '' && payload.amount != null)
        ? (payload.amount.currencyFormat(removeInsignificantZeros: true)) + ''
        : '';
    final formattedPrice = (payload.price != '' && payload.price != null)
        ? (payload.price.currencyFormat(removeInsignificantZeros: true)) &
            otherCoin &
            '|'
        : 'Market price' & '|';
    final t = [
      payload.type == 'buy' ? "BUY |" : "SELL |",
      "$formattedAmount $coin ${payload.amount != '' ? ' ON' : ' '}",
      "$formattedPrice [${payload.status}]",
    ];

    final toastText = t[0] & t[1] & t[2];
    switch (payload.status.toLowerCase()) {
      case 'placed':
        toastInfo(toastText, duration: 6000);
        break;
      case 'canceled':
        toastWarning(toastText, duration: 4000);
        break;
      case 'filled':
        toastSuccess(toastText, duration: 6000);
        break;
      default:
        break;
    }
  }
}
