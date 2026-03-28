import 'package:meta/meta.dart' show required;
import '../signup_model.dart';

import '../../../../services/apiService.dart';

class SignupProvider {
  Future signup({@required SignupModel data}) async {
    final response = await apiService.post(
      url: "auth/register",
      data: data,
    );
    return response;
  }
}
