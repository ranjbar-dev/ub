import 'package:meta/meta.dart' show required;

import '../../../../services/apiService.dart';
import '../toggle2fa_model.dart';

class TwoFactorAuthenticationProviderProvider {
  Future getCharacterCode() async {
    final response = await apiService.get(
      url: 'user/google-2fa-barcode',
    );
    return response;
  }

  Future toggle2Fa(
      {@required Toggle2faModel data, @required bool setEnabled}) async {
    final response = await apiService.post(
      data: data,
      url: "user/google-2fa-${setEnabled == true ? 'enable' : 'disable'}",
    );
    return response;
  }
}
