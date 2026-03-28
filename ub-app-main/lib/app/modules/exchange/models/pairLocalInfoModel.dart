import 'package:get/get.dart';
import '../../../global/autocompleteModel.dart';

class PairLocalInfoModel {
  PairLocalInfoModel(
      {this.activePairID,
      this.activePairName,
      this.basisBalance,
      this.dependentBalance,
      this.pairPrecision,
      this.basisCoin,
      this.dependantCoin,
      this.possiblePairs,
      this.type});

  RxInt activePairID;
  RxString activePairName;
  RxDouble basisBalance;
  RxDouble dependentBalance;
  RxInt pairPrecision;
  Rx<AutoCompleteItem> basisCoin;
  Rx<AutoCompleteItem> dependantCoin;
  RxList<AutoCompleteItem> possiblePairs;
  RxString type;
}
