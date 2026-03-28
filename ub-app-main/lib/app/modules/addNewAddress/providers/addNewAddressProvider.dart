import '../add_new_address_model.dart';
import '../../../../services/apiService.dart';
import 'package:meta/meta.dart' show required;

class AddNewAddressProvider {
  Future addNewAddress({@required AddNewAddressModel data}) async {
    return await apiService.post(
      url: "withdraw-address/new",
      data: data,
    );
  }
}
