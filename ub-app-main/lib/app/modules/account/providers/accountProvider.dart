import '../../../../services/apiService.dart';

class AccountProvider {
  Future requestForEmail() async {
    final response = await apiService.post(
      url: 'user/send-verification-email',
      data: {},
    );
    return response;
  }
}
