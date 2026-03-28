import 'dart:convert';

import "package:meta/meta.dart" show required;

class AutoCompleteItem {
  final String name;
  final String desc;
  final bool favourite;
  final String image;
  final String code;
  final String value;
  final String searchPhrase;
  final String inPerentesis;
  final int id;
  final dynamic otherNetworks;
  final dynamic mainNetwork;
  AutoCompleteItem(
      {this.searchPhrase,
      this.value,
      this.otherNetworks,
      this.code,
      this.inPerentesis,
      this.mainNetwork,
      this.image,
      this.id,
      @required this.name,
      this.desc,
      this.favourite});

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'desc': desc,
      'favourite': favourite,
      'image': image,
      'code': code,
      'searchPhrase': searchPhrase,
      'value': value,
      'inPerentesis': inPerentesis,
      'id': id,
      'otherNetworks': otherNetworks,
      'mainNetwork': mainNetwork,
    };
  }

  factory AutoCompleteItem.fromMap(Map<String, dynamic> map) {
    return AutoCompleteItem(
      name: map['name'],
      desc: map['desc'],
      favourite: map['favourite'],
      image: map['image'],
      code: map['code'],
      searchPhrase: map['searchPhrase'],
      value: map['value'],
      inPerentesis: map['inPerentesis'],
      id: map['id'],
      otherNetworks: map['otherNetworks'],
      mainNetwork: map['mainNetwork'],
    );
  }

  String toJson() => json.encode(toMap());

  factory AutoCompleteItem.fromJson(String source) =>
      AutoCompleteItem.fromMap(json.decode(source));
}
