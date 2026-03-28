import 'package:dio/dio.dart';
import 'package:get/get.dart';

import '../../../../../../utils/logger.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../../utils/mixins/toast.dart';
import '../../../../../global/autocompleteModel.dart';
import '../../../../../global/providers/commonDataProvider.dart';
import '../../../../account/user_model.dart';
import '../../../../exchange/controllers/exchange_controller.dart';
import '../../../balance_response_model_model.dart';
import '../../../controllers/funds_controller.dart';
import '../../balance/providers/balanceProvider.dart';
import '../provider/auto_exchange_provider.dart';

class AutoExchangeController extends GetxController with Toaster, Popups {
  final switchValue = false.obs;
  final isSubmitLoading = false.obs;
  final isCoinsListLoading = false.obs;
  Balance currentBalance = Balance();
  final autoExchangeProvider = AutoExchangeProvider();
  final balanceProvider = BalanceProvider();
  final commonDataProvider = CommonDataProvider();
  ExchangeController exchangeController = ExchangeController();

  Rx<BalanceResponseModel> coinsList = BalanceResponseModel().obs;
  var balances = [].obs;
  AutoExchangeModel model = AutoExchangeModel();
  var canSubmitAutoExchange = false.obs;

  var searchedCoin = AutoCompleteItem(name: "").obs;

  @override
  void onInit() {
    exchangeController = Get.find<ExchangeController>();
    getBalances();

    exchangeController.pairLocalInfo.dependantCoin.stream.listen((value) async {
      if (!["", null].contains(currentBalance.autoExchangeCode) &&
          currentBalance.autoExchangeCode != value.code)
        canSubmitAutoExchange.value = true;
    });

    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void getBalances() async {
    try {
      isCoinsListLoading.value = true;
      final data = await Future.wait(
          [balanceProvider.getBalances(), commonDataProvider.getUserData()]);

      final userData = data[1];
      if (userData['status'] == true) {
        final userInfo = UserModel.fromJson(
          userData["data"],
        );
        final FundsController fundsController = Get.find();
        fundsController.isUserVerified.value = userInfo.isAccountVerified;
      }
      final response = data[0];
      if (response['status'] == true) {
        final balanceData = BalanceResponseModel.fromJson(response['data']);
        coinsList.value = balanceData;
        balances = [].obs;
        balances.addAll(coinsList.value.balances);

        return;
      }
      log.e('error  getting balances ' + 'status=false');
    } catch (e) {
      log.e('error while getting balances' + e.toString());
    } finally {
      isCoinsListLoading.value = false;
    }
  }

  void toggleAutoExchange({Balance balance}) async {
    switchValue.toggle();
    //In Case That User Changed Any Value Related To AutoExchange, The Apply Buttom Should Be Enabled
    if ((switchValue.value &&
            balance.autoExchangeCode !=
                exchangeController.pairLocalInfo.dependantCoin.value.code) ||
        (!switchValue.value &&
            (balance.autoExchangeCode != "" &&
                balance.autoExchangeCode != null)))
      canSubmitAutoExchange.value = true;
    else
      canSubmitAutoExchange.value = false;
    if (switchValue.value) {
      model.autoExchangeCode.value =
          exchangeController.pairLocalInfo.dependantCoin.value.code;
    }
  }

  void setCurrentBalance(Balance balance) {
    currentBalance = balance;
    canSubmitAutoExchange.value = false;
    if (!["", null].contains(balance.autoExchangeCode))
      switchValue.value = true;
    else {
      switchValue.value = false;
      canSubmitAutoExchange.value = false;
    }
    exchangeController.handleAutoExchangeBaseCoinSelected(
        coin: AutoCompleteItem(
            code: balance.code,
            name: balance.name,
            image: balance.image,
            desc: balance.name),
        autoExchangeCode: balance.autoExchangeCode);
  }

  void submitAutoExchange() async {
    try {
      isSubmitLoading.value = true;

      model
        ..code.value = currentBalance.code
        ..autoExchangeCode.value =
            exchangeController.pairLocalInfo.dependantCoin.value.code
        ..mode.value = switchValue.value ? "add" : "delete";

      final response =
          await autoExchangeProvider.submitAutoExchange(model: model);
      if (response['status'] == true) {
        isSubmitLoading.value = false;

        Get.back();
        canSubmitAutoExchange.value = false;
        getBalances();
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isSubmitLoading.value = false;
    }
  }

  Future handleCoinSelected(
      {AutoCompleteItem coin, bool isFromSearch = false}) async {
    isCoinsListLoading.value = true;
    setCurrentBalance(
        balances.firstWhere((element) => element.code == coin.code));
    balances = [].obs;
    balances.add(currentBalance);
    if (isFromSearch) searchedCoin.value = coin;
    isCoinsListLoading.value = false;
  }
}

class AutoExchangeModel {
  final code = ''.obs;
  final autoExchangeCode = ''.obs;
  final mode = ''.obs;
}
