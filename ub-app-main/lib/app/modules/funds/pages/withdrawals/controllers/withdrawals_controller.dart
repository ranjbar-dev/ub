import 'package:dio/dio.dart' show DioError;
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../services/storageKeys.dart';
import '../../../../../../utils/commonUtils.dart';
import '../../../../../../utils/extentions/basic.dart';
import '../../../../../../utils/logger.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../../utils/mixins/toast.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../global/autocompleteModel.dart';
import '../../../../../routes/app_pages.dart';
import '../../../../separateMessagePage/views/separate_message_page_view.dart';
import '../../../../withdrawAddressManagement/controllers/withdraw_address_management_controller.dart';
import '../../../../withdrawAddressManagement/views/withdraw_address_management_view.dart';
import '../../../withdraw_deposit_data_model.dart';
import '../../deposits/providers/depositsProvider.dart';
import '../pre_withdraw_model.dart';
import '../providers/withdrawalProvider.dart';
import '../views/withdrawalEditPopup.dart';

class WithdrawalsController extends GetxController
    with Formatter, Toaster, Popups {
  final storage = GetStorage();
  final selectedCoin = AutoCompleteItem(name: '').obs;
  final selectedNetwork = OtherNetworksConfigsAndAddresses().obs;
  final withdrawAndDepositData = WithdrawDepositDataModel().obs;
  final WithdrawalProvider withdrawalProvider = WithdrawalProvider();
  final isScannerOpen = false.obs;
  final isSubmitting = false.obs;
  final showLoadingInsidePopup = false.obs;
  final isLoadingWithdrawAndDepositData = false.obs;
  final savedCoins = <AutoCompleteItem>[].obs;
  final address = ''.obs;
  final amount = ''.obs;
  final fee = '0.00'.obs;
  final youGet = ''.obs;
  final selectedNetworkIndex = 0.obs;
  final selectedPercentIndex = (-1).obs;
  final numberOfPercentSegments = 4;
  final depositProvider = DepositsProvider();

  void handleCoinSelected(AutoCompleteItem coin) {
    _reset();
    selectedCoin.value = coin;
    saveCoinToHistory(
      coin: coin,
      stream: savedCoins,
      storageKey: StorageKeys.savedWithdrawalCoins,
    );
    getUserDepositData(code: coin.name);
  }

  void handleNetworkChange(int i) {
    selectedNetworkIndex.value = i;
    selectedNetwork.value =
        withdrawAndDepositData.value.networksConfigsAndAddresses[i];
    _placeFee();
  }

  void getUserDepositData(
      {@required String code, bool openPopupAfter = true}) async {
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
          selectedNetworkIndex.value = 0;
        }
        if (openPopupAfter) {
          _openWithdrawalEditPopup();
        }
        _placeFee();
      }
    } catch (e) {
      log.e('error while getting withDrawData');
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
  void onClose() {}

  void handleAddressChange(String v) {
    address.value = v;
  }

  void handleAmountChange(String v) {
    amount.value = v;
  }

  void handleScanButtonTap() async {
    final dataFromQrScan = await Get.toNamed(
      AppPages.QR_SCAN,
    );
    debugPrint(dataFromQrScan);

    if (dataFromQrScan != null) {
      address.value = dataFromQrScan;
    } else {
      Get.back();
    }
  }

  void handleAllAmountClick() {
    handlePercentClick(numberOfPercentSegments - 1);
    //final balance = withdrawAndDepositData.value.balance;
    //if (balance != null) {
    //  amount.value = balance.availableAmount.currencyFormat();

    //}
  }

  void handleSubmitClick() async {
    var youWillGet;
    final lfee = withdrawAndDepositData.value.balance.fee;
    if (amount.value != '' && lfee != '') {
      youWillGet =
          (double.parse(amount.value.replaceAll(',', '')) - double.parse(lfee))
              .toString();
    }

    openWithdrawSubmitPopup(
      coin: selectedCoin.value,
      onSubmit: finalSubmit,
      address: address.value,
      amount: amount.value,
      network: selectedNetwork.value.completedNetworkName != ''
          ? selectedNetwork.value.completedNetworkName
          : withdrawAndDepositData.value.completedNetworkName,
      transactionFee: lfee,
      youWillGetAmount: youWillGet,
    );
  }

  void finalSubmit() async {
    try {
      final networkText = selectedCoin.value.code != selectedNetwork.value.code
          ? selectedNetwork.value.code
          : "";
      isSubmitting.value = true;
      final sendingData = PreWithdrawModel(
        address: address.value,
        amount: amount.value.replaceAll(',', ''),
        code: selectedCoin.value.code,
        network: networkText,
      );
      final response = await withdrawalProvider.preWithdraw(data: sendingData);
      if (response['status'] == true) {
        if (response['data']['need2fa'] != true &&
            response['data']['needEmailCode'] != true) {
          requestWithdraw();
        } else {
          final data = response['data'];
          isSubmitting.value = false;
          openVerificationPopup(
              need2fa: data['need2fa'],
              needEmailCode: data['needEmailCode'],
              isNewDevice: false,
              onSubmit: (v) {
                //print(v);
                Get.back();
                requestWithdraw(
                  emailCode: v['emailCode'],
                  g2faCode: v['s2fa'],
                );
              });
        }
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isSubmitting.value = false;
    }
  }

  void handleSearchAddressClick() async {
    Get.put(WithdrawAddressManagementController());
    final result = await Get.to(
      () => WithdrawAddressManagementView(
          selectAddress: true, code: selectedCoin.value.code),
    );
    if (result != null) {
      address.value = result;
    }
  }

  void _openWithdrawalEditPopup() {
    Get.bottomSheet(
      WithdrawalEditPopup(),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  void requestWithdraw({String g2faCode, String emailCode}) async {
    try {
      final networkText = selectedCoin.value.code != selectedNetwork.value.code
          ? selectedNetwork.value.code
          : "";
      isSubmitting.value = true;
      final sendingData = PreWithdrawModel(
        address: address.value,
        amount: amount.value.replaceAll(',', ''),
        code: selectedCoin.value.code,
        network: networkText,
      );
      if (g2faCode != null && g2faCode != '') {
        sendingData.g2faCode = g2faCode;
      }
      if (emailCode != null && emailCode != '') {
        sendingData.emailCode = emailCode;
      }

      final response = await withdrawalProvider.withdraw(data: sendingData);
      if (response['status'] == true) {
        toastSuccess('Withdrawal completed');
        if (Get.isBottomSheetOpen) {
          Get.back();
        }

        Get.to(() => SeparateMessagePageView(
              image: Assets.images.withdrawalSuccessful.image(width: 140.0),
              texts: [
                SeparateMessagePageText(
                  text: 'Your Withdraw request process has been',
                ),
                SeparateMessagePageText(
                  text: 'Started',
                  color: ColorName.green,
                ),
                SeparateMessagePageText(
                  text: 'You can check the progress in transaction history',
                  color: ColorName.orange,
                ),
              ],
            ));

        getUserDepositData(
          code: selectedCoin.value.code,
          openPopupAfter: false,
        );
        //handleCoinSelected(selectedCoin.value);
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isSubmitting.value = false;
    }
  }

  void handleEmptyBalanceClick() async {
    toastAction(
        "You have 0 ${selectedCoin.value.name} balance!",
        UBButton(
          height: 26,
          fontSize: 11.0,
          width: 85,
          onClick: () async {
            Get.back();
            isLoadingWithdrawAndDepositData.value = true;
            await openDepositPopup(coin: selectedCoin.value);
            isLoadingWithdrawAndDepositData.value = false;
          },
          text: 'Deposit ${selectedCoin.value.name}',
        ),
        duration: 4000);
  }

  handlePercentClick(int index) {
    final balance = withdrawAndDepositData.value.balance;
    double available = 0.0;
    if (balance != null) {
      available = double.parse(balance.availableAmount);
    }

    if (index == selectedPercentIndex.value) {
      selectedPercentIndex.value = -1;
      amount.value = '';
    } else {
      selectedPercentIndex.value = index;

      selectedPercentIndex.value = index;
      amount.value = (available * ((index + 1) / numberOfPercentSegments))
          .toString()
          .currencyFormat();
    }
  }

  void _reset() {
    amount.value = '';
    address.value = '';
    selectedPercentIndex.value = -1;
    selectedNetwork.value = OtherNetworksConfigsAndAddresses();
    selectedNetworkIndex.value = 0;
    _placeFee(placeDirectly: "0.00");
  }

  void _placeFee({String placeDirectly}) {
    if (placeDirectly != null) {
      fee.value = placeDirectly;
      return;
    }
    final config = withdrawAndDepositData.value.networksConfigsAndAddresses
        .firstWhere((element) => element.code == selectedNetwork.value.code);
    if (config != null) {
      fee.value = config.fee;
    }
  }
}
