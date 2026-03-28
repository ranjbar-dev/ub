class UpdateUserInfoModel {
  String address;
  int countryId;
  String dateOfBirth;
  String firstName;
  String gender;
  String lastName;
  String postalCode;
  String regionAndCity;

  UpdateUserInfoModel(
      {address,
      countryId,
      dateOfBirth,
      firstName,
      gender,
      lastName,
      postalCode,
      regionAndCity});

  UpdateUserInfoModel.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    countryId = json['country_id'];
    dateOfBirth = json['date_of_birth'];
    firstName = json['first_name'];
    gender = json['gender'];
    lastName = json['last_name'];
    postalCode = json['postal_code'];
    regionAndCity = json['region_and_city'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['country_id'] = countryId;
    data['date_of_birth'] = dateOfBirth;
    data['first_name'] = firstName;
    data['gender'] = gender;
    data['last_name'] = lastName;
    data['postal_code'] = postalCode;
    data['region_and_city'] = regionAndCity;
    return data;
  }
}
