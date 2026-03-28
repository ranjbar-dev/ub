import 'package:meta/meta.dart' show required;

import '../../../../services/apiService.dart';
import '../update_user_info_model.dart';

class IdentityInfoProvider {
  Future updateUserInfo({@required UpdateUserInfoModel data}) async {
    final response = await apiService.post(
      url: 'user/set-user-profile',
      data: data,
    );
    return response;
  }
}
