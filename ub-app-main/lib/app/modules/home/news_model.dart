class NewsModel {
  int id;
  String title;
  String date;
  String content;
  String mainContent;
  CreatedBy createdBy;
  CreatedBy updatedBy;
  String createdAt;
  String updatedAt;
  String slug;
  MainImage mainImage;
  List<dynamic> categories;
  List<Tags> tags;

  NewsModel(
      {id,
      title,
      date,
      content,
      mainContent,
      createdBy,
      updatedBy,
      slug,
      createdAt,
      updatedAt,
      mainImage,
      categories,
      tags});

  NewsModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    title = json['title'];
    slug = json['slug'];
    date = json['date'];
    content = json['content'];
    mainContent = json['main_content'];
    createdBy = json['created_by'] != null
        ? CreatedBy.fromJson(json['created_by'])
        : null;
    updatedBy = json['updated_by'] != null
        ? CreatedBy.fromJson(json['updated_by'])
        : null;
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
    mainImage = json['mainImage'] != null
        ? MainImage.fromJson(json['mainImage'])
        : null;

    if (json['tags'] != null) {
      tags = <Tags>[];
      json['tags'].forEach((v) {
        tags.add(Tags.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['title'] = title;
    data['slug'] = slug;
    data['date'] = date;
    data['content'] = content;
    data['main_content'] = mainContent;
    if (createdBy != null) {
      data['created_by'] = createdBy.toJson();
    }
    if (updatedBy != null) {
      data['updated_by'] = updatedBy.toJson();
    }
    data['created_at'] = createdAt;
    data['updated_at'] = updatedAt;
    if (mainImage != null) {
      data['mainImage'] = mainImage.toJson();
    }

    if (tags != null) {
      data['tags'] = tags.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class CreatedBy {
  int id;
  String firstname;
  String lastname;
  String username;

  CreatedBy({id, firstname, lastname, username});

  CreatedBy.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    firstname = json['firstname'];
    lastname = json['lastname'];
    username = json['username'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['firstname'] = firstname;
    data['lastname'] = lastname;
    data['username'] = username;
    return data;
  }
}

class MainImage {
  int id;
  String name;
  String alternativeText;
  String caption;
  int width;
  int height;
  Formats formats;
  String hash;
  String ext;
  String mime;
  double size;
  String url;
  String previewUrl;
  String provider;
  String providerMetadata;
  int createdBy;
  int updatedBy;
  String createdAt;
  String updatedAt;

  MainImage(
      {id,
      name,
      alternativeText,
      caption,
      width,
      height,
      formats,
      hash,
      ext,
      mime,
      size,
      url,
      previewUrl,
      provider,
      providerMetadata,
      createdBy,
      updatedBy,
      createdAt,
      updatedAt});

  MainImage.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    name = json['name'];
    alternativeText = json['alternativeText'];
    caption = json['caption'];
    width = json['width'];
    height = json['height'];
    formats =
        json['formats'] != null ? Formats.fromJson(json['formats']) : null;
    hash = json['hash'];
    ext = json['ext'];
    mime = json['mime'];
    size = json['size'];
    url = json['url'];
    previewUrl = json['previewUrl'];
    provider = json['provider'];
    providerMetadata = json['provider_metadata'];
    createdBy = json['created_by'];
    updatedBy = json['updated_by'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['name'] = name;
    data['alternativeText'] = alternativeText;
    data['caption'] = caption;
    data['width'] = width;
    data['height'] = height;
    if (formats != null) {
      data['formats'] = formats.toJson();
    }
    data['hash'] = hash;
    data['ext'] = ext;
    data['mime'] = mime;
    data['size'] = size;
    data['url'] = url;
    data['previewUrl'] = previewUrl;
    data['provider'] = provider;
    data['provider_metadata'] = providerMetadata;
    data['created_by'] = createdBy;
    data['updated_by'] = updatedBy;
    data['created_at'] = createdAt;
    data['updated_at'] = updatedAt;
    return data;
  }
}

class Formats {
  Thumbnail thumbnail;
  Thumbnail medium;
  Thumbnail small;

  Formats({thumbnail, medium, small});

  Formats.fromJson(Map<String, dynamic> json) {
    thumbnail = json['thumbnail'] != null
        ? Thumbnail.fromJson(json['thumbnail'])
        : null;
    medium = json['medium'] != null ? Thumbnail.fromJson(json['medium']) : null;
    small = json['small'] != null ? Thumbnail.fromJson(json['small']) : null;
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    if (thumbnail != null) {
      data['thumbnail'] = thumbnail.toJson();
    }
    if (medium != null) {
      data['medium'] = medium.toJson();
    }
    if (small != null) {
      data['small'] = small.toJson();
    }
    return data;
  }
}

class Thumbnail {
  String name;
  String hash;
  String ext;
  String mime;
  int width;
  int height;
  double size;
  String path;
  String url;

  Thumbnail({name, hash, ext, mime, width, height, size, path, url});

  Thumbnail.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    hash = json['hash'];
    ext = json['ext'];
    mime = json['mime'];
    width = json['width'];
    height = json['height'];
    path = json['path'];
    url = json['url'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['name'] = name;
    data['hash'] = hash;
    data['ext'] = ext;
    data['mime'] = mime;
    data['width'] = width;
    data['height'] = height;
    data['size'] = size;
    data['path'] = path;
    data['url'] = url;
    return data;
  }
}

class Tags {
  int id;
  String tagName;
  String blogPost;
  int createdBy;
  int updatedBy;
  String createdAt;
  String updatedAt;

  Tags({id, tagName, blogPost, createdBy, updatedBy, createdAt, updatedAt});

  Tags.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    tagName = json['tagName'];
    blogPost = json['blog_post'];
    createdBy = json['created_by'];
    updatedBy = json['updated_by'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['tagName'] = tagName;
    data['blog_post'] = blogPost;
    data['created_by'] = createdBy;
    data['updated_by'] = updatedBy;
    data['created_at'] = createdAt;
    data['updated_at'] = updatedAt;
    return data;
  }
}
