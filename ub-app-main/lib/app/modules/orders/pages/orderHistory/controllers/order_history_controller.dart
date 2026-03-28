import 'dart:async';

import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import '../../../../../common/components/UBWrappedButtons.dart';
import '../../../../../global/autocompleteModel.dart';
import '../../../../../global/controller/authorizedMqttController.dart';
import '../../../../../global/response_model.dart';
import '../../../order_history_detail_model.dart';
import '../../../order_model.dart';
import '../providers/order_history_provider.dart';
import '../../../../../../services/constants.dart';
import '../../../../../../utils/logger.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../../utils/mixins/toast.dart';

class OrderHistoryController extends GetxController with Toaster, Popups {
  final OrderHistoryProvider orderHistoryProvider = OrderHistoryProvider();
  final AuthorizedMqttController authorizedMqttController = Get.find();
  bool disableInfiniteScroll = false;
  final loadingData = false.obs;
  final filtered = false.obs;
  final silentLoadingData = false.obs;
  final orderHistory = <OrderModel>[].obs;
  final loadingId = 0.obs;
  final showFilterButton = true.obs;
  //filter
  final selectedDateButtonIndex = (-1).obs;
  final selecteTypeButtonIndex = 0.obs;
  final filterStartDate = ''.obs;
  final filterEndDate = ''.obs;
  final showCanceledOrders = true.obs;
  final filterPair = 'all-all'.obs;
  final dateButtons = [
    WrappedButtonModel(text: '1 Day', value: '1day'),
    WrappedButtonModel(text: '1 Week', value: '1week'),
    WrappedButtonModel(text: '1 Month', value: '1month'),
    WrappedButtonModel(text: '3 Month', value: '3month'),
  ];
  final filterTypeButtons = [
    WrappedButtonModel(text: 'all', value: 'all'),
    WrappedButtonModel(text: 'Buy', value: 'buy'),
    WrappedButtonModel(text: 'Sell', value: 'sell'),
  ];

  LightSubscription<List<RxUpdateables>> updateSubscription;

  @override
  void onInit() {
    super.onInit();
    updateSubscription =
        authorizedMqttController.updateDataSubject.listen((value) {
      if (value is List && (value.indexOf(RxUpdateables.OrderHistory) != -1)) {
        getOrderHistory(silent: true, andResetFilters: true);
      }
      return;
    });
    getOrderHistory();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    if (updateSubscription != null) {
      updateSubscription.cancel();
    }
  }

  Future getOrderHistory({
    bool silent,
    Map<String, dynamic> filters,
    bool andResetFilters,
    int fromId,
  }) async {
    final applyingFilters = filters ?? {};
    if (fromId != null) {
      applyingFilters["last_id"] = fromId;
    }

    if (andResetFilters == true) {
      resetFilters();
    }
    if (silent != true) {
      loadingData.value = true;
    } else {
      silentLoadingData.value = true;
    }
    try {
      final response =
          await orderHistoryProvider.getOrderHistory(data: applyingFilters);
      final responseObject = ResponseModel.fromJson(response);
      if (responseObject.status == true) {
        if (showCanceledOrders.value != true) {
          response['hideCancelled'] = true;
        }
        final data = await compute(parseOrderHistory, response);
        if (data.length < 50) {
          disableInfiniteScroll = true;
        } else {
          disableInfiniteScroll = false;
        }
        if (fromId == null) {
          orderHistory.assignAll(data);
        } else {
          // ignore: invalid_use_of_protected_member
          final newData = [...orderHistory.value, ...data];
          orderHistory.assignAll(newData);
        }
        orderHistory.refresh();
      }
    } catch (e) {
      log.e(e.toString());
    } finally {
      loadingData.value = false;
      silentLoadingData.value = false;
      filtered.value = filters != null && andResetFilters != true;
    }
    return Future.value();
  }

  hadleOrderDetailsClick(OrderModel data) async {
    if (loadingId.value == 0) {
      try {
        loadingId.value = data.id;
        final response = await orderHistoryProvider.getOrderDetails(data.id);
        if (response['status'] == true) {
          final orderDetails =
              OrderHistoryDetailModel.fromJson(response['data'][0]);
          openOrderDetailsPopup(details: orderDetails, originalData: data);
        }
      } on DioError catch (e) {
        toastDioError(e);
      } catch (e) {
        print(e.toString());
      } finally {
        loadingId.value = 0;
      }
    }
  }

  void handleDateButtonClick({int index}) {
    if (selectedDateButtonIndex.value == index) {
      selectedDateButtonIndex.value = -1;
      return;
    }
    selectedDateButtonIndex.value = index;
  }

  void handleFilterSubmitClick() {
    final Map<String, dynamic> data = {};
    data['pair_currency_name'] = filterPair.value;
    if (filterStartDate.value != '') {
      data['start_date'] = filterStartDate.value + ' ' + '00:00:00';
    }
    if (filterEndDate.value != '') {
      data['end_date'] = filterEndDate.value + ' ' + '00:00:00';
    }
    if (selectedDateButtonIndex.value != -1) {
      data['period'] = dateButtons[selectedDateButtonIndex.value].value;
    }
    if (selecteTypeButtonIndex.value != 0) {
      data['type'] = filterTypeButtons[selecteTypeButtonIndex.value].value;
    }
    getOrderHistory(silent: true, filters: data);
    Get.back();
  }

  void handleResetFiltersClick({bool andPop = true}) {
    getOrderHistory(silent: true, andResetFilters: true);
    if (andPop) {
      Get.back();
    }
  }

  void resetFilters() {
    filterEndDate.value = '';
    filterPair.value = 'all-all';
    showCanceledOrders.value = true;
    filterStartDate.value = '';
    selecteTypeButtonIndex.value = 0;
    selectedDateButtonIndex.value = -1;
  }

  handleStartDateSelect(String date) {
    filterStartDate.value = date.split(' ')[0];
  }

  handleEndDateSelect(String date) {
    filterEndDate.value = date.split(' ')[0];
  }

  void handleTypeButtonClick({int index}) {
    selecteTypeButtonIndex.value = index;
  }

  void handleShowCanceledOrdersToggle() {
    showCanceledOrders.toggle();
  }

  void handlePairSelect(AutoCompleteItem v) {
    filterPair.value = v.name;
  }

  void handleListScroll({ScrollDirection direction}) {
    if (direction == ScrollDirection.Up) {
      if (showFilterButton.value == false) {
        showFilterButton.value = true;
      }
      return;
    }
    if (showFilterButton.value == true) {
      showFilterButton.value = false;
    }
  }

  void onListEndRiched() async {
    final itemsLength = orderHistory.length;
    if (itemsLength >= 50 && !disableInfiniteScroll) {
      await getOrderHistory(silent: true, fromId: orderHistory.last.id);
    }
  }
}

List<OrderModel> parseOrderHistory(response) {
  final list = List<OrderModel>.from(
    response["data"].map(
      (model) => OrderModel.fromJson(model),
    ),
  );
  if (response['hideCancelled'] == true) {
    final filtered = list.where((element) => element.status != 'canceled');
    return filtered.toList();
  }
  return list;
}
