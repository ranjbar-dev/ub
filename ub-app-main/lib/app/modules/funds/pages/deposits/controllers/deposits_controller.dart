import 'package:dio/dio.dart' show DioError;
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:meta/meta.dart' show required;
import '../../../../../global/autocompleteModel.dart';
import '../../../../../global/controller/globalController.dart';
import '../providers/depositsProvider.dart';
import '../../../withdraw_deposit_data_model.dart';
import '../../../../../../services/constants.dart';
import '../../../../../../services/storageKeys.dart';
import '../../../../../../utils/commonUtils.dart';
import '../../../../../../utils/logger.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../../utils/mixins/toast.dart';

class DepositsController extends GetxController with Popups, Toaster {
  final storage = GetStorage();
  final isLoadingWithdrawAndDepositData = false.obs;
  final selectedNetworkIndex = 0.obs;
  final withdrawAndDepositData = WithdrawDepositDataModel().obs;
  final selectedNetwork = OtherNetworksConfigsAndAddresses().obs;
  final coinsList = Constants.currencyArray();
  final depositProvider = DepositsProvider();
  final selectedCoin = AutoCompleteItem(name: '').obs;
  final GlobalController globalController = Get.find();
  final savedCoins = <AutoCompleteItem>[].obs;

  Future handleCoinSelected(AutoCompleteItem coin) async {
    selectedCoin.value = coin;
    selectedNetwork.value = OtherNetworksConfigsAndAddresses();
    selectedNetworkIndex.value = 0;
    saveCoinToHistory(
      coin: coin,
      storageKey: StorageKeys.savedDepositCoins,
      stream: savedCoins,
    );
    await getUserDepositData(code: coin.name);
  }

  void handleNetworkChange(int i) {
    selectedNetworkIndex.value = i;
    selectedNetwork.value =
        withdrawAndDepositData.value.networksConfigsAndAddresses[i];
  }

  void resetSelectedCoin() {
    selectedCoin.value = AutoCompleteItem(name: '');
    selectedNetworkIndex.value = 0;
    selectedNetwork.value = OtherNetworksConfigsAndAddresses();
    withdrawAndDepositData.value = WithdrawDepositDataModel();
  }

  Future getUserDepositData({@required String code}) async {
    isLoadingWithdrawAndDepositData.value = true;
    try {
      final response = await depositProvider.getUserDepositData(code: code);
      if (response["status"] == true) {
        withdrawAndDepositData.value =
            WithdrawDepositDataModel.fromJson(response["data"]);
        if (withdrawAndDepositData.value.networksConfigsAndAddresses.length >
            0) {
          selectedNetwork.value =
              withdrawAndDepositData.value.networksConfigsAndAddresses[0];
        }
      }
    } on DioError catch (e) {
      toastDioError(e);
    } catch (e) {
      log.e('error while getting depositData');
    } finally {
      isLoadingWithdrawAndDepositData.value = false;
    }
  }

  @override
  void onInit() {
    final List<dynamic> storedCoins =
        storage.read<List>(StorageKeys.savedDepositCoins);
    if (storedCoins != null) {
      savedCoins.assignAll(
          storedCoins.map((e) => AutoCompleteItem.fromJson(e)).toList());
    }
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    super.onClose();
  }
}
