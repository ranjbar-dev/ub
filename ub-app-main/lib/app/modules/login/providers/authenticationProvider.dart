import 'package:meta/meta.dart' show required;

import '../../../../services/apiService.dart';
import '../models/login_model.dart';

class AuthenticationProvider {
  Future login({@required LoginModel data}) async {
    final response = await apiService.post(
      url: "auth/login",
      data: data,
    );
    return response;
  }
}
