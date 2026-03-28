import 'package:flutter/widgets.dart';
import 'package:get/get.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../services/constants.dart';
import '../../../../utils/mixins/popups.dart';
import '../providers/addressMaganagementProvider.dart';
import '../withdraw_address_model.dart';

class WithdrawAddressManagementController extends GetxController with Popups {
  final AddressManagementProvider addressManagementProvider =
      AddressManagementProvider();
  final currencyArray = Constants.currencyArray();

  final loadingData = true.obs;
  final isSilentLoading = true.obs;
  final isRefreshing = false.obs;
  final withdrawAddresses = <WithdrawAddressModel>[].obs;

  @override
  void onInit() {
    super.onInit();
    getAddresses();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  Future getAddresses({bool silent}) async {
    if (silent != true) {
      loadingData.value = true;
    } else {
      isSilentLoading.value = true;
    }
    try {
      final response = await addressManagementProvider.getAddresses();
      if (response['status'] == true) {
        withdrawAddresses.assignAll(
          List<WithdrawAddressModel>.from(
            response["data"].map(
              (model) {
                if (currencyArray.length > 0) {
                  for (var currency in currencyArray) {
                    if (currency.name == model['code']) {
                      model["icon"] = currency.image;
                      break;
                    }
                  }
                }
                return WithdrawAddressModel.fromJson(model);
              },
            ),
          ),
        );
      }
    } catch (e) {
    } finally {
      isSilentLoading.value = false;

      loadingData.value = false;
    }
    return Future.value();
  }

  handlePullToRefresh() {
    getAddresses(silent: true);
  }

  handleDeleteClick(int index) {
    openConfirmation(
      onConfirm: () {
        deleteAddresses([withdrawAddresses[index].id]);
      },
      titleWidget: RichText(
        text: TextSpan(
          children: [
            TextSpan(
              text: 'Delete ',
              style: TextStyle(
                color: ColorName.red,
              ),
            ),
            TextSpan(
              text: '"' + withdrawAddresses[index].label + '" ?',
              style: TextStyle(
                color: ColorName.white,
              ),
            )
          ],
        ),
      ),
      confirmText: 'Delete',
    );
  }

  void deleteAddresses(List<int> list) async {
    try {
      isRefreshing.value = true;
      final response =
          await addressManagementProvider.deleteAddresses(ids: list);
      if (response['status'] == true) {
        await getAddresses(silent: true);
      }
    } catch (e) {
    } finally {
      loadingData.value = false;
      isRefreshing.value = false;
    }
  }
}
