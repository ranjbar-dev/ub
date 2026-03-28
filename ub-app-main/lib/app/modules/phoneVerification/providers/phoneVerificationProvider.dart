import 'package:meta/meta.dart' show required;
import '../phone_verification_post_model.dart';
import '../smsRequestModel.dart';

import '../../../../services/apiService.dart';

class PhoneVerificationProvider {
  Future requestSMS({@required SMSRequestModel data}) async {
    final response = await apiService.post(
      url: 'user/sms-send',
      data: data,
    );
    return response;
  }

  Future submitSMSVerificationCode(
      {@required PhoneVerificationPostModel data}) async {
    final response = await apiService.post(
      url: 'user/sms-enable',
      data: data,
    );
    return response;
  }
}
