import 'package:dio/dio.dart';
import 'package:get/get_rx/get_rx.dart';
import 'package:meta/meta.dart' show required;
import '../../../../services/apiService.dart';

class IdentityDocumentsProvider {
  Future getUserProfile() async {
    return await apiService.get(
      url: "user/get-user-profile",
    );
  }

  Future upload({
    @required FormData form,
    @required RxInt stream,
    @required CancelToken cancelToken,
  }) async {
    const url = 'user-profile-image/multiple-upload';
    return await apiService.upload(
        form: form, stream: stream, url: url, cancelToken: cancelToken);
  }
}
