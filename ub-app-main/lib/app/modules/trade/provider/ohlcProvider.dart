import '../../../../services/constants.dart';

import '../../../../services/apiService.dart';

class OHLCProvider {
  Future getOHLC({
    String symbol,
    String resolution,
    String from,
    String to,
  }) async {
    final response = await apiService.get(
      url: "get-bars?symbol=$symbol&resolution=$resolution&from=$from&to=$to",
      urlGenerator: Constants.generatetvUrl,
    );
    return response;
  }
}
