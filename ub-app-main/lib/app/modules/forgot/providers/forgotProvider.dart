import 'package:meta/meta.dart' show required;
import '../forgot_password_model.dart';

import '../../../../services/apiService.dart';

class ForgotProvider {
  Future forgot({@required ForgotPasswordModel data}) async {
    final response = await apiService.post(
      url: "auth/forgot-password",
      data: data,
    );
    return response;
  }
}
