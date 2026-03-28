import 'package:get/get.dart';

import '../../../../../../services/constants.dart';
import '../../../../../../utils/logger.dart';
import '../../../../../global/controller/authorizedMqttController.dart';
import '../../../../../global/providers/commonDataProvider.dart';
import '../../../../account/user_model.dart';
import '../../../balance_response_model_model.dart';
import '../../../controllers/funds_controller.dart';
import '../providers/balanceProvider.dart';

class BalanceController extends GetxController {
  final AuthorizedMqttController authorizedMqttController = Get.find();
  final isLoading = false.obs;
  final isSilentLoading = false.obs;
  final showSmallBalances = true.obs;
  final showAvailableData = true.obs;
  final isHeadOpen = false.obs;
  LightSubscription<List<RxUpdateables>> updateSubscription;
  final balanceProvider = BalanceProvider();
  final balancesAllData = BalanceResponseModel().obs;
  final commonDataProvider = CommonDataProvider();

  void handelShowAvailableDataToggle() {
    showAvailableData.toggle();
    return;
  }

  void toggleHeadOpen() {
    isHeadOpen.toggle();
  }

  void handleShowSmallBalancesChange(bool val) {
    showSmallBalances.value = val;
  }

  @override
  void onInit() {
    updateSubscription =
        authorizedMqttController.updateDataSubject.listen((value) {
      if (value is List && (value.indexOf(RxUpdateables.Balances) != -1)) {
        getBalances(silent: true);
      }
      return;
    });
    getBalances();
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  void getBalances({bool silent}) async {
    if (silent == true) {
      isSilentLoading.value = true;
    } else {
      isLoading.value = true;
    }
    try {
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
        balancesAllData.value = balanceData;

        return;
      }
      log.e('error  getting balances ' + 'status=false');
    } catch (e) {
      log.e('error while getting balances' + e.toString());
    } finally {
      isLoading.value = false;
      isSilentLoading.value = false;
    }
  }

  @override
  void onClose() {
    if (updateSubscription != null) {
      updateSubscription.cancel();
    }
  }

  handleRefreshBalances() {
    getBalances(silent: true);
  }
}
