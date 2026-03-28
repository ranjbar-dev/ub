import 'package:meta/meta.dart' show required;
import '../change_password_model.dart';

import '../../../../services/apiService.dart';

class ChangePasswordProvider {
  Future changePassword({@required ChangePasswordModel data}) async {
    final response = await apiService.post(
      url: 'user/change-password',
      data: data,
    );
    return response;
  }
}
