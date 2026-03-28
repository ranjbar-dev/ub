import '../../../services/apiService.dart';

class CommonDataProvider {
  Future getCountries() async {
    final response = await apiService.get(
      url: "main-data/country-list",
    );
    return response;
  }

  Future getCurrencies() async {
    final response = await apiService.get(
      url: "currencies",
    );
    return response;
  }

  Future getUserData() async {
    final response = await apiService.get(
      url: "user/user-data",
    );
    return response;
  }

  Future getFavoritePairs() async {
    final response = await apiService.get(
      url: "currencies/favorite-pairs",
    );
    return response;
  }

  Future getPairs() async {
    final response = await apiService.get(
      url: "currencies/pairs-list",
    );
    return response;
  }

  Future getAppVerion({String platform, String currentVersion}) async {
    final response = await apiService.get(
      url:
          'main-data/version?platform=$platform&current_version=$currentVersion',
    );
    return response;
  }
}
