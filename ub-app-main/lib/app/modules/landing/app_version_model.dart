class AppVersionModel {
  String version;
  bool forceUpdate;
  List<String> keyFeatures;
  List<String> bugFixes;
  String releaseDate;
  String url;

  AppVersionModel(
      {version, forceUpdate, keyFeatures, bugFixes, releaseDate, url});

  AppVersionModel.fromJson(Map<String, dynamic> json) {
    version = json['version'];
    forceUpdate = json['forceUpdate'];
    if (json['keyFeatures'] != null) {
      keyFeatures = json['keyFeatures'].cast<String>();
    }
    if (json['bugFixes'] != null) {
      bugFixes = json['bugFixes'].cast<String>();
    }
    releaseDate = json['releaseDate'];
    url = json['url'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['version'] = version;
    data['forceUpdate'] = forceUpdate;
    data['keyFeatures'] = keyFeatures;
    data['bugFixes'] = bugFixes;
    data['releaseDate'] = releaseDate;
    data['url'] = url;
    return data;
  }
}
