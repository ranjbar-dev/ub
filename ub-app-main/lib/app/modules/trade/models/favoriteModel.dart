class FavoritePairModel {
  int id;
  String name;
  FavoritePairModel({
    this.id,
    this.name,
  });

  FavoritePairModel.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    id = json['id'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['name'] = name;
    return data;
  }
}
