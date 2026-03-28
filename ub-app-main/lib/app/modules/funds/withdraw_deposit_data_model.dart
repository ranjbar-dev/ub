class WithdrawDepositDataModel {
  Balance balance;
  String completedNetworkName;
  CurrencyExtraInfo currencyExtraInfo;
  bool isDepositPermissionGranted;
  bool isWithdrawPermissionGranted;
  List<String> withdrawComments;
  List<String> depositComments;
  String mainNetwork;
  List<OtherNetworksConfigsAndAddresses> otherNetworksConfigsAndAddresses;
  List<OtherNetworksConfigsAndAddresses> networksConfigsAndAddresses;
  bool supportsDeposit;
  bool supportsWithdraw;
  String walletAddress;

  WithdrawDepositDataModel(
      {balance,
      completedNetworkName,
      currencyExtraInfo,
      isDepositPermissionGranted,
      isWithdrawPermissionGranted,
      mainNetwork,
      otherNetworksConfigsAndAddresses,
      networksConfigsAndAddresses,
      supportsDeposit,
      withdrawComments,
      depositComments,
      supportsWithdraw,
      walletAddress});

  WithdrawDepositDataModel.fromJson(Map<String, dynamic> json) {
    balance =
        json['balance'] != null ? Balance.fromJson(json['balance']) : null;
    completedNetworkName = json['completedNetworkName'];
    currencyExtraInfo = json['currencyExtraInfo'] != null
        ? CurrencyExtraInfo.fromJson(json['currencyExtraInfo'])
        : null;
    isDepositPermissionGranted = json['isDepositPermissionGranted'];

    isWithdrawPermissionGranted = json['isWithdrawPermissionGranted'];
    mainNetwork = json['mainNetwork'];
    if (json['otherNetworksConfigsAndAddresses'] != null) {
      otherNetworksConfigsAndAddresses = <OtherNetworksConfigsAndAddresses>[];
      json['otherNetworksConfigsAndAddresses'].forEach((v) {
        otherNetworksConfigsAndAddresses
            .add(OtherNetworksConfigsAndAddresses.fromJson(v));
      });
    }
    if (json['depositComments'] != null) {
      depositComments = <String>[];
      json['depositComments'].forEach((v) {
        depositComments.add(v);
      });
    }

    if (json['withdrawComments'] != null) {
      withdrawComments = <String>[];
      json['withdrawComments'].forEach((v) {
        withdrawComments.add(v);
      });
    }

    if (json['networksConfigsAndAddresses'] != null) {
      networksConfigsAndAddresses = <OtherNetworksConfigsAndAddresses>[];
      json['networksConfigsAndAddresses'].forEach((v) {
        networksConfigsAndAddresses
            .add(OtherNetworksConfigsAndAddresses.fromJson(v));
      });
    }
    supportsDeposit = json['supportsDeposit'];
    supportsWithdraw = json['supportsWithdraw'];
    walletAddress = json['walletAddress'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    if (balance != null) {
      data['balance'] = balance.toJson();
    }
    data['completedNetworkName'] = completedNetworkName;
    if (currencyExtraInfo != null) {
      data['currencyExtraInfo'] = currencyExtraInfo.toJson();
    }
    data['isDepositPermissionGranted'] = isDepositPermissionGranted;
    data['isWithdrawPermissionGranted'] = isWithdrawPermissionGranted;
    data['mainNetwork'] = mainNetwork;
    if (otherNetworksConfigsAndAddresses != null) {
      data['otherNetworksConfigsAndAddresses'] =
          otherNetworksConfigsAndAddresses.map((v) => v.toJson()).toList();
    }
    if (depositComments != null) {
      data['depositComments'] = depositComments.map((v) => v).toList();
    }
    if (withdrawComments != null) {
      data['withdrawComments'] = withdrawComments.map((v) => v).toList();
    }

    if (networksConfigsAndAddresses != null) {
      data['networksConfigsAndAddresses'] =
          networksConfigsAndAddresses.map((v) => v.toJson()).toList();
    }
    data['supportsDeposit'] = supportsDeposit;

    data['supportsWithdraw'] = supportsWithdraw;
    data['walletAddress'] = walletAddress;
    return data;
  }
}

class Balance {
  String address;
  String availableAmount;
  String backgroundImage;
  String btcAvailableEquivalentAmount;
  String btcInOrderEquivalentAmount;
  String btcTotalEquivalentAmount;
  String code;
  String equivalentAvailableAmount;
  String equivalentInOrderAmount;
  String equivalentTotalAmount;
  String fee;
  String image;
  String inOrderAmount;
  String minimumWithdraw;
  String name;
  String price;
  int subUnit;
  String totalAmount;

  Balance(
      {address,
      availableAmount,
      backgroundImage,
      btcAvailableEquivalentAmount,
      btcInOrderEquivalentAmount,
      btcTotalEquivalentAmount,
      code,
      equivalentAvailableAmount,
      equivalentInOrderAmount,
      equivalentTotalAmount,
      fee,
      image,
      inOrderAmount,
      minimumWithdraw,
      name,
      price,
      subUnit,
      totalAmount});

  Balance.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    availableAmount = json['availableAmount'];
    backgroundImage = json['backgroundImage'];
    btcAvailableEquivalentAmount = json['btcAvailableEquivalentAmount'];
    btcInOrderEquivalentAmount = json['btcInOrderEquivalentAmount'];
    btcTotalEquivalentAmount = json['btcTotalEquivalentAmount'];
    code = json['code'];
    equivalentAvailableAmount = json['equivalentAvailableAmount'];
    equivalentInOrderAmount = json['equivalentInOrderAmount'];
    equivalentTotalAmount = json['equivalentTotalAmount'];
    fee = json['fee'];
    image = json['image'];
    inOrderAmount = json['inOrderAmount'];
    minimumWithdraw = json['minimumWithdraw'];
    name = json['name'];
    price = json['price'];
    subUnit = json['subUnit'];
    totalAmount = json['totalAmount'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['availableAmount'] = availableAmount;
    data['backgroundImage'] = backgroundImage;
    data['btcAvailableEquivalentAmount'] = btcAvailableEquivalentAmount;
    data['btcInOrderEquivalentAmount'] = btcInOrderEquivalentAmount;
    data['btcTotalEquivalentAmount'] = btcTotalEquivalentAmount;
    data['code'] = code;
    data['equivalentAvailableAmount'] = equivalentAvailableAmount;
    data['equivalentInOrderAmount'] = equivalentInOrderAmount;
    data['equivalentTotalAmount'] = equivalentTotalAmount;
    data['fee'] = fee;
    data['image'] = image;
    data['inOrderAmount'] = inOrderAmount;
    data['minimumWithdraw'] = minimumWithdraw;
    data['name'] = name;
    data['price'] = price;
    data['subUnit'] = subUnit;
    data['totalAmount'] = totalAmount;
    return data;
  }
}

class CurrencyExtraInfo {
  String circulation;
  String description;
  String image;
  String issueDate;
  String name;
  String totalAmount;

  CurrencyExtraInfo(
      {circulation, description, image, issueDate, name, totalAmount});

  CurrencyExtraInfo.fromJson(Map<String, dynamic> json) {
    circulation = json['circulation'];
    description = json['description'];
    image = json['image'];
    issueDate = json['issueDate'];
    name = json['name'];
    totalAmount = json['totalAmount'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['circulation'] = circulation;
    data['description'] = description;
    data['image'] = image;
    data['issueDate'] = issueDate;
    data['name'] = name;
    data['totalAmount'] = totalAmount;
    return data;
  }
}

class OtherNetworksConfigsAndAddresses {
  String address;
  String code;
  String completedNetworkName;
  bool supportsDeposit;
  bool supportsWithdraw;
  String fee;

  OtherNetworksConfigsAndAddresses(
      {this.address,
      this.code,
      this.completedNetworkName,
      this.fee,
      this.supportsDeposit,
      this.supportsWithdraw});

  OtherNetworksConfigsAndAddresses.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    code = json['code'];
    completedNetworkName = json['completedNetworkName'];
    supportsDeposit = json['supportsDeposit'];
    supportsWithdraw = json['supportsWithdraw'];
    fee = json['fee'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['code'] = code;
    data['completedNetworkName'] = completedNetworkName;
    data['supportsDeposit'] = supportsDeposit;
    data['supportsWithdraw'] = supportsWithdraw;
    data['fee'] = fee;
    return data;
  }
}
