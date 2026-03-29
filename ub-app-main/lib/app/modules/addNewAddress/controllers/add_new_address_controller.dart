import 'package:dio/dio.dart' show DioError;
import 'package:get/get.dart';
import '../../../global/autocompleteModel.dart';
import '../../../global/currency_model.dart';
import '../add_new_address_model.dart';
import '../providers/addNewAddressProvider.dart';
import '../../funds/withdraw_deposit_data_model.dart';
import '../../withdrawAddressManagement/controllers/withdraw_address_management_controller.dart';
import '../../withdrawAddressManagement/withdraw_address_model.dart';
import '../../../routes/app_pages.dart';
import '../../../../services/constants.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/pairAndCurrencyUtils.dart';

class AddNewAddressController extends GetxController with Toaster {
  final selectedCoin = AutoCompleteItem(name: '').obs;
  final currencyArray = Constants.currencyArray();
  final addNewAddressProvider = AddNewAddressProvider();
  final WithdrawAddressManagementController
      withdrawAddressManagementController = Get.find();
  final selectedNetworkIndex = 0.obs;
  final selectedNetwork = OtherNetworksConfigsAndAddresses().obs;
  final newAddressLabel = ''.obs;
  final address = ''.obs;
  final isAddingNewAddress = false.obs;
  final networks = <OtherNetworksConfigsAndAddresses>[].obs;
  final coinMap =
      Map<String, CurrencyModel>.from(PairAndCurrencyUtils.coinsMap.value);
  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void handleCoinSelected(AutoCompleteItem coin) {
    selectedNetwork.value = OtherNetworksConfigsAndAddresses();
    selectedCoin.value = coin;
    final coinData = coinMap[coin.code];
    if (coinData.otherBlockChainNetworks.length > 0) {
      final tmp = [
        OtherNetworksConfigsAndAddresses(
          code: coinData.mainNetwork,
          completedNetworkName: coinData.completedNetworkName,
        )
      ];
      coinData.otherBlockChainNetworks.forEach((element) {
        tmp.add(OtherNetworksConfigsAndAddresses(
          code: element["code"],
          completedNetworkName: element["completedNetworkName"],
        ));
      });
      networks.assignAll(tmp);
    } else {
      networks.assignAll([]);
      selectedNetwork.value = OtherNetworksConfigsAndAddresses();
    }
    print(coin);
  }

  void handleNetworkChange(int i) {
    selectedNetworkIndex.value = i;
  }

  void handleNewAddressLabelChange(String v) {
    newAddressLabel.value = v;
  }

  void handleAddressChange(String v) {
    address.value = v;
  }

  void handleScanButtonTap() async {
    final dataFromQrScan = await Get.toNamed(
      AppPages.QR_SCAN,
    );
    if (dataFromQrScan != null) {
      address.value = dataFromQrScan;
    }
  }

  void handleCreateClick() async {
    try {
      isAddingNewAddress.value = true;
      final data = AddNewAddressModel(
        address: address.value,
        code: selectedCoin.value.code,
        label: newAddressLabel.value,
      );
      if (selectedNetwork.value.code != null) {
        data.network = selectedNetwork.value.code;
      }

      final response = await addNewAddressProvider.addNewAddress(data: data);
      if (response['status'] == true) {
        final List responseData = response['data'];
        if (responseData == null || responseData.length == 0) return;
        final code = responseData[0]['code'];
        var model = {};
        for (var currency in currencyArray) {
          if (currency.name == code) {
            model = responseData[0];
            model["icon"] = currency.image;
            break;
          }
        }
        final WithdrawAddressModel wamodel =
            WithdrawAddressModel.fromJson(model);
        withdrawAddressManagementController.withdrawAddresses
            .insert(0, wamodel);
        Get.back();
        toastSuccess('Added new withdraw address');
        reset();
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isAddingNewAddress.value = false;
    }
  }

  void reset() {
    address.value = '';
    newAddressLabel.value = '';
    selectedCoin.value = AutoCompleteItem(name: '');
  }
}
