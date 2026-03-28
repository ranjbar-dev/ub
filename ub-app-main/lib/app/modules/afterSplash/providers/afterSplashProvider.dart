import '../../../../services/apiService.dart';

class AfterSplashProvider {
  Future getAppVerion({String platform, String currentVersion}) async {
    final response = await apiService.get(
      url:
          'main-data/version?platform=$platform&current_version=$currentVersion',
    );
    return response;
  }
}
