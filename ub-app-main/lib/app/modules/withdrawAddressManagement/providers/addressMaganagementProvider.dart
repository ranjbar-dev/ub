import '../../../../services/apiService.dart';

class AddressManagementProvider {
  Future getAddresses() async {
    final response = await apiService.get(
      url: "withdraw-address",
    );
    return response;
  }

  Future deleteAddresses({List<int> ids}) async {
    final response = await apiService.post(
      data: {"ids": ids},
      url: "withdraw-address/delete",
    );
    return response;
  }
}
