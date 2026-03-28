import 'package:get/get_utils/src/get_utils/get_utils.dart';

String validateEmail(String email) {
  if (email == '') {
    return 'Email is required';
  }
  if (!GetUtils.isEmail(email)) {
    return 'Email is invalid';
  } else {
    return '';
  }
}
