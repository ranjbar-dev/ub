import '../../../../services/apiService.dart';
import 'package:meta/meta.dart' show required;

class FavoritePairsProvider {
  Future toggleFav({
    @required int pairId,
    @required bool add,
  }) async {
    final response = await apiService.post(
      data: {
        'pair_currency_id': pairId,
        'action': add == true ? 'add' : 'remove'
      },
      url: 'currencies/favorite',
    );
    return response;
  }
}
