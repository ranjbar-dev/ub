import 'dart:async';

import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
//import 'package:unitedbit/app/modules/trade/controllers/trade_controller.dart';
import 'package:unitedbit/generated/locales.g.dart';
import 'package:unitedbit/services/constants.dart';
import 'package:unitedbit/utils/mixins/popups.dart';
import 'package:unitedbit/utils/mixins/toast.dart';

import '../../../../../common/components/UBText.dart';
import '../../../../../global/controller/authorizedCentrifugoController.dart';
import '../../../order_model.dart';
import '../../../widgets/orderRow/openOrderRow.dart';
import '../cancel_order_model.dart';
import '../providers/open_orders_provider.dart';

class OpenOrdersController extends GetxController with Toaster, Popups {
  final AuthorizedCentrifugoController authorizedCentrifugoController = Get.find();

  final OpenOrdersProvider openOrdersProvider = OpenOrdersProvider();
  GlobalKey<AnimatedListState> openOrderListKey = GlobalKey();
  GlobalKey<AnimatedListState> fullScreenKey = GlobalKey();
  final selectedOpenOrderFilterText = (LocaleKeys.allOpenOrders.tr).obs;
  final openOrders = <OrderModel>[].obs;
  final isFullScreen = false.obs;
  final isSilentLoading = false.obs;
  bool isCancelling = false;
  final List<OrderModel> initialOpenOrders = [];

  final loadingIds = <int>[].obs;

  final loadingData = false.obs;

  LightSubscription<List<RxUpdateables>> updateSubscription;

  @override
  void onInit() {
    super.onInit();
    updateSubscription =
        authorizedCentrifugoController.updateDataSubject.listen((value) {
      if (value is List && (value.indexOf(RxUpdateables.OpenOrders) != -1)) {
        getOpenOrders(silent: true);
      }
      return;
    });
    isFullScreen.listen((v) {
      if (v == false) {
        final initial = [...initialOpenOrders];
        openOrders.assignAll(initial);
      }
    });
    getOpenOrders();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    super.onClose();
    if (updateSubscription != null) {
      updateSubscription.cancel();
    }
  }

  Future getOpenOrders({bool silent}) async {
    if (silent != true) {
      loadingData.value = true;
    } else {
      isSilentLoading.value = true;
    }
    try {
      final response = await openOrdersProvider.getOpenOrders();
      if (response['status'] == true) {
        final list = await compute(parseOpenOrders, response);
        initialOpenOrders.assignAll(list);
        openOrders.assignAll(list);
        isSilentLoading.value = false;
      }
    } catch (e) {
      print(e.toString());
    } finally {
      isSilentLoading.value = false;
      loadingData.value = false;
    }
    return Future.value();
  }

  void handleCancelClick({int id, int index}) async {
    if (isCancelling) {
      return;
    }
    openConfirmation(
      onConfirm: () => _cancelOrder(id: id, index: index),
      titleWidget: UBText(
        text: 'Cancel this Order?',
      ),
      cancelText: 'No',
      confirmText: 'Yes',
    );
  }

  _cancelOrder({int id, int index}) async {
    isCancelling = true;
    loadingIds.add(id);
    try {
      final response = await openOrdersProvider.cancelOrders(
        model: CancelOrderModel(orderId: id),
      );
      if (response['status'] == true) {
        loadingIds.remove(id);
        // ignore: invalid_use_of_protected_member
        final tmp = List<OrderModel>.from(openOrders.value);
        tmp.removeAt(index);
        initialOpenOrders.assignAll(tmp);
        openOrders.assignAll(tmp);
      }
    } on DioError catch (e) {
      toastDioError(e);
    } catch (e) {
      loadingIds.remove(id);
    } finally {
      isSilentLoading.value = false;
      isCancelling = false;
      loadingIds.remove(id);
    }
  }

  filterOpenOrders(String v) {
    selectedOpenOrderFilterText.value = v;
    setFilter(v);
  }

  void setFilter(String v) {
    final initial = [...initialOpenOrders];
    switch (v) {
      case 'All Open Orders':
        openOrders.assignAll(initial);
        return;
      case 'Only Sell Orders':
        final filtered = initial.where((element) => element.side == 'sell');
        openOrders.assignAll(filtered);
        return;
      case 'Only Buy Orders':
        final filtered = initial.where((element) => element.side == 'buy');
        openOrders.assignAll(filtered);
        return;
        break;
      default:
        openOrders.assignAll(initial);
        return;
    }
  }

  GlobalKey<AnimatedListState> keySelector(bool isFullScreen,
      {bool returnOldOne}) {
    if (isFullScreen == true) {
      if (returnOldOne == true) {
        return fullScreenKey;
      }
      fullScreenKey = GlobalKey();
      return fullScreenKey;
    }
    if (returnOldOne == true) {
      return openOrderListKey;
    }
    openOrderListKey = GlobalKey();
    return openOrderListKey;
  }

  void removeFromList({OrderModel order, int index, bool fullScreen}) {
    keySelector(fullScreen, returnOldOne: true).currentState.removeItem(
      index,
      (BuildContext context, Animation<double> animation) {
        return FadeTransition(
          opacity:
              CurvedAnimation(parent: animation, curve: Interval(0.5, 1.0)),
          child: SizeTransition(
            sizeFactor:
                CurvedAnimation(parent: animation, curve: Interval(0.0, 1.0)),
            axisAlignment: 0.0,
            child: OpenOrderRow(
              isCancelLoading: true,
              order: order,
            ),
          ),
        );
      },
      duration: Duration(milliseconds: 300),
    );
  }
}

FutureOr parseOpenOrders(dynamic response) {
  return List<OrderModel>.from(
    response["data"].map(
      (model) => OrderModel.fromJson(model),
    ),
  );
}
